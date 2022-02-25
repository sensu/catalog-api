package catalogloader

import (
	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
)

type Loader interface {
	LoadIntegrations() (types.Integrations, error)
	NewIntegrationLoader(namespace string, integration string, version string) integrationloader.Loader
}
