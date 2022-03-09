package catalogcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sensu/catalog-api/internal/catalogloader"
)

func (c *Config) ValidateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("catalog-api catalog validate", flag.ExitOnError)

	// register catalog & global flags
	c.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "validate",
		ShortUsage: "catalog-api catalog validate [flags]",
		ShortHelp:  "Validate a catalog directory and its integrations",
		FlagSet:    fs,
		Exec:       c.rootConfig.PreExec(c.execValidate),
	}
}

func (c *Config) execValidate(context.Context, []string) error {
	loader := catalogloader.NewPathLoader(c.repoDir, c.integrationsDirName)

	cm, err := c.newCatalogManager(loader)
	if err != nil {
		return err
	}

	// validate the catalog & all its integrations
	if err := cm.ValidateCatalog(); err != nil {
		return fmt.Errorf("error validating catalog: %w", err)
	}

	return nil
}
