package v1

import (
	"fmt"
	"path"

	catalogv1 "github.com/sensu/catalog-api/api/catalog/v1"
)

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version.json
type IntegrationVersionEndpoint struct {
	outputPath string
	data       catalogv1.Integration
}

type IntegrationVersion struct {
	Integration catalogv1.Integration `json:"integration" yaml:"integration"`
	Version     string                `json:"version" yaml:"version"`
}

func NewIntegrationVersionEndpoint(basePath string, iv IntegrationVersion) IntegrationVersionEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		iv.Integration.Metadata.Namespace,
		iv.Integration.Metadata.Name,
		fmt.Sprintf("%s.json", iv.Version))

	return IntegrationVersionEndpoint{
		outputPath: outputPath,
		data:       iv.Integration,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/versions.json
type IntegrationVersionsEndpoint struct {
	outputPath string
	data       []string
}

type IntegrationWithVersions struct {
	catalogv1.Integration
	Versions []string `json:"versions" yaml:"versions"`
}

func NewIntegrationVersionsEndpoint(basePath string, ivs IntegrationWithVersions) IntegrationVersionsEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		ivs.Metadata.Namespace,
		ivs.Metadata.Name,
		"versions.json")

	return IntegrationVersionsEndpoint{
		outputPath: outputPath,
		data:       ivs.Versions,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name.json
type IntegrationEndpoint struct {
	outputPath string
	data       IntegrationWithVersions
}

func NewIntegrationEndpoint(basePath string, ivs IntegrationWithVersions) IntegrationEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		ivs.Integration.Metadata.Namespace,
		fmt.Sprintf("%s.json", ivs.Integration.Metadata.Name))

	return IntegrationEndpoint{
		outputPath: outputPath,
		data:       ivs,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace.json
type IntegrationNamespaceEndpoint struct {
	outputPath string
	data       IntegrationNamespace
}

type IntegrationNamespace struct {
	Name         string   `json:"name" yaml:"name"`
	Integrations []string `json:"integrations" yaml:"integrations"`
}

func NewIntegrationNamespaceEndpoint(basePath string, ns IntegrationNamespace) IntegrationNamespaceEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		fmt.Sprintf("%s.json", ns.Name))

	return IntegrationNamespaceEndpoint{
		outputPath: outputPath,
		data:       ns,
	}
}

// GET /api/:generated_sha/v1/integrations.json
type IntegrationNamespacesEndpoint struct {
	outputPath string
	data       []string
}

func NewIntegrationNamespacesEndpoint(basePath string, ns []string) IntegrationNamespacesEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		"integrations.json")

	return IntegrationNamespacesEndpoint{
		outputPath: outputPath,
		data:       ns,
	}
}
