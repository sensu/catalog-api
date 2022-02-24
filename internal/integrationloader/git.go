package integrationloader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"path"
	"regexp"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	"github.com/sensu/catalog-api/internal/types"
	"gopkg.in/yaml.v3"
)

type GitLoader struct {
	repo *git.Repository
	ref  string
	path string
}

func NewGitLoader(repo *git.Repository, ref string, path string) GitLoader {
	return GitLoader{
		repo: repo,
		ref:  ref,
		path: path,
	}
}

func (l GitLoader) LoadConfig() (catalogv1.Integration, error) {
	var integration catalogv1.Integration

	// TODO(jk): support both .yaml & .yml extensions
	contents, err := l.getFileContentsAsBytes(defaultConfigName)
	if err != nil {
		return integration, err
	}

	raw, err := types.RawWrapperFromYAMLBytes(contents)
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

func (l GitLoader) LoadChangelog() (string, error) {
	return l.getFileContentsAsString(defaultChangelogName)
}

func (l GitLoader) LoadImages() (Images, error) {
	images := Images{}

	imagesPath := path.Join(l.path, defaultImagesDirName)
	files, err := ioutil.ReadDir(imagesPath)
	if _, ok := err.(*fs.PathError); ok {
		// no images found; skipping
		return images, nil
	} else if err != nil {
		return images, fmt.Errorf("error reading directory %s: %w", imagesPath, err)
	}

	for _, f := range files {
		relativePath := path.Join(defaultImagesDirName, f.Name())

		if !f.IsDir() {
			match, _ := regexp.MatchString(`.*\.(jpg|gif|png)$`, f.Name())
			if match {
				data, err := l.getFileContentsAsString(relativePath)
				if err != nil {
					return images, err
				}
				images[f.Name()] = data
			}
		}
	}
	return Images{}, nil
}

func (l GitLoader) LoadLogo() (string, error) {
	// TODO(jk): add basic validation using http.DetectContentType
	// https://pkg.go.dev/net/http#DetectContentType
	return l.getFileContentsAsString(defaultLogoName)
}

func (l GitLoader) LoadReadme() (string, error) {
	return l.getFileContentsAsString(defaultReadmeName)
}

func (l GitLoader) LoadResources() (string, error) {
	// TODO(jk): support both .yaml & .yml extensions
	b, err := l.getFileContentsAsBytes(defaultResourcesName)
	if err != nil {
		return "", err
	}

	// attempt to unmarshal yaml to verify that the yaml is valid
	// TODO(jk): iterate through & validate each resource against the supported
	// versions of Sensu that the integration defines
	resources := []map[string]interface{}{}
	dec := yaml.NewDecoder(bytes.NewReader(b))
	for {
		doc := new(map[string]interface{})

		if err := dec.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("error parsing %s: %w", defaultResourcesName, err)
		}
		if doc == nil {
			return "", fmt.Errorf("error parsing %s", defaultResourcesName)
		}
		resources = append(resources, *doc)
	}

	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		return "", fmt.Errorf("error json marshalling after unmarshalling %s: %w", defaultResourcesName, err)
	}

	return string(resourcesJSON), nil
}

func (l GitLoader) getFileContentsAsBytes(relativePath string) ([]byte, error) {
	contents, err := l.getFileContentsAsString(relativePath)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

func (l GitLoader) getFileContentsAsString(relativePath string) (string, error) {
	filePath := path.Join(l.path, relativePath)

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
