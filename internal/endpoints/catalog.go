package endpoints

import (
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
)

// GET /api/:generated_sha/v1/catalog.json
func GenerateCatalogEndpoint(basePath string, nis map[string][]catalogapiv1.IntegrationVersion) error {
	// TODO(jk): Implement a less hacky way of determing what fields are
	// rendered for an endpoint. Perhaps using gotemplates for views.
	//
	// Set prompts & resource_patches fields to empty strings to prevent them
	// from being shown in the catalog endpoint.
	for namespace, versions := range nis {
		for i, version := range versions {
			version.Prompts = nil
			version.ResourcePatches = nil

			nis[namespace][i] = version
		}
	}

	catalog := catalogapiv1.Catalog{
		NamespacedIntegrations: nis,
	}
	endpoint := catalogapiv1.NewCatalogEndpoint(basePath, catalog)
	return renderJSON(endpoint)
}
