package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	catalogv1 "github.com/sensu/catalog-api/api/catalog/v1"
	catalogapiv1 "github.com/sensu/catalog-api/api/catalogapi/v1"
)

type apiEndpoint interface {
	GetOutputPath() string
	GetData() interface{}
}

func generateJSONEndpoint(endpoint apiEndpoint) error {
	contents, err := json.Marshal(endpoint.GetData())
	if err != nil {
		return fmt.Errorf("error generating endpoint: %w", err)
	}

	outputPath := endpoint.GetOutputPath()

	// ensure the parent directory exists
	parent := filepath.Dir(outputPath)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return fmt.Errorf("error creating endpoint parent directory: %w", err)
	}

	// write the endpoint contents to the output path
	if err := os.WriteFile(outputPath, contents, 0600); err != nil {
		return fmt.Errorf("error creating endpoint file: %w", err)
	}

	return nil
}

func generateEndpoint(endpoint apiEndpoint) error {
	contents, ok := endpoint.GetData().(string)
	if !ok {
		return fmt.Errorf("endpoint data is not a string")
	}

	outputPath := endpoint.GetOutputPath()

	// ensure the parent directory exists
	parent := filepath.Dir(outputPath)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return fmt.Errorf("error creating endpoint parent directory: %w", err)
	}

	// write the endpoint contents to the output path
	if err := os.WriteFile(outputPath, []byte(contents), 0600); err != nil {
		return fmt.Errorf("error creating endpoint file: %w", err)
	}

	return nil
}

// GET /api/:generated_sha/v1/catalog.json
func generateCatalogEndpoint(basePath string, nis namespacedIntegrations) error {
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
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace.json
func generateIntegrationNamespaceEndpoint(basePath string, namespace string, ivs []integrationVersion) error {
	integrations := []string{}
	for _, iv := range ivs {
		integrations = append(integrations, iv.Name)
	}
	ns := catalogapiv1.IntegrationNamespace{
		Name:         namespace,
		Integrations: integrations,
	}
	endpoint := catalogapiv1.NewIntegrationNamespaceEndpoint(basePath, ns)
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name.json
func generateIntegrationEndpoint(basePath string, integration catalogv1.Integration, ivs []integrationVersion) error {
	versions := []string{}
	for _, iv := range ivs {
		versions = append(versions, iv.SemVer())
	}
	integrationWithVersions := catalogapiv1.IntegrationWithVersions{
		Integration: integration,
		Versions:    versions,
	}
	endpoint := catalogapiv1.NewIntegrationEndpoint(basePath, integrationWithVersions)
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/versions.json
func generateIntegrationVersionsEndpoint(basePath string, namespace string, integration string, ivs []integrationVersion) error {
	versions := []string{}
	for _, iv := range ivs {
		versions = append(versions, iv.SemVer())
	}
	endpoint := catalogapiv1.NewIntegrationVersionsEndpoint(basePath, namespace, integration, versions)
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version.json
func generateIntegrationVersionEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionEndpoint(basePath, iv)
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/resources.json
func generateIntegrationVersionResourcesEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionResourcesEndpoint(basePath, iv, data)
	return generateJSONEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/logo.png
func generateIntegrationVersionLogoEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionLogoEndpoint(basePath, iv, data)
	return generateEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/README.md
func generateIntegrationVersionReadmeEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionReadmeEndpoint(basePath, iv, data)
	return generateEndpoint(endpoint)
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/CHANGELOG.md
func generateIntegrationVersionChangelogEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion, data string) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionChangelogEndpoint(basePath, iv, data)
	return generateEndpoint(endpoint)
}

// GET /api/version.json
func generateVersionEndpoint(basePath string, sha256 string) error {
	version := catalogapiv1.ReleaseVersion{
		ReleaseSHA256: sha256,
		LastUpdated:   time.Now().Unix(),
	}
	endpoint := catalogapiv1.NewVersionEndpoint(basePath, version)
	return generateJSONEndpoint(endpoint)
}
