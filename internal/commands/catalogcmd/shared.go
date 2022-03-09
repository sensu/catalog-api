package catalogcmd

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/catalogloader"
	"github.com/sensu/catalog-api/internal/catalogmanager"
)

func (c *Config) newCatalogManagerFromRepo(ctx context.Context) (catalogmanager.CatalogManager, string, error) {
	var cm catalogmanager.CatalogManager

	repo, err := git.PlainOpen(c.repoDir)
	if err != nil {
		return cm, "", err
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
	cm, outdir, err := c.newCatalogManagerFromRepo(ctx)
	if err != nil {
		return outdir, err
	}

	// process the catalog & all its integrations
	err = cm.ProcessCatalog()
	return outdir, err
}

func (c *Config) watchRepo(ctx context.Context) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, fmt.Errorf("error configuring watcher: %w", err)
	}
	if err := watcher.Add(c.repoDir); err != nil {
		return watcher, fmt.Errorf("error watching repo dir: %w", err)
	}
	log.Info().Msg("watching repo for changes")
	return watcher, err
}
