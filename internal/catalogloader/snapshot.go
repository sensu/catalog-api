package catalogloader

import (
	git "github.com/go-git/go-git/v5"
	"github.com/sensu/catalog-api/internal/integrationloader"
	"github.com/sensu/catalog-api/internal/types"
)

type SnapshotLoader struct {
	gitLoader  GitLoader
	pathLoader PathLoader
}

func NewSnapshotLoader(repo *git.Repository, repoPath string, integrationsDirName string) SnapshotLoader {
	return SnapshotLoader{
		gitLoader:  NewGitLoader(repo, integrationsDirName),
		pathLoader: NewPathLoader(repoPath, integrationsDirName),
	}
}

func (l SnapshotLoader) NewIntegrationLoader(integration types.IntegrationVersion) integrationloader.Loader {
	switch integration.Source {
	case "git":
		return l.gitLoader.NewIntegrationLoader(integration)
	case "path":
		return l.pathLoader.NewIntegrationLoader(integration)
	}
	panic("invalid source field in integration")
}

func (l SnapshotLoader) LoadIntegrations() (types.Integrations, error) {
	integrations := types.Integrations{}

	gitIntegrations, err := l.gitLoader.LoadIntegrations()
	if err != nil {
		return integrations, err
	}
	integrations = append(integrations, gitIntegrations...)

	pathIntegrations, err := l.pathLoader.LoadIntegrations()
	if err != nil {
		return integrations, err
	}
	integrations = append(integrations, pathIntegrations...)

	return integrations, nil
}
