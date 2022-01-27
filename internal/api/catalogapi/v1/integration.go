package v1

import (
	"fmt"
	"path"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
)

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/resources.json
type IntegrationVersionResourcesEndpoint struct {
	outputPath string
	data       string
}

func (e IntegrationVersionResourcesEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionResourcesEndpoint) GetData() interface{}  { return e.data }

func NewIntegrationVersionResourcesEndpoint(basePath string, iv IntegrationVersion, data string) IntegrationVersionResourcesEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		iv.Integration.Metadata.Namespace,
		iv.Integration.Metadata.Name,
		iv.Version,
		"sensu-resources.json")

	return IntegrationVersionResourcesEndpoint{
		outputPath: outputPath,
		data:       data,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/logo.png
type IntegrationVersionLogoEndpoint struct {
	outputPath string
	data       string
}

func (e IntegrationVersionLogoEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionLogoEndpoint) GetData() interface{}  { return e.data }

func NewIntegrationVersionLogoEndpoint(basePath string, iv IntegrationVersion, data string) IntegrationVersionLogoEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		iv.Integration.Metadata.Namespace,
		iv.Integration.Metadata.Name,
		iv.Version,
		"logo.png")

	return IntegrationVersionLogoEndpoint{
		outputPath: outputPath,
		data:       data,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/README.md
type IntegrationVersionReadmeEndpoint struct {
	outputPath string
	data       string
}

func (e IntegrationVersionReadmeEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionReadmeEndpoint) GetData() interface{}  { return e.data }

func NewIntegrationVersionReadmeEndpoint(basePath string, iv IntegrationVersion, data string) IntegrationVersionReadmeEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		iv.Integration.Metadata.Namespace,
		iv.Integration.Metadata.Name,
		iv.Version,
		"README.md")

	return IntegrationVersionReadmeEndpoint{
		outputPath: outputPath,
		data:       data,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version/CHANGELOG.md
type IntegrationVersionChangelogEndpoint struct {
	outputPath string
	data       string
}

func (e IntegrationVersionChangelogEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionChangelogEndpoint) GetData() interface{}  { return e.data }

func NewIntegrationVersionChangelogEndpoint(basePath string, iv IntegrationVersion, data string) IntegrationVersionChangelogEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		iv.Integration.Metadata.Namespace,
		iv.Integration.Metadata.Name,
		iv.Version,
		"CHANGELOG.md")

	return IntegrationVersionChangelogEndpoint{
		outputPath: outputPath,
		data:       data,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/:version.json
type IntegrationVersionEndpoint struct {
	outputPath string
	data       IntegrationVersion
}

func (e IntegrationVersionEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionEndpoint) GetData() interface{}  { return e.data }

type IntegrationVersion struct {
	catalogv1.Integration
	Version string `json:"version" yaml:"version"`
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
		data:       iv,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name/versions.json
type IntegrationVersionsEndpoint struct {
	outputPath string
	data       []string
}

func (e IntegrationVersionsEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationVersionsEndpoint) GetData() interface{}  { return e.data }

type IntegrationVersions []string

func NewIntegrationVersionsEndpoint(basePath string, namespace string, integration string, versions IntegrationVersions) IntegrationVersionsEndpoint {
	outputPath := path.Join(
		basePath,
		apiVersion,
		namespace,
		integration,
		"versions.json")

	return IntegrationVersionsEndpoint{
		outputPath: outputPath,
		data:       versions,
	}
}

// GET /api/:generated_sha/v1/integrations/:namespace/:name.json
type IntegrationEndpoint struct {
	outputPath string
	data       IntegrationWithVersions
}

func (e IntegrationEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationEndpoint) GetData() interface{}  { return e.data }

type IntegrationWithVersions struct {
	catalogv1.Integration
	Versions []string `json:"versions" yaml:"versions"`
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

func (e IntegrationNamespaceEndpoint) GetOutputPath() string { return e.outputPath }
func (e IntegrationNamespaceEndpoint) GetData() interface{}  { return e.data }

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
