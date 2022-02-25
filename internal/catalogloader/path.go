package catalogloader

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
)

type PathLoader struct {
	repoPath         string
	integrationsPath string
}

func NewPathLoader(repoPath string, integrationsPath string) PathLoader {
	return PathLoader{
		repoPath:         repoPath,
		integrationsPath: integrationsPath,
	}
}

func (l PathLoader) IntegrationsPath() string {
	return l.integrationsPath
}

func (l PathLoader) NewIntegrationLoader(namespace string, integration string, version string) integrationloader.Loader {
	integrationPath := path.Join(l.IntegrationsPath(), namespace, integration)
	return integrationloader.NewPathLoader(integrationPath)
}

func (l PathLoader) LoadIntegrations() (types.Integrations, error) {
	integrations := types.Integrations{}

	// get a list of namespaces & the integrations that belong to them from the
	// directory structure
	files, err := ioutil.ReadDir(l.IntegrationsPath())
	if err != nil {
		return integrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			namespace := file.Name()
			namespaceDir := path.Join(l.IntegrationsPath(), namespace)

			namespaceFiles, err := ioutil.ReadDir(namespaceDir)
			if err != nil {
				return integrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
			}

			for _, namespaceFile := range namespaceFiles {
				if namespaceFile.IsDir() {
					integration := types.IntegrationVersion{
						Name:          namespaceFile.Name(),
						Namespace:     namespace,
						Major:         0,
						Minor:         0,
						Patch:         0,
						Prerelease:    "",
						BuildMetadata: "",
						GitTag:        "",
						GitRef:        "",
					}
					integrations = append(integrations, integration)
				}
			}
		}
	}

	return integrations, nil
}
