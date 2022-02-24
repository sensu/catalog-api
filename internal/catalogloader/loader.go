package catalogloader

import "github.com/sensu/catalog-api/internal/integrationloader"

type Loader interface {
	LoadIntegrations() (map[string][]string, error)
	IntegrationLoader(namespace string, integration string) integrationloader.Loader
}
