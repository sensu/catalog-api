package integrationloader

import (
	"encoding/json"
	"fmt"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	"github.com/sensu/catalog-api/internal/types"
)

var (
	defaultConfigName    = "sensu-integration.yaml"
	defaultResourcesName = "sensu-resources.yaml"
	defaultLogoName      = "logo.png"
	defaultReadmeName    = "README.md"
	defaultChangelogName = "CHANGELOG.md"
	defaultImagesDirName = "img"
)

type Loader interface {
	LoadConfig() (catalogv1.Integration, error)
	LoadChangelog() (string, error)
	LoadImages() (Images, error)
	LoadLogo() (string, error)
	LoadReadme() (string, error)
	LoadResources() (string, error)
	GetFileContentsAsBytes(string) ([]byte, error)
	GetFileContentsAsString(string) (string, error)
}

func loadConfig(loader Loader) (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	b, err := loader.GetFileContentsAsBytes(defaultConfigName)
	if err != nil {
		return integration, err
	}

	raw, err := types.RawWrapperFromYAMLBytes(b)
	if err != nil {
		return integration, err
	}

	wrap, err := types.WrapperFromRawWrapper(raw)
	if err != nil {
		return integration, err
	}
	integration = wrap.Value.(catalogv1.Integration)

	return integration, nil
}

func loadChangelog(loader Loader) (string, error) {
	// TODO(jk): support both .yaml & .yml extensions
	return loader.GetFileContentsAsString(defaultChangelogName)
}

func loadLogo(loader Loader) (string, error) {
	// TODO(jk): add basic validation using http.DetectContentType
	// https://pkg.go.dev/net/http#DetectContentType
	return loader.GetFileContentsAsString(defaultLogoName)
}

func loadReadme(loader Loader) (string, error) {
	return loader.GetFileContentsAsString(defaultReadmeName)
}

func loadResources(loader Loader) (string, error) {
	b, err := loader.GetFileContentsAsBytes(defaultResourcesName)
	if err != nil {
		return "", err
	}

	// attempt to unmarshal yaml to verify that the yaml is valid
	// TODO(jk): iterate through & validate each resource against the supported
	// versions of Sensu that the integration defines
	resources, err := catalogv1.ResourcesFromYAML(b)
	if err != nil {
		return "", fmt.Errorf("error parsing %s: %w", defaultResourcesName, err)
	}

	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		return "", fmt.Errorf("error json marshalling after unmarshalling %s: %w", defaultResourcesName, err)
	}

	return string(resourcesJSON), nil
}
