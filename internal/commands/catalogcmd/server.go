package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog/log"
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
	fs.BoolVar(&c.snapshot, "snapshot", defaultSnapshot, "generate a catalog api for the current catalog branch")
	fs.BoolVar(&c.watch, "watch", defaultWatchMode, "enter watch mode, which rebuilds on file change")
}

func (c *Config) execServer(ctx context.Context, _ []string) error {
	if c.watch {
		return c.startServerWatcher(ctx)
	}
	return c.startServer(ctx)
}

func (c *Config) startServerWatcher(ctx context.Context) error {
	// ...
	return nil
}

func (c *Config) startServer(ctx context.Context) error {
	// configure
	cm, outdir, err := c.newCatalogManagerFromRepo(ctx)
	if err != nil {
		return err
	}
	defer os.RemoveAll(outdir)

	// generate
	if err := cm.ProcessCatalog(); err != nil {
		return err
	}

	// validate
	if err := cm.ValidateCatalog(); err != nil {
		log.Warn().Err(err)
	}

	// symlink
	symdir := filepath.Join(c.tempDir, "current")
	if err := os.RemoveAll(symdir); err != nil {
		return err
	}
	if err := os.Symlink(filepath.Join(outdir, "release"), symdir); err != nil {
		return err
	}
	defer os.RemoveAll(symdir)

	// start server
	router := http.FileServer(http.Dir(symdir))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.port),
		Handler: router,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// configure channel for exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	return server.Shutdown(ctx)
}
