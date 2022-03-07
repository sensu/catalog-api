package catalogloader

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
)

var sourcePath = "path"

type PathLoader struct {
	repoPath            string
	integrationsDirName string
}

func NewPathLoader(repoPath string, integrationsDirName string) PathLoader {
	return PathLoader{
		repoPath:            repoPath,
		integrationsDirName: integrationsDirName,
	}
}

func (l PathLoader) IntegrationsAbsPath() string {
	return path.Join(l.repoPath, l.integrationsDirName)
}

func (l PathLoader) NewIntegrationLoader(integration types.IntegrationVersion) integrationloader.Loader {
	integrationPath := integration.Path(l.IntegrationsAbsPath())
	return integrationloader.NewPathLoader(integrationPath)
}

func (l PathLoader) LoadIntegrations() (types.Integrations, error) {
	integrations := types.Integrations{}

	// get a list of namespaces & the integrations that belong to them from the
	// directory structure
	files, err := ioutil.ReadDir(l.IntegrationsAbsPath())
	if err != nil {
		return integrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			namespace := file.Name()
			namespaceDir := path.Join(l.IntegrationsAbsPath(), namespace)

			namespaceFiles, err := ioutil.ReadDir(namespaceDir)
			if err != nil {
				return integrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
			}

			for _, namespaceFile := range namespaceFiles {
				if namespaceFile.IsDir() {
					integration := types.IntegrationVersion{
						Name:          namespaceFile.Name(),
						Namespace:     namespace,
						Major:         99991231,
						Minor:         0,
						Patch:         0,
						Prerelease:    "",
						BuildMetadata: "",
						GitTag:        "",
						GitRef:        "",
						Source:        sourcePath,
					}
					integrations = append(integrations, integration)

					log.Info().
						Str("name", integration.Name).
						Str("namespace", integration.Namespace).
						Str("version", integration.SemVer()).
						Str("source", sourcePath).
						Msg("Found integration version")
				}
			}
		}
	}

	return integrations, nil
}
