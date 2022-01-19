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
