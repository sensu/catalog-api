package main

import (
	"fmt"
	"path"
	"regexp"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
	catalogv1 "github.com/sensu/catalog-api/api/catalog/v1"
	"gopkg.in/yaml.v3"
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
	Processed     bool
}

func (i integrationVersion) BaseVersion() string {
	return fmt.Sprintf("%d.%d.%d", i.Major, i.Minor, i.Patch)
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
			Str("service", "test").
			Msgf("Failed to determine git tags: %w", err)
	}

	nsIntegrations := namespacedIntegrations{}
	tags.ForEach(func(tagRef *plumbing.Reference) error {
		gitTag := tagRef.Name().Short()
		gitRef := tagRef.Hash().String()

		expr := fmt.Sprintf("^integrations/(?P<IntegrationNamespace>[a-z0-9_-]+)/(?P<IntegrationName>[a-z0-9_-]+)/%s$", semverRegex)
		r := regexp.MustCompile(expr)
		groupNames := r.SubexpNames()
		groupValues := r.FindStringSubmatch(gitTag)
		if len(groupValues) == 0 {
			log.Warn().Str("tag", gitTag).Msg("Skipping unmatched git tag")
			return nil
		}
		log.Debug().Str("tag", gitTag).Msg("Found matching git tag")

		groups := map[string]string{}
		for i, groupValue := range groupValues {
			groups[groupNames[i]] = groupValue
		}

		name, ok := groups["IntegrationName"]
		if !ok {
			log.Error().Str("tag", gitTag).Str("key", "IntegrationName").Msg("Key not found in regex match")
			return nil
		}
		namespace, ok := groups["IntegrationNamespace"]
		if !ok {
			log.Error().Str("tag", gitTag).Str("key", "IntegrationNamespace").Msg("Key not found in regex match")
			return nil
		}
		major, err := strconv.Atoi(groups["Major"])
		if err != nil {
			log.Err(err).Str("tag", gitTag).Str("major", groups["Major"]).Msg("Failed to convert major version to integer")
			return nil
		}
		minor, err := strconv.Atoi(groups["Minor"])
		if err != nil {
			log.Err(err).Str("tag", gitTag).Str("minor", groups["Minor"]).Msg("Failed to convert minor version to integer")
			return nil
		}
		patch, err := strconv.Atoi(groups["Patch"])
		if err != nil {
			log.Err(err).Str("tag", gitTag).Str("patch", groups["Patch"]).Msg("Failed to convert patch version to integer")
			return nil
		}
		prerelease, ok := groups["Prerelease"]
		if !ok {
			log.Error().Str("tag", gitTag).Str("key", "Prerelease").Msg("Key not found in regex match")
			return nil
		}
		buildMetadata, ok := groups["BuildMetadata"]
		if !ok {
			log.Error().Str("tag", gitTag).Str("key", "BuildMetadata").Msg("Key not found in regex match")
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
			Processed:     false,
		}

		versions = append(versions, iv)
		nsIntegration[name] = versions
		nsIntegrations[namespace] = nsIntegration

		log.Info().
			Str("name", name).
			Str("namespace", namespace).
			Str("version", iv.BaseVersion()).
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

func noop(_ interface{}) {}

func (m integrationManager) ProcessIntegration(version integrationVersion, integrationPath string) error {
	integration, err := m.getIntegrationConfig(version, integrationPath)
	if err != nil {
		return err
	}
	noop(integration)

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

	return nil
}
