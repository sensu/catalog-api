package catalogcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/peterbourgon/ff/v3/ffcli"
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
}

func (c *Config) execGenerate(ctx context.Context, _ []string) error {
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
