package catalogcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
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
	walker := func(path string, de os.DirEntry, err error) error {
		if de.Type().IsDir() {
			return watcher.Add(path)
		}
		return nil
	}
	if err := filepath.WalkDir(dir, walker); err != nil {
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

type Server interface {
	Start(context.Context)
	Stop(context.Context) error
}

type WatchableServer interface {
	Server
	HandleWatchEvent()
}

func (c *Config) startServerWithWatcher(ctx context.Context, symlink string, server WatchableServer) (err error) {
	// prepare
	cleanup, err := c.prepare(ctx, symlink)
	if err != nil {
		return err
	}
	defer func() {
		cleanup() // inside closure to ensure we deref the correct func
	}()

	// start server
	go server.Start(ctx)

	// watch
	if c.watch {
		process := func() error {
			log.Info().Msg("Filesystem change detected")
			cleanup()
			cleanup, err = c.prepare(ctx, symlink)
			if err != nil {
				return err
			}
			server.HandleWatchEvent()
			return nil
		}
		if err = c.watchRepo(ctx, process); err != nil {
			return
		}
	}

	// configure channel for exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	return server.Stop(ctx)
}

func (c *Config) watchRepo(ctx context.Context, process func() error) (err error) {
	// setup file-system notifications
	watcher, err := c.createWatcher(ctx)
	if err != nil {
		return err
	}

	// spin up process
	go func() {
		defer watcher.Close()
		for {
			select {
			case err := <-dedupeWatchEvents(ctx, watcher):
				if err != nil {
					log.Error().Err(err)
				}
				if err := process(); err != nil {
					log.Error().Err(err).Msg("error occurred while processing")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return
}

func (c *Config) prepare(ctx context.Context, symlink string) (cleanup func(), err error) {
	// configure
	cm, err := c.newCatalogManagerFromRepo(ctx)
	if err != nil {
		return cleanup, err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(cm.tmpdir)
		}
	}()

	// generate
	if err := cm.ProcessCatalog(); err != nil {
		return cleanup, err
	}

	// validate
	if err := cm.ValidateCatalog(); err != nil {
		log.Warn().Err(err)
	}

	// symlink
	if err := os.RemoveAll(symlink); err != nil {
		return cleanup, err
	}
	if err := os.Symlink(filepath.Join(cm.tmpdir, "release"), symlink); err != nil {
		return cleanup, err
	}
	log.Info().Str("path", cm.tmpdir).Msg("API generated")
	return func() {
		_ = os.RemoveAll(cm.tmpdir)
		_ = os.RemoveAll(symlink)
	}, err
}
