package catalogcmd

import (
	"context"
	"flag"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sensu/catalog-api/internal/commands/rootcmd"
)

var (
	defaultIntegrationsDirName = "integrations"
	defaultRepoDir             = "."
	defaultTempDir             = os.TempDir()
	defaultSnapshot            = false
	defaultWatchMode           = false
	defaultApiURL              = "http://localhost:8080"
)

type Config struct {
	rootConfig          rootcmd.Config
	repoDir             string
	tempDir             string
	integrationsDirName string
	snapshot            bool
	watch               bool
	port                int
	apiURL              string
}

func New(rootConfig rootcmd.Config) *ffcli.Command {
	cfg := Config{
		rootConfig: rootConfig,
	}

	fs := flag.NewFlagSet("catalog-api catalog", flag.ExitOnError)
	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "catalog",
		ShortUsage: "catalog-api catalog [flags] <subcommand> [flags]",
		ShortHelp:  "Validate a Catalog and its integrations",
		FlagSet:    fs,
		Exec:       cfg.Exec,
		Subcommands: []*ffcli.Command{
			cfg.GenerateCommand(),
			cfg.ValidateCommand(),
			cfg.ServerCommand(),
			cfg.PreviewCommand(),
		},
	}
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	// register catalog flags
	c.RegisterCatalogFlags(fs)

	// register global flags
	c.rootConfig.RegisterFlags(fs)
}

func (c *Config) RegisterCatalogFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.repoDir, "repo-dir", defaultRepoDir, "path to the catalog repository")
	fs.StringVar(&c.integrationsDirName, "integrations-dir-name", defaultIntegrationsDirName, "path to the directory containing namespaced integrations")
}

func (c *Config) Exec(context.Context, []string) error {
	// The catalog command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}
