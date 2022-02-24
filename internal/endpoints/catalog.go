package endpoints

import (
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
)

// GET /api/:generated_sha/v1/catalog.json
func GenerateCatalogEndpoint(basePath string, nsIntegrations map[string][]catalogapiv1.IntegrationVersion) error {
	catalog := catalogapiv1.Catalog{
		NamespacedIntegrations: nsIntegrations,
	}
	endpoint := catalogapiv1.NewCatalogEndpoint(basePath, catalog)
	return renderJSON(endpoint)
}
