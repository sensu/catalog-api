package main

import (
	"encoding/base64"
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
	zerolog "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	catalogv1 "github.com/sensu/catalog-api/api/catalog/v1"
	"golang.org/x/mod/sumdb/dirhash"
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

func (m integrationManager) GetNamespacedIntegrations() (namespacedIntegrations, error) {
	tags, err := m.repo.Tags()
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Failed to determine git tags")
	}

	nsIntegrations := namespacedIntegrations{}
	tags.ForEach(func(tagRef *plumbing.Reference) error {
		gitTag := tagRef.Name().Short()
		gitRef := tagRef.Hash().String()
		logger := log.With().Str("tag", gitTag).Logger()

		expr := fmt.Sprintf("^integrations/(?P<IntegrationNamespace>[a-z0-9_-]+)/(?P<IntegrationName>[a-z0-9_-]+)/%s$", semverRegex)
		r := regexp.MustCompile(expr)
		groupNames := r.SubexpNames()
		groupValues := r.FindStringSubmatch(gitTag)
		if len(groupValues) == 0 {
			logger.Warn().Msg("Skipping unmatched git tag")
			return nil
		}
		logger.Debug().Msg("Found matching git tag")

		groups := map[string]string{}
		for i, groupValue := range groupValues {
			groups[groupNames[i]] = groupValue
		}

		name, ok := groups["IntegrationName"]
		if !ok {
			logger.Error().Str("key", "IntegrationName").Msg("Key not found in regex match")
			return nil
		}
		namespace, ok := groups["IntegrationNamespace"]
		if !ok {
			logger.Error().Str("key", "IntegrationNamespace").Msg("Key not found in regex match")
			return nil
		}
		major, err := strconv.Atoi(groups["Major"])
		if err != nil {
			logger.Err(err).Str("major", groups["Major"]).Msg("Failed to convert major version to integer")
			return nil
		}
		minor, err := strconv.Atoi(groups["Minor"])
		if err != nil {
			logger.Err(err).Str("minor", groups["Minor"]).Msg("Failed to convert minor version to integer")
			return nil
		}
		patch, err := strconv.Atoi(groups["Patch"])
		if err != nil {
			logger.Err(err).Str("patch", groups["Patch"]).Msg("Failed to convert patch version to integer")
			return nil
		}
		prerelease, ok := groups["Prerelease"]
		if !ok {
			logger.Error().Str("key", "Prerelease").Msg("Key not found in regex match")
			return nil
		}
		buildMetadata, ok := groups["BuildMetadata"]
		if !ok {
			logger.Error().Str("key", "BuildMetadata").Msg("Key not found in regex match")
			return nil
		}

		nsIntegration, ok := nsIntegrations[namespace]
		if !ok {
			nsIntegration = versionedIntegrations{}
		}
		versions, ok := nsIntegration[name]
		if !ok {
			versions = []integrationVersion{}
		}

		iv := integrationVersion{
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

		versions = append(versions, iv)
		nsIntegration[name] = versions
		nsIntegrations[namespace] = nsIntegration

		logger.Info().
			Str("name", name).
			Str("namespace", namespace).
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

	return contents, nil
}

func (m integrationManager) getIntegrationLogo(version integrationVersion, integrationPath string) (string, error) {
	filePath := path.Join(integrationPath, "logo.png")
	contents, err := m.getFileContentsAtGitRef(version.GitRef, filePath)
	if err != nil {
		return "", fmt.Errorf("error reading contents of integration logo")
	}

	return contents, nil
}

func (m integrationManager) ProcessIntegrationNamepaces() (fErr error) {
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

	// loop through the list of namespaces & integrations, unmarshal the configs
	// & resource files, and then generate the static api
	for namespace, vis := range nsIntegrations {
		logger := log.With().Str("namespace", namespace).Logger()

		processedIntegrations, err := m.ProcessIntegrations(logger, vis, stagingDir)
		if err != nil {
			// no-op
		}

		noop(processedIntegrations)
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

func calculateDirChecksum(path string, prefix string) (string, error) {
	// calculate sha256 checksum, which is returned as a base64 encoded string
	// prefixed with "h1:"
	h1, err := dirhash.HashDir(path, prefix, dirhash.Hash1)
	if err != nil {
		return "", fmt.Errorf("error calculating checksum of dir: %w", err)
	}

	// remove the "h1:" prefix
	re := regexp.MustCompile(`^h1:`)
	b64 := re.ReplaceAllString(h1, "")

	// base64 decode string
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("error base64 decoding checksum: %w", err)
	}

	return fmt.Sprintf("%x", bytes), nil
}

func (m integrationManager) ProcessIntegrations(logger zerolog.Logger, vis versionedIntegrations, basePath string) (processedIntegrations [][]integrationVersion, fErr error) {
	for name, versions := range vis {
		logger := log.With().Str("name", name).Logger()

		processedVersions, err := m.ProcessIntegration(logger, versions, basePath)
		if err != nil {
			logger.Err(err).Msg("Failed to process integration")
			if fErr == nil {
				fErr = errors.New("failed to process one or more integrations")
			}
			continue
		}

		noop(processedVersions)

		processedIntegrations = append(processedIntegrations, versions)
	}
	return
}

func (m integrationManager) ProcessIntegration(logger zerolog.Logger, versions []integrationVersion, basePath string) (processedVersions []integrationVersion, fErr error) {
	for _, version := range versions {
		logger := log.With().Str("version", version.SemVer()).Logger()
		logger.Debug().Msg("Processing integration version")

		// attempt to get files from git
		integrationPath := path.Join("integrations", version.Namespace, version.Name)
		err := m.ProcessIntegrationVersion(version, integrationPath, basePath)
		if err != nil {
			logger.Err(err).Msg("Failed to process integration version")
			if fErr == nil {
				fErr = errors.New("failed to process one or more integration versions")
			}
		}
		processedVersions = append(processedVersions, version)
	}
	return
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

	generateIntegrationVersionEndpoint(basePath, integration, version)

	return nil
}

func noop(_ interface{}) {}
