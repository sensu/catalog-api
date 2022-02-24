package integrationloader

import (
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
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
}
