package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/catalogserver"
)

const (
	defaultServerPort = 8083
)

func (c *Config) ServerCommand() *ffcli.Command {
	fs := flag.NewFlagSet("catalog-api catalog server", flag.ExitOnError)

	// register catalog generate flags
	c.RegisterServerFlags(fs)

	// register catalog & global flags
	c.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "server",
		ShortUsage: "catalog-api catalog server [flags]",
		ShortHelp:  "Serves static catalog API for development purposes",
		FlagSet:    fs,
		Exec:       c.rootConfig.PreExec(c.execServer),
	}
}

func (c *Config) RegisterServerFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.port, "port", defaultServerPort, "port to use for dev server")
	fs.StringVar(&c.tempDir, "temp-dir", defaultTempDir, "path to a temporary directory for generated files")
	fs.BoolVar(&c.snapshot, "without-snapshot", defaultSnapshot, "generate a catalog api using tags only")
	fs.BoolVar(&c.watch, "watch", defaultWatchMode, "enter watch mode, which rebuilds on file change")
}

func (c *Config) execServer(ctx context.Context, _ []string) error {
	// treat snapshot as if it were without-snapshot
	c.snapshot = !c.snapshot

	return c.startServer(ctx)
}

func (c *Config) startServer(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prepare
	symlink := filepath.Join(c.tempDir, "current")
	cleanup, err := c.prepare(ctx, symlink)
	if err != nil {
		return err
	}
	defer func() {
		cleanup() // inside closure to ensure we deref the correct func
	}()

	// start server
	listenAddr := fmt.Sprintf(":%d", c.port)
	server := catalogserver.NewCatalogServer(listenAddr, symlink)
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
