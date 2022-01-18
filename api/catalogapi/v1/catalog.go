package v1

import (
	"path"
)

// GET /api/:generated_sha/v1/catalog.json
type CatalogEndpoint struct {
	outputPath string
	data       Catalog
}

type Catalog struct {
	Name         string                 `json:"name" yaml:"name"`
	Integrations []IntegrationNamespace `json:"integrations" yaml:"integrations"`
}

func NewCatalogEndpoint(basePath string, catalog Catalog) CatalogEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		"catalog.json")

	return CatalogEndpoint{
		outputPath: outputPath,
		data:       catalog,
	}
}
