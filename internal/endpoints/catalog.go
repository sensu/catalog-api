package endpoints

import (
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/types"
)

// GET /api/:generated_sha/v1/catalog.json
func GenerateCatalogEndpoint(basePath string, nis types.NamespacedIntegrations) error {
	nsIntegrations := map[string][]string{}
	for namespace, vis := range nis {
		integrations := []string{}
		for integration, _ := range vis {
			integrations = append(integrations, integration)
		}
		nsIntegrations[namespace] = integrations
	}
	catalog := catalogapiv1.Catalog{
		NamespacedIntegrations: nsIntegrations,
	}
	endpoint := catalogapiv1.NewCatalogEndpoint(basePath, catalog)
	return renderJSON(endpoint)
}
