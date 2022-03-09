package catalogcmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/sensu/catalog-api/internal/catalogloader"
	"github.com/sensu/catalog-api/internal/catalogmanager"
)

var (
	throttleDelay = 1250 * time.Millisecond
)

type tmpCatalogManager struct {
	catalogmanager.CatalogManager
	tmpdir string
}

func (c *Config) newCatalogManager(loader catalogloader.Loader) (cm tmpCatalogManager, err error) {
	// create a temporay directory within c.tempDir to hold the generated api
	// files
	cm.tmpdir, err = os.MkdirTemp(c.tempDir, "")
	if err != nil {
		return cm, fmt.Errorf("error creating temp directory: %w", err)
	}
	cm.tmpdir, err = filepath.Abs(cm.tmpdir)
	if err != nil {
		return cm, err
	}

	// create a staging dir to hold the generated api files used to calculate
	// the checksum of the release
	stagingDir := path.Join(cm.tmpdir, "staging")
	if err := os.Mkdir(stagingDir, 0700); err != nil {
		return cm, fmt.Errorf("error creating staging directory: %w", err)
	}

	// create a release dir to hold the complete set of generated api files
	releaseDir := path.Join(cm.tmpdir, "release")
	if err := os.Mkdir(releaseDir, 0700); err != nil {
		return cm, fmt.Errorf("error creating release directory: %w", err)
	}

	mCfg := catalogmanager.Config{
		StagingDir: stagingDir,
		ReleaseDir: releaseDir,
	}

	// create a new catalog manager which is used to determine versions from git
	// tags, unmarshal resources, and generate the api
	cm.CatalogManager, err = catalogmanager.New(mCfg, loader)
	return cm, err
}

func (c *Config) newCatalogManagerFromRepo(ctx context.Context) (cm tmpCatalogManager, err error) {
	repo, err := git.PlainOpen(c.repoDir)
	if err != nil {
		return cm, err
	}

	var loader catalogloader.Loader
	if c.snapshot {
		loader = catalogloader.NewSnapshotLoader(repo, c.repoDir, c.integrationsDirName)
	} else {
		loader = catalogloader.NewGitLoader(repo, c.integrationsDirName)
	}
	return c.newCatalogManager(loader)
}

func (c *Config) generate(ctx context.Context) (string, error) {
	cm, err := c.newCatalogManagerFromRepo(ctx)
	if err != nil {
		return cm.tmpdir, err
	}

	// process the catalog & all its integrations
	err = cm.ProcessCatalog()
	return cm.tmpdir, err
}

func (c *Config) createWatcher(ctx context.Context) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, fmt.Errorf("error configuring watcher: %w", err)
	}
	dir := filepath.Join(c.repoDir, "integrations")
	walker := func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}
		return nil
	}
	if err := filepath.Walk(dir, walker); err != nil {
		return watcher, fmt.Errorf("error adding paths to watcher: %w", err)
	}
	if err := watcher.Add(c.repoDir); err != nil {
		return watcher, fmt.Errorf("error watching repo dir: %w", err)
	}
	return watcher, err
}

func dedupeWatchEvents(ctx context.Context, watcher *fsnotify.Watcher) chan error {
	ch := make(chan error, 1)
	go func() {
		timer := time.NewTimer(60 * time.Minute)
		stoptimer := func() {
			if !timer.Stop() {
				<-timer.C
			}
		}
		for {
			select {
			case <-watcher.Events:
				stoptimer()
				timer.Reset(throttleDelay)
			case err := <-watcher.Errors:
				ch <- err
			case <-ctx.Done():
				stoptimer()
				return
			case <-timer.C:
				ch <- nil
			}
		}
	}()
	return ch
}
