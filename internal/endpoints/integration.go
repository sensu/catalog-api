package endpoints

import (
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	catalogapiv1 "github.com/sensu/catalog-api/internal/api/catalogapi/v1"
	"github.com/sensu/catalog-api/internal/types"
)

// GET /api/:generated_sha/v1/integrations/:namespace/:name.json
func GenerateIntegrationEndpoint(basePath string, integration catalogv1.Integration, ivs []types.IntegrationVersion) error {
	versions := []string{}
	for _, iv := range ivs {
		versions = append(versions, iv.SemVer())
	}
	integrationWithVersions := catalogapiv1.IntegrationWithVersions{
		Integration: integration,
		Versions:    versions,
	}
	endpoint := catalogapiv1.NewIntegrationEndpoint(basePath, integrationWithVersions)
	return renderJSON(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/versions.json
func GenerateIntegrationVersionsEndpoint(basePath string, namespace string, integration string, ivs []types.IntegrationVersion) error {
	versions := []string{}
	for _, iv := range ivs {
		versions = append(versions, iv.SemVer())
	}
	endpoint := catalogapiv1.NewIntegrationVersionsEndpoint(basePath, namespace, integration, versions)
	return renderJSON(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version.json
func GenerateIntegrationVersionEndpoint(basePath string, integration catalogv1.Integration, version types.IntegrationVersion) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionEndpoint(basePath, iv)
	return renderJSON(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/resources.json
func GenerateIntegrationVersionResourcesEndpoint(basePath string, integration catalogv1.Integration, version types.IntegrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionResourcesEndpoint(basePath, iv, data)
	return renderJSON(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/logo.png
func GenerateIntegrationVersionLogoEndpoint(basePath string, integration catalogv1.Integration, version types.IntegrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionLogoEndpoint(basePath, iv, data)
	return renderRaw(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/README.md
func GenerateIntegrationVersionReadmeEndpoint(basePath string, integration catalogv1.Integration, version types.IntegrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionReadmeEndpoint(basePath, iv, data)
	return renderRaw(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/CHANGELOG.md
func GenerateIntegrationVersionChangelogEndpoint(basePath string, integration catalogv1.Integration, version types.IntegrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionChangelogEndpoint(basePath, iv, data)
	return renderRaw(endpoint)
}
