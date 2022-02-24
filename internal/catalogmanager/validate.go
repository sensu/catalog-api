package catalogmanager

import (
	"errors"
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/catalogloader"
	"github.com/sensu/catalog-api/internal/integrationloader"
)

func (m CatalogManager) ValidateCatalogDir() error {
	catalogLoader := catalogloader.NewPathLoader(m.IntegrationsDir())
	nsIntegrations, err := catalogLoader.LoadIntegrations()
	if err != nil {
		return fmt.Errorf("error loading integrations from catalog: %w", err)
	}

	// loop through the list of namespaces & integrations, and unmarshal the
	// configs & resource files
	validationFailed := false
	for namespace, integrations := range nsIntegrations {
		for _, integration := range integrations {
			integrationPath := path.Join(m.IntegrationsDir(), namespace, integration)
			integrationLoader := integrationloader.NewPathLoader(integrationPath)

			logger := log.With().
				Str("namespace", namespace).
				Str("integration", integration).
				Logger()

			// load & validate the integration config
			integrationConfig, err := integrationLoader.LoadConfig()
			if err != nil {
				logger.Err(err).Msg("Failed to load integration config")
				validationFailed = true
				continue
			}
			if err := integrationConfig.Validate(); err != nil {
				logger.Err(err).Msg("Failed to validate integration config")
				validationFailed = true
				continue
			}

			// load & validate sensu resources
			if _, err = integrationLoader.LoadResources(); err != nil {
				logger.Err(err).Msg("Failed to load resources file")
			}
			// TODO(jk): call resouces.Validate() once it's implemented

			// load & validate logo
			_, err = integrationLoader.LoadLogo()
			if err != nil {
				logger.Err(err).Msg("Failed to load logo")
			}

			// load & validate readme
			_, err = integrationLoader.LoadReadme()
			if err != nil {
				logger.Err(err).Msg("Failed to load readme")
			}

			// load & validate changelog
			_, err = integrationLoader.LoadReadme()
			if err != nil {
				logger.Err(err).Msg("Failed to load changelog")
			}

			// load & validate images
			_, err = integrationLoader.LoadImages()
			if err != nil {
				logger.Err(err).Msg("Failed to load images")
			}
		}
	}

	if validationFailed {
		return errors.New("one or more integrations failed validation")
	}
	return nil
}

// func (m CatalogManager) ValidateCatalogRepo() error {
// 	integrationsDir := path.Join(m.config.RepoDir, m.config.IntegrationsDirName)

// 	// get a list of namespaces & the integrations that belong to them from the
// 	// list of git tags
// 	nsIntegrations, err := m.GetNamespacedIntegrations()
// 	if err != nil {
// 		return fmt.Errorf("error retrieving list of integrations from git tags: %w", err)
// 	}

// 	// loop through the list of namespaces & integrations, and unmarshal the
// 	// configs & resource files
// 	validationFailed := false
// 	for namespace, integrations := range nsIntegrations {
// 		for integration, versions := range integrations {
// 			integrationPath := path.Join(m.config.IntegrationsDirName, namespace, integration)

// 			for _, version := range versions {
// 				logger := log.With().
// 					Str("namespace", namespace).
// 					Str("integration", integration).
// 					Str("version", version.SemVer()).
// 					Logger()

// 				// retrieve & validate the integration config
// 				integrationConfig, err := m.getIntegrationConfigFromPathAtRef(integrationPath)
// 				if err != nil {
// 					logger.Error().Err(err).Msg("Failed to retrieve integration config")
// 					validationFailed = true
// 					continue
// 				}
// 				if err := integrationConfig.Validate(); err != nil {
// 					logger.Err(err).Msg("Failed to validate integration config")
// 					validationFailed = true
// 					continue
// 				}
// 			}
// 		}
// 	}

// 	if validationFailed {
// 		return errors.New("one or more integrations failed validation")
// 	}
// 	return nil
// }
