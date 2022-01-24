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

func generateEndpoint(endpoint apiEndpoint) error {
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
	return generateEndpoint(endpoint)
}

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
	return generateEndpoint(endpoint)
}

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
	return generateEndpoint(endpoint)
}

func generateIntegrationVersionsEndpoint(basePath string, namespace string, integration string, ivs []integrationVersion) error {
	versions := []string{}
	for _, iv := range ivs {
		versions = append(versions, iv.SemVer())
	}
	endpoint := catalogapiv1.NewIntegrationVersionsEndpoint(basePath, namespace, integration, versions)
	return generateEndpoint(endpoint)
}

func generateIntegrationVersionEndpoint(basePath string, integration catalogv1.Integration, version integrationVersion) error {
	iv := catalogapiv1.IntegrationVersion{
		Integration: integration,
		Version:     version.SemVer(),
	}
	endpoint := catalogapiv1.NewIntegrationVersionEndpoint(basePath, iv)
	return generateEndpoint(endpoint)
}

func generateVersionEndpoint(basePath string, sha256 string) error {
	version := catalogapiv1.ReleaseVersion{
		ReleaseSHA256: sha256,
		LastUpdated:   time.Now().Unix(),
	}
	endpoint := catalogapiv1.NewVersionEndpoint(basePath, version)
	return generateEndpoint(endpoint)
}
