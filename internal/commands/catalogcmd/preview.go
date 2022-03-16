package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sensu/catalog-api/internal/catalogpreview"
)

const (
	defaultPreviewPort = 3003
)

func (c *Config) PreviewCommand() *ffcli.Command {
	fs := flag.NewFlagSet("catalog-api catalog preview", flag.ExitOnError)

	// register catalog generate flags
	c.RegisterPreviewFlags(fs)

	// register catalog & global flags
	c.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "preview",
		ShortUsage: "catalog-api catalog preview [flags]",
		ShortHelp:  "Serves static catalog API & preview catalog web application for development purposes",
		FlagSet:    fs,
		Exec:       c.rootConfig.PreExec(c.execPreview),
	}
}

func (c *Config) RegisterPreviewFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.port, "port", defaultPreviewPort, "port to use for dev server")
	fs.StringVar(&c.tempDir, "temp-dir", defaultTempDir, "path to a temporary directory for generated files")
	fs.StringVar(&c.apiURL, "api-url", defaultApiURL, "host URL of Sensu installation; optional")
	fs.BoolVar(&c.snapshot, "without-snapshot", defaultSnapshot, "generate a catalog api using tags only")
	fs.BoolVar(&c.watch, "without-watch", defaultWatchMode, "enter watch mode, which rebuilds on file change")
}

func (c *Config) execPreview(ctx context.Context, _ []string) error {
	// treat snapshot as if it were without-snapshot
	c.snapshot = !c.snapshot
	c.watch = !c.watch

	return c.startPreviewServer(ctx)
}

func (c *Config) startPreviewServer(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// start server
	symlink := filepath.Join(c.tempDir, "current")
	listenAddr := fmt.Sprintf(":%d", c.port)
	server := catalogpreview.NewPreviewServer(listenAddr, symlink, c.apiURL)
	return c.startServerWithWatcher(ctx, symlink, &server)
}
