package generatecmd

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
	rootConfig rootcmd.Config
	repoDir    string
	tempDir    string
}

func New(rootConfig rootcmd.Config) (*ffcli.Command, error) {
	cfg := Config{
		rootConfig: rootConfig,
	}

	fs := flag.NewFlagSet("catalog-api generate", flag.ExitOnError)

	// register generatecmd flags
	if err := cfg.RegisterFlags(fs); err != nil {
		return nil, fmt.Errorf("error registering flags: %w", err)
	}

	// register global flags
	rootConfig.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "catalog-api generate [flags] <subcommand> [flags]",
		ShortHelp:  "Generate a static Catalog API",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}, nil
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) error {
	defaultRepoDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error retrieving current working directory: %w", err)
	}

	defaultTempDir := os.TempDir()

	fs.StringVar(&c.repoDir, "repo-dir", defaultRepoDir, "path to the catalog repository")
	fs.StringVar(&c.tempDir, "temp-dir", defaultTempDir, "path to a temporary directory for generated files")

	return nil
}

func (c *Config) Exec(context.Context, []string) error {
	// create a temporay directory within c.tempDir to hold the generated api
	// files
	tmpDir, err := os.MkdirTemp(c.tempDir, "")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %w", err)
	}

	// create a staging dir to hold the generated api files used to calculate
	// the checksum of the release
	stagingDir := path.Join(tmpDir, "staging")
	if err := os.Mkdir(stagingDir, 0700); err != nil {
		return fmt.Errorf("error creating staging directory: %w", err)
	}

	// create a release dir to hold the complete set of generated api files
	releaseDir := path.Join(tmpDir, "release")
	if err := os.Mkdir(releaseDir, 0700); err != nil {
		return fmt.Errorf("error creating release directory: %w", err)
	}

	mConfig := catalogmanager.Config{
		RepoDir:    c.repoDir,
		StagingDir: stagingDir,
		ReleaseDir: releaseDir,
	}

	// create a new catalog manager which is used to determine versions from git
	// tags, unmarshal resources, and generate the api
	m, err := catalogmanager.New(mConfig)
	if err != nil {
		return fmt.Errorf("error creating catalog manager: %w", err)
	}

	// process the catalog & all its integrations
	if err := m.ProcessCatalog(); err != nil {
		return fmt.Errorf("error processing catalog: %w", err)
	}

	// print outputs for github actions
	// TODO(jk): enable this with a command flag,
	// e.g. --with-github-action-outputs
	fmt.Printf("::set-output name=release-dir::%s\n", releaseDir)

	return nil
}
