package catalogmanager

import (
	"fmt"
	"os"
	"path"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	"github.com/sensu/catalog-api/internal/types"
)

var (
	defaultIntegrationConfigName    = "sensu-integration.yaml"
	defaultIntegrationResourcesName = "sensu-resources.yaml"
	defaultIntegrationLogoName      = "logo.png"
	defaultIntegrationReadmeName    = "README.md"
	defaultIntegrationChangelogName = "CHANGELOG.md"
	defaultIntegrationImagesDirName = "img"
)

type IntegrationLoader interface {
	LoadConfig() (catalogv1.Integration, error)
	LoadResources() (string, error)
	LoadLogo() (string, error)
	LoadReadme() (string, error)
	LoadChangelog() (string, error)
	LoadImages() (Images, error)
}

type IntegrationGitLoader struct {
	repo *git.Repository
	ref  string
	path string
}

func NewIntegrationGitLoader(repo *git.Repository, ref string, path string) IntegrationGitLoader {
	return IntegrationGitLoader{
		repo: repo,
		ref:  ref,
		path: path,
	}
}

func (l IntegrationGitLoader) LoadConfig() (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	filePath := path.Join(l.path, defaultIntegrationConfigName)
	contents, err := l.getFileContents(filePath)
	if err != nil {
		return integration, fmt.Errorf("error reading contents of integration config: %w", err)
	}

	raw, err := types.RawWrapperFromYAMLBytes([]byte(contents))
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

func (l IntegrationGitLoader) getFileContents(path string) (string, error) {
	// attempt to resolve the git ref to a revision
	hash, err := l.repo.ResolveRevision(plumbing.Revision(l.ref))
	if err != nil {
		return "", fmt.Errorf("error resolving git revision: %w", err)
	}

	// attempt to retrieve the commit object for the hash
	commit, err := l.repo.CommitObject(*hash)
	if err != nil {
		return "", fmt.Errorf("error retrieving commit for hash: %w", err)
	}

	// attempt to retrieve the directory tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("error retrieving tree for commit: %w", err)
	}

	// attempt to load the file from the tree
	file, err := tree.File(path)
	if err != nil {
		return "", fmt.Errorf("error loading file from tree: %w", err)
	}

	// attempt to read the contents of the file
	contents, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("error reading contents of file: %w", err)
	}

	return contents, nil
}

type IntegrationPathLoader struct {
	path string
}

func NewIntegrationPathLoader(path string) IntegrationPathLoader {
	return IntegrationPathLoader{
		path: path,
	}
}

func (l IntegrationPathLoader) LoadConfig() (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	filePath := path.Join(l.path, defaultIntegrationConfigName)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return integration, fmt.Errorf("error reading contents of integration config: %w", err)
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
