package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sensu/catalog-api/internal/catalogmanager"
	"github.com/sensu/catalog-api/internal/commands/rootcmd"
)

type Config struct {
	rootConfig      rootcmd.Config
	repoDir         string
	tempDir         string
	integrationsDir string
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
		},
	}
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) error {
	// register catalog flags
	if err := c.RegisterCatalogFlags(fs); err != nil {
		return fmt.Errorf("error registering catalog flags: %w", err)
	}

	// register global flags
	c.rootConfig.RegisterFlags(fs)

	return nil
}

func (c *Config) RegisterCatalogFlags(fs *flag.FlagSet) error {
	defaultRepoDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error retrieving current working directory: %w", err)
	}

	defaultTempDir := os.TempDir()
	defaultIntegrationsDir := "integrations"

	fs.StringVar(&c.repoDir, "repo-dir", defaultRepoDir, "path to the catalog repository")
	fs.StringVar(&c.tempDir, "temp-dir", defaultTempDir, "path to a temporary directory for generated files")
	fs.StringVar(&c.integrationsDir, "integrations-dir", defaultIntegrationsDir, "path to the directory containing namespaced integrations")

	return nil
}

func (c *Config) newCatalogManager() (catalogmanager.CatalogManager, error) {
	var cm catalogmanager.CatalogManager

	// create a temporay directory within c.tempDir to hold the generated api
	// files
	tmpDir, err := os.MkdirTemp(c.tempDir, "")
	if err != nil {
		return cm, fmt.Errorf("error creating temp directory: %w", err)
	}

	// create a staging dir to hold the generated api files used to calculate
	// the checksum of the release
	stagingDir := path.Join(tmpDir, "staging")
	if err := os.Mkdir(stagingDir, 0700); err != nil {
		return cm, fmt.Errorf("error creating staging directory: %w", err)
	}

	// create a release dir to hold the complete set of generated api files
	releaseDir := path.Join(tmpDir, "release")
	if err := os.Mkdir(releaseDir, 0700); err != nil {
		return cm, fmt.Errorf("error creating release directory: %w", err)
	}

	mCfg := catalogmanager.Config{
		RepoDir:         c.repoDir,
		StagingDir:      stagingDir,
		ReleaseDir:      releaseDir,
		IntegrationsDir: c.integrationsDir,
	}

	// create a new catalog manager which is used to determine versions from git
	// tags, unmarshal resources, and generate the api
	return catalogmanager.New(mCfg)
}

func (c *Config) Exec(context.Context, []string) error {
	// The catalog command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}
