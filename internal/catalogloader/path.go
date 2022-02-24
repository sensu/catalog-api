package catalogloader

import (
	"fmt"
	"io/ioutil"
	"path"
)

type PathLoader struct {
	path string
}

func NewPathLoader(path string) PathLoader {
	return PathLoader{
		path: path,
	}
}

func (l PathLoader) LoadIntegrations() (map[string][]string, error) {
	nsIntegrations := map[string][]string{}

	// get a list of namespaces & the integrations that belong to them from the
	// directory structure
	files, err := ioutil.ReadDir(l.path)
	if err != nil {
		return nsIntegrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			namespace := file.Name()
			namespaceDir := path.Join(l.path, namespace)

			namespaceFiles, err := ioutil.ReadDir(namespaceDir)
			if err != nil {
				return nsIntegrations, fmt.Errorf("error retrieving integrations directory listing: %w", err)
			}

			for _, namespaceFile := range namespaceFiles {
				if namespaceFile.IsDir() {
					integration := namespaceFile.Name()
					nsIntegrations[namespace] = append(nsIntegrations[namespace], integration)
				}
			}
		}
	}

	return nsIntegrations, nil
}
