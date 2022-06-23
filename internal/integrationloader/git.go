package integrationloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
)

type GitLoader struct {
	repo            *git.Repository
	ref             string
	integrationPath string
}

func NewGitLoader(repo *git.Repository, ref string, integrationPath string) GitLoader {
	return GitLoader{
		repo:            repo,
		ref:             ref,
		integrationPath: integrationPath,
	}
}

func (l GitLoader) LoadConfig() (catalogv1.Integration, error) {
	return loadConfig(l)
}

func (l GitLoader) LoadChangelog() (string, error) {
	return loadChangelog(l)
}

func (l GitLoader) LoadImages() (Images, error) {
	images := Images{}
	imagesPath := path.Join(l.integrationPath, defaultImagesDirName)

	// attempt to resolve the git ref to a revision
	hash, err := l.repo.ResolveRevision(plumbing.Revision(l.ref))
	if err != nil {
		return images, fmt.Errorf("error resolving git revision %s: %w", l.ref, err)
	}

	// attempt to retrieve the commit object for the hash
	commit, err := l.repo.CommitObject(*hash)
	if err != nil {
		return images, fmt.Errorf("error retrieving commit for hash %s: %w", hash, err)
	}

	// attempt to retrieve the directory tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return images, fmt.Errorf("error retrieving tree for commit %s: %w", hash, err)
	}

	// attempt to load the file from the tree
	// file, err := tree.FindEntry(imagesPath)
	// if err != nil {
	// 	if errors.Is(err, plumbingobject.ErrEntryNotFound) {
	// 		return images, nil
	// 	}
	// 	return images, fmt.Errorf("error accessing %s for hash %s: %w", imagesPath, hash, err)
	// }

	imagesTree, err := tree.Tree(imagesPath)
	if err != nil {
		if errors.Is(err, object.ErrDirectoryNotFound) {
			return images, nil
		}
		return images, err
	}

	if err := imagesTree.Files().ForEach(func(f *object.File) error {
		match, _ := regexp.MatchString(reImageExtensions, f.Name)
		if match {
			data, err := f.Contents()
			if err != nil {
				return err
			}
			images[f.Name] = data
		}
		return nil
	}); err != nil {
		return images, err
	}

	return images, nil
}

func (l GitLoader) LoadDashboards() (Dashboards, error) {
	dashboards := Dashboards{}
	dashboardsPath := path.Join(l.integrationPath, defaultDashboardsDirName)

	// attempt to resolve the git ref to a revision
	hash, err := l.repo.ResolveRevision(plumbing.Revision(l.ref))
	if err != nil {
		return dashboards, fmt.Errorf("error resolving git revision %s: %w", l.ref, err)
	}

	// attempt to retrieve the commit object for the hash
	commit, err := l.repo.CommitObject(*hash)
	if err != nil {
		return dashboards, fmt.Errorf("error retrieving commit for hash %s: %w", hash, err)
	}

	// attempt to retrieve the directory tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return dashboards, fmt.Errorf("error retrieving tree for commit %s: %w", hash, err)
	}

	dashboardsTree, err := tree.Tree(dashboardsPath)
	if err != nil {
		if errors.Is(err, object.ErrDirectoryNotFound) {
			return dashboards, nil
		}
		return dashboards, err
	}

	if err := dashboardsTree.Files().ForEach(func(f *object.File) error {
		match, _ := regexp.MatchString(reDashboardExtensions, f.Name)
		if match {
			data, err := f.Contents()
			if err != nil {
				return err
			}
			dashboardPath := path.Join(dashboardsPath, f.Name)
			var anyMap map[string]interface{}
			if err := json.Unmarshal([]byte(data), &anyMap); err != nil {
				return fmt.Errorf("error unmarshaling dashboard: %w (%s)", err, dashboardPath)
			}
			dashboards[f.Name] = data
		}
		return nil
	}); err != nil {
		return dashboards, err
	}

	return dashboards, nil
}

func (l GitLoader) LoadLogo() (string, error) {
	return loadLogo(l)
}

func (l GitLoader) LoadReadme() (string, error) {
	return loadReadme(l)
}

func (l GitLoader) LoadResources() (string, error) {
	return loadResources(l)
}

func (l GitLoader) GetFileContentsAsBytes(relativePath string) ([]byte, error) {
	contents, err := l.GetFileContentsAsString(relativePath)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

func (l GitLoader) GetFileContentsAsString(relativePath string) (string, error) {
	filePath := path.Join(l.integrationPath, relativePath)

	// attempt to resolve the git ref to a revision
	hash, err := l.repo.ResolveRevision(plumbing.Revision(l.ref))
	if err != nil {
		return "", fmt.Errorf("error resolving git revision %s: %w", l.ref, err)
	}

	// attempt to retrieve the commit object for the hash
	commit, err := l.repo.CommitObject(*hash)
	if err != nil {
		return "", fmt.Errorf("error retrieving commit for hash %s: %w", hash, err)
	}

	// attempt to retrieve the directory tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("error retrieving tree for commit %s: %w", hash, err)
	}

	// attempt to load the file from the tree
	file, err := tree.File(filePath)
	if err != nil {
		return "", fmt.Errorf("error accessing %s for hash %s: %w", filePath, hash, err)
	}

	// attempt to read the contents of the file
	contents, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("error reading contents of %s for hash %s: %w", filePath, hash, err)
	}

	return contents, nil
}
