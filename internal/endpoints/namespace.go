package endpoints

import (
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/types"
)

// GET /api/:generated_sha/v1/integrations/:namespace.json
func generateIntegrationNamespaceEndpoint(basePath string, namespace string, ivs []types.IntegrationVersion) error {
	integrations := []string{}
	for _, iv := range ivs {
		integrations = append(integrations, iv.Name)
	}
	ns := catalogapiv1.IntegrationNamespace{
		Name:         namespace,
		Integrations: integrations,
	}
	endpoint := catalogapiv1.NewIntegrationNamespaceEndpoint(basePath, ns)
	return renderJSON(endpoint)
}
