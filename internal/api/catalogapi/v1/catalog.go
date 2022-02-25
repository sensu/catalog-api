package catalogapiv1

import (
	"path"
)

// GET /api/:generated_sha/v1/catalog.json
type CatalogEndpoint struct {
	outputPath string
	data       Catalog
}

func (e CatalogEndpoint) GetOutputPath() string { return e.outputPath }
func (e CatalogEndpoint) GetData() interface{}  { return e.data }

type Catalog struct {
	NamespacedIntegrations map[string][]IntegrationVersion `json:"namespaced_integrations" yaml:"namespaced_integrations"`
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
