package catalogmanager

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"

	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/catalogloader"
	"github.com/sensu/catalog-api/internal/endpoints"
	"github.com/sensu/catalog-api/internal/types"
	"github.com/sensu/catalog-api/internal/util"
)

type CatalogManager struct {
	config Config
	loader catalogloader.Loader
}

func (m CatalogManager) GetConfig() Config {
	return m.config
}

func New(config Config, loader catalogloader.Loader) (CatalogManager, error) {
	m := CatalogManager{
		config: config,
		loader: loader,
	}

	if err := config.validate(); err != nil {
		return m, fmt.Errorf("catalog manager config validation failed: %w", err)
	}

	return m, nil
}

func (m CatalogManager) ProcessCatalog() error {
	integrations, err := m.loader.LoadIntegrations()
	if err != nil {
		return fmt.Errorf("error loading integrations: %w", err)
	}

	integrationsByNamespace := integrations.ByNamespace()
	for namespace, nsIntegrations := range integrationsByNamespace {
		if err := m.ProcessNamespace(namespace, nsIntegrations); err != nil {
			return err
		}
	}

	latestNsIntegrations := map[string][]catalogapiv1.IntegrationVersion{}
	for namespace, nsIntegrations := range integrationsByNamespace.FilterByLatestVersions() {
		for _, integration := range nsIntegrations {
			integrationLoader := m.loader.NewIntegrationLoader(namespace, integration.Name, integration.SemVer())
			config, err := integrationLoader.LoadConfig()
			if err != nil {
				return err
			}
			if err := config.Validate(); err != nil {
				return err
			}

			// Set prompts & resource_patches fields to empty strings to prevent
			// them from being shown in the catalog endpoint.
			config.Prompts = nil
			config.ResourcePatches = nil

			iv := catalogapiv1.IntegrationVersion{
				Integration: config,
				Version:     integration.SemVer(),
			}

			latestNsIntegrations[namespace] = append(latestNsIntegrations[namespace], iv)
		}
	}

	if err := endpoints.GenerateCatalogEndpoint(m.config.StagingDir, latestNsIntegrations); err != nil {
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

func (m CatalogManager) ProcessNamespace(namespace string, integrations types.Integrations) error {
	for integration, versions := range integrations.ByName() {
		if err := m.ProcessIntegrationVersions(namespace, integration, versions); err != nil {
			return err
		}
		for _, version := range versions {
			if err := m.ProcessIntegrationVersion(version); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m CatalogManager) ProcessIntegrationVersions(namespace string, integrationName string, integrations types.Integrations) error {
	for _, integration := range integrations {
		if err := m.ProcessIntegrationVersion(integration); err != nil {
			log.Err(err).
				Str("namespace", namespace).
				Str("integration", integrationName).
				Str("version", integration.SemVer()).
				Msg("Failed to process integration version")
			return fmt.Errorf("error processing integration version: %w", err)
		}
	}

	if err := endpoints.GenerateIntegrationVersionsEndpoint(m.config.StagingDir, namespace, integrationName, integrations); err != nil {
		return fmt.Errorf("error generating integration versions endpoint: %w", err)
	}

	latestVersion := integrations.LatestVersion()
	integrationLoader := m.loader.NewIntegrationLoader(latestVersion.Namespace, latestVersion.Name, latestVersion.SemVer())
	integrationConfig, err := integrationLoader.LoadConfig()
	if err != nil {
		return err
	}
	if err := integrationConfig.Validate(); err != nil {
		return err
	}
	if err := endpoints.GenerateIntegrationEndpoint(m.config.StagingDir, integrationConfig, integrations); err != nil {
		return fmt.Errorf("error generating integration endpoint: %w", err)
	}

	return nil
}

func (m CatalogManager) ProcessIntegrationVersion(version types.IntegrationVersion) error {
	integrationLoader := m.loader.NewIntegrationLoader(version.Namespace, version.Name, version.SemVer())

	config, err := integrationLoader.LoadConfig()
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
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

	if err := endpoints.GenerateIntegrationVersionEndpoint(m.config.StagingDir, config, version); err != nil {
		return fmt.Errorf("error generating integration version endpoint: %w", err)
	}
	if err := endpoints.GenerateIntegrationVersionResourcesEndpoint(m.config.StagingDir, config, version, resourcesJSON); err != nil {
		return fmt.Errorf("error generating integration version resources endpoint: %w", err)
	}
	if logo != "" {
		if err := endpoints.GenerateIntegrationVersionLogoEndpoint(m.config.StagingDir, config, version, logo); err != nil {
			return fmt.Errorf("error generating integration version logo endpoint: %w", err)
		}
	}
	if err := endpoints.GenerateIntegrationVersionReadmeEndpoint(m.config.StagingDir, config, version, readme); err != nil {
		return fmt.Errorf("error generating integration version readme endpoint: %w", err)
	}
	if err := endpoints.GenerateIntegrationVersionChangelogEndpoint(m.config.StagingDir, config, version, changelog); err != nil {
		return fmt.Errorf("error generating integration version changelog endpoint: %w", err)
	}

	// iterate through each .jpg file in the img directory and create an
	// endpoint for it
	images, err := integrationLoader.LoadImages()
	if err != nil {
		return fmt.Errorf("error loading integration images: %w", err)
	}
	for imageName, imageData := range images {
		if err := endpoints.GenerateIntegrationVersionImageEndpoint(m.config.StagingDir, config, version, imageName, imageData); err != nil {
			return fmt.Errorf("error generating integration version image endpoint: %w", err)
		}
	}

	return nil
}
