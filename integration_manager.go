package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"

	catalogv1 "github.com/sensu/catalog-api/api/catalog/v1"
)

const semverRegex = `(?P<Major>0|[1-9]\d*)\.(?P<Minor>0|[1-9]\d*)\.(?P<Patch>0|[1-9]\d*)(?:-(?P<Prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<BuildMetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`

// integrationVersion is a representation of a single integration version
type integrationVersion struct {
	Name          string
	Namespace     string
	Major         int
	Minor         int
	Patch         int
	Prerelease    string
	BuildMetadata string
	GitTag        string
	GitRef        string
}

func (i integrationVersion) String() string {
	return fmt.Sprintf("%s/%s:%s", i.Namespace, i.Name, i.SemVer())
}

func (i integrationVersion) SemVer() string {
	version := fmt.Sprintf("%d.%d.%d", i.Major, i.Minor, i.Patch)

	if i.Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, i.Prerelease)
	}

	if i.BuildMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, i.BuildMetadata)
	}

	return version
}

// versionedIntegrations is a mapping of integration names to
// integrationVersions
type versionedIntegrations map[string][]integrationVersion

// namespacedIntegrations is a mapping of namespaces to versionedIntegrations
type namespacedIntegrations map[string]versionedIntegrations

type integrationManager struct {
	repo *git.Repository
}

func newIntegrationManager(path string) (integrationManager, error) {
	var im integrationManager

	repo, err := git.PlainOpen(path)
	if err != nil {
		return im, err
	}

	im = integrationManager{
		repo: repo,
	}
	return im, nil
}

func getIntegrationVersionFromGitTag(tagRef *plumbing.Reference) (integrationVersion, error) {
	var iv integrationVersion

	gitTag := tagRef.Name().Short()
	gitRef := tagRef.Hash().String()

	expr := fmt.Sprintf("^integrations/(?P<IntegrationNamespace>[a-z0-9_-]+)/(?P<IntegrationName>[a-z0-9_-]+)/%s$", semverRegex)
	r := regexp.MustCompile(expr)
	groupNames := r.SubexpNames()
	groupValues := r.FindStringSubmatch(gitTag)
	if len(groupValues) == 0 {
		return iv, errors.New("skipping unmatched git tag")
	}

	groups := map[string]string{}
	for i, groupValue := range groupValues {
		groups[groupNames[i]] = groupValue
	}

	// verify that all of the required regex group keys were set
	namespace, ok := groups["IntegrationNamespace"]
	if !ok {
		return iv, errors.New("key not found in regex match: IntegrationNamespace")
	}
	name, ok := groups["IntegrationName"]
	if !ok {
		return iv, errors.New("key not found in regex match: IntegrationName")
	}
	majorStr, ok := groups["Major"]
	if !ok {
		return iv, errors.New("key not found in regex match: Major")
	}
	minorStr, ok := groups["Minor"]
	if !ok {
		return iv, errors.New("key not found in regex match: Minor")
	}
	patchStr, ok := groups["Patch"]
	if !ok {
		return iv, errors.New("key not found in regex match: Patch")
	}
	prerelease, ok := groups["Prerelease"]
	if !ok {
		return iv, errors.New("key not found in regex match: Prerelease")
	}
	buildMetadata, ok := groups["BuildMetadata"]
	if !ok {
		return iv, errors.New("key not found in regex match: BuildMetadata")
	}

	// convert major, minor & patch versions to integers
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return iv, fmt.Errorf("error converting major version to integer: %w", err)
	}
	minor, err := strconv.Atoi(minorStr)
	if err != nil {
		return iv, fmt.Errorf("error converting minor version to integer: %w", err)
	}
	patch, err := strconv.Atoi(patchStr)
	if err != nil {
		return iv, fmt.Errorf("error converting patch version to integer: %w", err)
	}

	iv = integrationVersion{
		Name:          name,
		Namespace:     namespace,
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		Prerelease:    prerelease,
		BuildMetadata: buildMetadata,
		GitTag:        gitTag,
		GitRef:        gitRef,
	}

	return iv, nil
}

func (m integrationManager) GetNamespacedIntegrations() (namespacedIntegrations, error) {
	tags, err := m.repo.Tags()
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Failed to determine git tags")
	}

	nsIntegrations := namespacedIntegrations{}
	tags.ForEach(func(tagRef *plumbing.Reference) error {
		logger := log.With().Str("tag", tagRef.Name().Short()).Logger()

		iv, err := getIntegrationVersionFromGitTag(tagRef)
		if err != nil {
			logger.Err(err).Msg("Failed to get integration version from git tag")
			return nil
		}

		nsIntegration, ok := nsIntegrations[iv.Namespace]
		if !ok {
			nsIntegration = versionedIntegrations{}
		}
		versions, ok := nsIntegration[iv.Name]
		if !ok {
			versions = []integrationVersion{}
		}
		versions = append(versions, iv)
		nsIntegration[iv.Name] = versions
		nsIntegrations[iv.Namespace] = nsIntegration

		logger.Info().
			Str("name", iv.Name).
			Str("namespace", iv.Namespace).
			Str("version", iv.SemVer()).
			Msg("Found integration version")

		return nil
	})

	return nsIntegrations, nil
}

func loadFileFromGitTreePaths(tree *object.Tree, filePaths []string) (*object.File, error) {
	for _, filePath := range filePaths {
		file, err := tree.File(filePath)
		if err != nil {
			continue
		}
		return file, nil
	}
	return nil, fmt.Errorf("failed to load file from any of the provided paths: %s", filePaths)
}

func (m integrationManager) getFileContentsAtGitRef(ref string, filePath string) (string, error) {
	// attempt to resolve the git ref to a revision
	hash, err := m.repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", fmt.Errorf("error resolving git revision: %w", err)
	}

	// attempt to retrieve the commit object for the hash
	commit, err := m.repo.CommitObject(*hash)
	if err != nil {
		return "", fmt.Errorf("error retrieving commit for hash: %w", err)
	}

	// attempt to retrieve the directory tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("error retrieving tree for commit: %w", err)
	}

	// attempt to load the file from the tree
	file, err := tree.File(filePath)
	if err != nil {
		return "", fmt.Errorf("error loading file from tree: %w", err)
	}

	// attempt to read the contents of the file
	contents, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("error reading contents of file: %w", err)
	}

	return contents, nil
}

func (m integrationManager) getIntegrationConfig(version integrationVersion, integrationPath string) (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	filePath := path.Join(integrationPath, "sensu-integration.yaml")
	contents, err := m.getFileContentsAtGitRef(version.GitRef, filePath)
	if err != nil {
		return integration, fmt.Errorf("error reading contents of integration config")
	}

	raw, err := rawWrapperFromYAMLBytes([]byte(contents))
	if err != nil {
		return integration, err
	}

	wrap, err := wrapperFromRawWrapper(raw)
	if err != nil {
		return integration, err
	}
	integration = wrap.Value.(catalogv1.Integration)

	return integration, nil
}

func (m integrationManager) getIntegrationResources(version integrationVersion, integrationPath string) (string, error) {
	// TODO(jk): support both .yaml & .yml extensions
	filePath := path.Join(integrationPath, "sensu-resources.yaml")
	contents, err := m.getFileContentsAtGitRef(version.GitRef, filePath)
	if err != nil {
		return "", fmt.Errorf("error reading contents of integration resources")
	}

	// attempt to unmarshal yaml to verify that the yaml is valid
	// TODO(jk): iterate through & validate each resource against the supported
	// versions of Sensu that the integration defines
	node := yaml.Node{}
	if err := yaml.Unmarshal([]byte(contents), &node); err != nil {
		return "", fmt.Errorf("error unmarshaling integration resources: %w", err)
	}

	nj, err := json.Marshal(node)
	if err != nil {
		return "", err
	}

	fmt.Println(string(nj))

	return contents, nil
}

func (m integrationManager) getIntegrationLogo(version integrationVersion, integrationPath string) (string, error) {
	filePath := path.Join(integrationPath, "logo.png")
	contents, err := m.getFileContentsAtGitRef(version.GitRef, filePath)
	if err != nil {
		return "", fmt.Errorf("error reading contents of integration logo")
	}

	// TODO(jk): add basic validation using https://pkg.go.dev/net/http#DetectContentType

	return contents, nil
}

func (m integrationManager) ProcessCatalog() (fErr error) {
	// get a list of namespaces & the integrations that belong to them from the
	// list of git tags
	nsIntegrations, err := m.GetNamespacedIntegrations()
	if err != nil {
		return fmt.Errorf("error retrieving list of integrations from git tags: %w", err)
	}

	// create a temp dir to hold the generated api files
	tmpDir, err := os.MkdirTemp("", "sensu-catalog-api-")
	if err != nil {
		return fmt.Errorf("error creating tmp directory: %w", err)
	}

	// create a staging dir to hold the generated api files used to calculate
	// the checksum of the release
	stagingDir := path.Join(tmpDir, "staging")
	if err := os.Mkdir(stagingDir, 0700); err != nil {
		return fmt.Errorf("error creating staging directory: %w", err)
	}

	processed := namespacedIntegrations{}

	// loop through the list of namespaces & integrations, unmarshal the configs
	// & resource files, and then generate the static api
	for namespace, vis := range nsIntegrations {
		logger := log.With().Str("namespace", namespace).Logger()

		if err := m.ProcessNamespace(stagingDir, namespace, vis); err != nil {
			logger.Err(err).Msg("Failed to process namespace")
			continue
		}

		processed[namespace] = vis
	}

	// calculate the sha256 checksum of the generated api
	checksum, err := calculateDirChecksum(stagingDir, "staging")
	if err != nil {
		return fmt.Errorf("error calculating checksum of release: %w", err)
	}

	// create a release dir to hold the complete set of generated api files
	releaseDir := path.Join(tmpDir, "release")
	if err := os.Mkdir(releaseDir, 0700); err != nil {
		return fmt.Errorf("error creating release directory: %w", err)
	}

	// copy the staging dir to the release dir
	dstPath := path.Join(releaseDir, checksum)
	cmd := exec.Command("cp", "-R", stagingDir, dstPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error copying staging files to release dir: %w", err)
	}

	generateVersionEndpoint(releaseDir, checksum)

	fmt.Println(releaseDir)

	return fErr
}

func (m integrationManager) ProcessNamespace(basePath string, namespace string, vis versionedIntegrations) error {
	processedIntegrations, err := m.ProcessIntegrations(basePath, namespace, vis)
	if err != nil {
		return fmt.Errorf("error processing integrations: %w", err)
	}

	for integration, versions := range processedIntegrations {
		generateIntegrationVersionsEndpoint(basePath, namespace, integration, versions)
	}

	return nil
}

func (m integrationManager) ProcessIntegrations(basePath string, namespace string, vis versionedIntegrations) (versionedIntegrations, error) {
	processed := versionedIntegrations{}
	failed := versionedIntegrations{}

	for integration, versions := range vis {
		if err := m.ProcessIntegration(basePath, namespace, integration, versions); err != nil {
			failed[integration] = versions
			continue
		}

		processed[integration] = versions
	}

	if len(failed) > 0 {
		return processed, fmt.Errorf("error processing integrations: %s", failed)
	}

	return processed, nil
}

func (m integrationManager) ProcessIntegration(basePath string, namespace string, integration string, versions []integrationVersion) error {
	processed := []integrationVersion{}
	failed := []integrationVersion{}
	latestVersion := integrationVersion{}
	integrationPath := path.Join("integrations", namespace, integration)

	for i, version := range versions {
		if err := m.ProcessIntegrationVersion(version, integrationPath, basePath); err != nil {
			failed = append(failed, version)
		}

		processed = append(processed, version)

		if i != 0 {
			if semver.Compare(latestVersion.SemVer(), version.SemVer()) == -1 {
				latestVersion = version
			}
		} else {
			latestVersion = version
		}
	}

	integrationConfig, err := m.getIntegrationConfig(latestVersion, integrationPath)
	if err != nil {
		return err
	}
	generateIntegrationEndpoint(basePath, integrationConfig, versions)

	if len(failed) > 0 {
		return fmt.Errorf("error processing integration versions: %s", failed)
	}

	return nil
}

func (m integrationManager) ProcessIntegrationVersion(version integrationVersion, integrationPath string, basePath string) error {
	integration, err := m.getIntegrationConfig(version, integrationPath)
	if err != nil {
		return err
	}

	resources, err := m.getIntegrationResources(version, integrationPath)
	if err != nil {
		return err
	}
	noop(resources)

	logo, err := m.getIntegrationLogo(version, integrationPath)
	if err != nil {
		return err
	}
	noop(logo)

	// TODO(jk): add README.md
	// TODO(jk): add CHANGELOG.md
	// TODO(jk): add images directory

	generateIntegrationVersionEndpoint(basePath, integration, version)
	//generateIntegration

	return nil
}

func noop(_ interface{}) {}
