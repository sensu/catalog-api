package catalogmanager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rs/zerolog/log"
)

func (m CatalogManager) ValidateCatalogDir() error {
	integrationsDir := path.Join(m.config.RepoDir, m.config.IntegrationsDirName)

	// get a list of namespaces & the integrations that belong to them from the
	// directory structure
	files, err := ioutil.ReadDir(m.IntegrationsDir())
	if err != nil {
		return fmt.Errorf("error retrieving integrations directory listing: %w", err)
	}

	nsIntegrations := map[string][]string{}
	for _, file := range files {
		if file.IsDir() {
			namespace := file.Name()
			namespaceDir := path.Join(integrationsDir, namespace)

			namespaceFiles, err := ioutil.ReadDir(namespaceDir)
			if err != nil {
				return fmt.Errorf("error retrieving integrations directory listing: %w", err)
			}

			for _, namespaceFile := range namespaceFiles {
				if namespaceFile.IsDir() {
					integration := namespaceFile.Name()
					nsIntegrations[namespace] = append(nsIntegrations[namespace], integration)
				}
			}
		}
	}

	// loop through the list of namespaces & integrations, and unmarshal the
	// configs & resource files
	validationFailed := false
	for namespace, integrations := range nsIntegrations {
		for _, integration := range integrations {
			integrationPath := path.Join(m.config.IntegrationsDirName, namespace, integration)
			loader := NewIntegrationPathLoader(integrationPath)

			logger := log.With().
				Str("namespace", namespace).
				Str("integration", integration).
				Logger()

			// retrieve & validate the integration config
			integrationConfig, err := loader.LoadConfig()
			if err != nil {
				logger.Err(err).Msg("Failed to retrieve integration config")
				validationFailed = true
				continue
			}
			if err := integrationConfig.Validate(); err != nil {
				logger.Err(err).Msg("Failed to validate integration config")
				validationFailed = true
				continue
			}

			// retrieve & validate sensu resources
			if _, err = m.getIntegrationResourcesFromPath(integrationPath); err != nil {
				logger.Err(err).Msg("Failed to retrieve resources file")
			}
			// TODO(jk): call resouces.Validate() once it's implemented

			// retrieve & validate logo

			// retrieve & validate readme
			// retrieve & validate changelog
			// retrieve & validate images
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
