package catalogloader

import (
	git "github.com/go-git/go-git/v5"
	"github.com/sensu/catalog-api/internal/integrationloader"
)

type GitLoader struct {
	repo *git.Repository
}

func NewGitLoader(repo *git.Repository) GitLoader {
	return GitLoader{
		repo: repo,
	}
}

func (l GitLoader) LoadIntegrations() (map[string][]string, error) {
	return map[string][]string{}, nil
}

func (l GitLoader) IntegrationLoader(namespace string, integration string, version string) integrationloader.Loader {
	ref := ""
	path := ""
	return integrationloader.NewGitLoader(l.repo, ref, path)
}
