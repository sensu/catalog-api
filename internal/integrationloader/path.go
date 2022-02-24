package integrationloader

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
)

type PathLoader struct {
	integrationPath string
}

func NewPathLoader(integrationPath string) PathLoader {
	return PathLoader{
		integrationPath: integrationPath,
	}
}

func (l PathLoader) LoadConfig() (catalogv1.Integration, error) {
	return loadConfig(l)
}

func (l PathLoader) LoadChangelog() (string, error) {
	return loadChangelog(l)
}

func (l PathLoader) LoadImages() (Images, error) {
	images := Images{}
	imagesPath := path.Join(l.integrationPath, defaultImagesDirName)

	files, err := ioutil.ReadDir(imagesPath)
	if _, ok := err.(*fs.PathError); ok {
		// no images found; skipping
		return images, nil
	} else if err != nil {
		return images, err
	}

	for _, f := range files {
		relativePath := path.Join(defaultImagesDirName, f.Name())

		if !f.IsDir() {
			match, _ := regexp.MatchString(reImageExtensions, f.Name())
			if match {
				data, err := l.GetFileContentsAsString(relativePath)
				if err != nil {
					return images, err
				}
				images[f.Name()] = data
			}
		}
	}

	return images, nil
}

func (l PathLoader) LoadLogo() (string, error) {
	return loadLogo(l)
}

func (l PathLoader) LoadReadme() (string, error) {
	return loadReadme(l)
}

func (l PathLoader) LoadResources() (string, error) {
	return loadResources(l)
}

func (l PathLoader) GetFileContentsAsBytes(relativePath string) ([]byte, error) {
	filePath := path.Join(l.integrationPath, relativePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (l PathLoader) GetFileContentsAsString(relativePath string) (string, error) {
	b, err := l.GetFileContentsAsBytes(relativePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
