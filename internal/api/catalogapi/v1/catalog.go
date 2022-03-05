package catalogapiv1

import (
	"errors"
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

func (c Catalog) GetIntegration(namespace, name, version string) (IntegrationVersion, error) {
	nsIntegrations, ok := c.NamespacedIntegrations[namespace]
	if !ok {
		return IntegrationVersion{}, errors.New("namespace not found in catalog")
	}

	for _, integration := range nsIntegrations {
		if integration.Metadata.Name == name && integration.Version == version {
			return integration, nil
		}
	}
	return IntegrationVersion{}, errors.New("integration not found in catalog")
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
