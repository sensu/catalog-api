package catalogmanager

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"

	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/endpoints"
	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
	"github.com/sensu/catalog-api/internal/util"
)

const semverRegex = `(?P<Major>0|[1-9]\d*)\.(?P<Minor>0|[1-9]\d*)\.(?P<Patch>0|[1-9]\d*)(?:-(?P<Prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<BuildMetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`

type CatalogManager struct {
	config Config
	repo   *git.Repository
}

func (m CatalogManager) GetConfig() Config {
	return m.config
}

func New(config Config) (CatalogManager, error) {
	m := CatalogManager{
		config: config,
	}

	if err := config.validate(); err != nil {
		return m, fmt.Errorf("catalog manager config validation failed: %w", err)
	}

	repo, err := git.PlainOpen(config.RepoDir)
	if err != nil {
		return m, err
	}
	m.repo = repo

	return m, nil
}

func getIntegrationVersionFromGitTag(tagRef *plumbing.Reference) (types.IntegrationVersion, error) {
	var iv types.IntegrationVersion

	gitTag := tagRef.Name().Short()
	gitRef := tagRef.Hash().String()

	expr := fmt.Sprintf("^(?P<IntegrationNamespace>[a-z0-9_-]+)/(?P<IntegrationName>[a-z0-9_-]+)/%s$", semverRegex)
	r := regexp.MustCompile(expr)
	groupNames := r.SubexpNames()
	groupValues := r.FindStringSubmatch(gitTag)
	if len(groupValues) == 0 {
		return iv, ErrUnmatchedGitTag
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

	iv = types.IntegrationVersion{
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

func (m CatalogManager) IntegrationsDir() string {
	return path.Join(m.config.RepoDir, m.config.IntegrationsDirName)
}

func (m CatalogManager) GetNamespacedIntegrations() (types.NamespacedIntegrations, error) {
	nsIntegrations := types.NamespacedIntegrations{}

	tags, err := m.repo.Tags()
	if err != nil {
		return nsIntegrations, fmt.Errorf("error determining git tags: %w", err)
	}

	tags.ForEach(func(tagRef *plumbing.Reference) error {
		logger := log.With().Str("tag", tagRef.Name().Short()).Logger()

		iv, err := getIntegrationVersionFromGitTag(tagRef)
		if err != nil {
			if errors.Is(err, ErrUnmatchedGitTag) {
				logger.Warn().Str("reason", err.Error()).Msg("Skipping integration version")
			} else {
				return fmt.Errorf("error parsing git tag - tag: %s, err: %w", tagRef.Name().Short(), err)
			}
			return nil
		}

		nsIntegration, ok := nsIntegrations[iv.Namespace]
		if !ok {
			nsIntegration = types.VersionedIntegrations{}
		}
		versions, ok := nsIntegration[iv.Name]
		if !ok {
			versions = []types.IntegrationVersion{}
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

func (m CatalogManager) ProcessCatalog() error {
	// get a list of namespaces & the integrations that belong to them from the
	// list of git tags
	nsIntegrations, err := m.GetNamespacedIntegrations()
	if err != nil {
		return fmt.Errorf("error retrieving list of integrations from git tags: %w", err)
	}

	// loop through the list of namespaces & integrations, unmarshal the configs
	// & resource files, and then generate the static api
	latestNsIntegrationVersions := map[string][]catalogapiv1.IntegrationVersion{}
	for namespace, vis := range nsIntegrations {
		if err := m.ProcessNamespace(namespace, vis); err != nil {
			return err
		}

		integrationVersions := []catalogapiv1.IntegrationVersion{}
		for integration, versions := range vis {
			integrationPath := path.Join(m.config.IntegrationsDirName, namespace, integration)
			latestVersion := versions.LatestVersion()
			loader := integrationloader.NewGitLoader(m.repo, latestVersion.GitRef, integrationPath)

			// retrieve the integration config for the latest integration version
			integrationConfig, err := loader.LoadConfig()
			if err != nil {
				return fmt.Errorf("error retrieving integration config: %w", err)
			}
			if err := integrationConfig.Validate(); err != nil {
				return fmt.Errorf("integration config: %w", err)
			}

			integrationVersions = append(integrationVersions, catalogapiv1.IntegrationVersion{
				Integration: integrationConfig,
				Version:     latestVersion.SemVer(),
			})
		}

		latestNsIntegrationVersions[namespace] = integrationVersions
	}

	if err := endpoints.GenerateCatalogEndpoint(m.config.StagingDir, latestNsIntegrationVersions); err != nil {
		return fmt.Errorf("error generating catalog endpoint: %w", err)
	}

	// calculate the sha256 checksum of the generated api
	checksum, err := util.CalculateDirChecksum(m.config.StagingDir, "staging")
	if err != nil {
		return fmt.Errorf("error calculating checksum of staging dir: %w", err)
	}

	// copy the staging dir to the release dir
	dstPath := path.Join(m.config.ReleaseDir, checksum)
	cmd := exec.Command("cp", "-R", m.config.StagingDir, dstPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error copying staging files to release dir: %w", err)
	}

	if err := endpoints.GenerateVersionEndpoint(m.config.ReleaseDir, checksum); err != nil {
		return fmt.Errorf("error generating version endpoint: %w", err)
	}

	return nil
}

func (m CatalogManager) ProcessNamespace(namespace string, vis types.VersionedIntegrations) error {
	processedIntegrations, err := m.ProcessIntegrations(namespace, vis)
	if err != nil {
		return fmt.Errorf("error processing namespace: %w", err)
	}

	for integration, versions := range processedIntegrations {
		if err := endpoints.GenerateIntegrationVersionsEndpoint(m.config.StagingDir, namespace, integration, versions); err != nil {
			return fmt.Errorf("error generating integration versions endpoint: %w", err)
		}
	}

	return nil
}

func (m CatalogManager) ProcessIntegrations(namespace string, vis types.VersionedIntegrations) (types.VersionedIntegrations, error) {
	integrations := types.VersionedIntegrations{}

	for integration, versions := range vis {
		if err := m.ProcessIntegration(namespace, integration, versions); err != nil {
			return integrations, fmt.Errorf("error processing integration: %s", err)
		}
		integrations[integration] = versions
	}

	return integrations, nil
}

func (m CatalogManager) ProcessIntegration(namespace string, integration string, versions []types.IntegrationVersion) error {
	integrationPath := path.Join("integrations", namespace, integration)

	var processed types.IntegrationVersionsSlice
	for _, version := range versions {
		if err := m.ProcessIntegrationVersion(version, integrationPath); err != nil {
			log.Err(err).
				Str("namespace", namespace).
				Str("integration", integration).
				Str("version", version.SemVer()).
				Msg("Failed to process integration version")
			return err
		}

		processed = append(processed, version)
	}

	latestVersion := processed.LatestVersion()
	loader := integrationloader.NewGitLoader(m.repo, latestVersion.GitRef, integrationPath)
	integrationConfig, err := loader.LoadConfig()
	if err != nil {
		return err
	}
	if err := integrationConfig.Validate(); err != nil {
		return fmt.Errorf("error validating integration config: %w", err)
	}
	if err := endpoints.GenerateIntegrationEndpoint(m.config.StagingDir, integrationConfig, versions); err != nil {
		return fmt.Errorf("error generating integration endpoint: %w", err)
	}

	return nil
}

func (m CatalogManager) ProcessIntegrationVersion(version types.IntegrationVersion, integrationPath string) error {
	integrationLoader := integrationloader.NewGitLoader(m.repo, version.GitRef, integrationPath)
	integration, err := integrationLoader.LoadConfig()
	if err != nil {
		return err
	}
	if err := integration.Validate(); err != nil {
		return fmt.Errorf("integration config: %w", err)
	}

	resourcesJSON, err := integrationLoader.LoadResources()
	if err != nil {
		return err
	}

	logo, err := integrationLoader.LoadLogo()
	if err != nil {
		// integration logo was found but an error occurred when reading it
		if _, ok := err.(*fs.PathError); !ok {
			return err
		}
	}

	readme, err := integrationLoader.LoadReadme()
	if err != nil {
		return err
	}

	changelog, err := integrationLoader.LoadChangelog()
	if err != nil {
		return err
	}

	if err := endpoints.GenerateIntegrationVersionEndpoint(m.config.StagingDir, integration, version); err != nil {
		return fmt.Errorf("error generating integration version endpoint: %w", err)
	}
	if err := endpoints.GenerateIntegrationVersionResourcesEndpoint(m.config.StagingDir, integration, version, resourcesJSON); err != nil {
		return fmt.Errorf("error generating integration version resources endpoint: %w", err)
	}
	if logo != "" {
		if err := endpoints.GenerateIntegrationVersionLogoEndpoint(m.config.StagingDir, integration, version, logo); err != nil {
			return fmt.Errorf("error generating integration version logo endpoint: %w", err)
		}
	}
	if err := endpoints.GenerateIntegrationVersionReadmeEndpoint(m.config.StagingDir, integration, version, readme); err != nil {
		return fmt.Errorf("error generating integration version readme endpoint: %w", err)
	}
	if err := endpoints.GenerateIntegrationVersionChangelogEndpoint(m.config.StagingDir, integration, version, changelog); err != nil {
		return fmt.Errorf("error generating integration version changelog endpoint: %w", err)
	}

	// iterate through each .jpg file in the img directory and create an
	// endpoint for it
	images, err := integrationLoader.LoadImages()
	if err != nil {
		return fmt.Errorf("error loading integration images: %w", err)
	}

	for imageName, imageData := range images {
		if err := endpoints.GenerateIntegrationVersionImageEndpoint(m.config.StagingDir, integration, version, imageName, imageData); err != nil {
			return fmt.Errorf("error generating integration version image endpoint: %w", err)
		}
	}

	return nil
}
