package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/catalogloader"
)

func (c *Config) GenerateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("catalog-api catalog generate", flag.ExitOnError)

	// register catalog generate flags
	c.RegisterGenerateFlags(fs)

	// register catalog & global flags
	c.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "catalog-api catalog generate [flags]",
		ShortHelp:  "Generate a static catalog API",
		FlagSet:    fs,
		Exec:       c.rootConfig.PreExec(c.execGenerate),
	}
}

func (c *Config) RegisterGenerateFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.tempDir, "temp-dir", defaultTempDir, "path to a temporary directory for generated files")
	fs.BoolVar(&c.snapshot, "snapshot", false, "generate a catalog api for the current catalog branch")
	fs.BoolVar(&c.watch, "watch", false, "enter watch mode, which rebuilds on file change")
}

func (c *Config) execGenerate(ctx context.Context, _ []string) error {
	if c.watch {
		return c.execGenerateWatcher(ctx)
	}
	return c.execRunGenerate(ctx)
}

func (c *Config) execGenerateWatcher(ctx context.Context) error {
	// produce initial build
	if err := c.execRunGenerate(ctx); err != nil {
		return err
	}

	// configure channel for exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	// setup file-system notifications
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error configuring watcher: %w", err)
	}
	defer watcher.Close()
	if err := watcher.Add(c.repoDir); err != nil {
		return fmt.Errorf("error watching repo dir: %w", err)
	}
	log.Info().Msg("watching repo for changes")

	// wait on fs events or exit signal
	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			log.Info().Msg("update detected, rebuilding...")
			err = c.execRunGenerate(ctx)
			if err != nil {
				log.Error().Err(err)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Error().Err(err)
		case <-exit:
			return nil
		}
	}
}

func (c *Config) execRunGenerate(ctx context.Context) error {
	repo, err := git.PlainOpen(c.repoDir)
	if err != nil {
		return err
	}

	var loader catalogloader.Loader
	if c.snapshot {
		loader = catalogloader.NewSnapshotLoader(repo, c.repoDir, c.integrationsDirName)
	} else {
		loader = catalogloader.NewGitLoader(repo, c.integrationsDirName)
	}

	cm, err := c.newCatalogManager(loader)
	if err != nil {
		return err
	}

	// process the catalog & all its integrations
	if err := cm.ProcessCatalog(); err != nil {
		return fmt.Errorf("error processing catalog: %w", err)
	}

	// print outputs for github actions
	// TODO(jk): enable this with a command flag,
	// e.g. --with-github-action-outputs
	fmt.Printf("::set-output name=release-dir::%s\n", cm.GetConfig().ReleaseDir)

	return nil
}
