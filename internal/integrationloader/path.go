package integrationloader

import (
	"fmt"
	"os"
	"path"

	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	"github.com/sensu/catalog-api/internal/types"
)

type PathLoader struct {
	path string
}

func NewPathLoader(path string) PathLoader {
	return PathLoader{
		path: path,
	}
}

func (l PathLoader) LoadConfig() (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	filePath := path.Join(l.path, defaultConfigName)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return integration, fmt.Errorf("error accessing %s: %w", filePath, err)
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

func (l PathLoader) LoadChangelog() (string, error) {
	// TODO(jk): support both .yaml & .yml extensions
	return l.getFileContentsAsString(defaultChangelogName)
}

func (l PathLoader) LoadImages() (Images, error) {
	return Images{}, nil
}

func (l PathLoader) LoadLogo() (string, error) {
	return l.getFileContentsAsString(defaultLogoName)
}

func (l PathLoader) LoadReadme() (string, error) {
	return l.getFileContentsAsString(defaultReadmeName)
}

func (l PathLoader) LoadResources() (string, error) {
	return l.getFileContentsAsString(defaultResourcesName)
}

func (l PathLoader) getFileContentsAsString(relativePath string) (string, error) {
	filePath := path.Join(l.path, relativePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error accessing %s: %w", filePath, err)
	}
	return string(b), nil
}
