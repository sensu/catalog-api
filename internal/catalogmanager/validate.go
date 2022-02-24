package catalogmanager

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (m CatalogManager) ValidateCatalog() error {
	integrations, err := m.loader.LoadIntegrations()
	if err != nil {
		return fmt.Errorf("error loading integrations from catalog: %w", err)
	}

	// loop through the list of namespaces & integrations, and unmarshal the
	// configs & resource files
	validationFailed := false
	for namespace, integrations := range integrations.ByNamespace() {
		for _, integration := range integrations {
			integrationLoader := m.loader.NewIntegrationLoader(namespace, integration.Name, "")

			logger := log.With().
				Str("namespace", namespace).
				Str("integration", integration.Name).
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
			// TODO(jk): call resouces.Validate() once it's implemented
			if _, err = integrationLoader.LoadResources(); err != nil {
				logger.Err(err).Msg("Failed to load resources file")
			}

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
