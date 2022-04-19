package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/peterbourgon/ff/v3/ffcli"
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
	// configure
	listenAddr := fmt.Sprintf(":%d", c.port)
	symlink := filepath.Join(c.tempDir, "current")
	server := catalogserver.NewCatalogServer(listenAddr, symlink)

	// start server
	return c.startServerWithWatcher(ctx, symlink, &server)
}
