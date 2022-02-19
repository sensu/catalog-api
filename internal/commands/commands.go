package commands

import (
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sensu/catalog-api/internal/commands/catalogcmd"
	"github.com/sensu/catalog-api/internal/commands/rootcmd"
)

func AddCommands() *ffcli.Command {
	// create root command
	rootCfg := rootcmd.Config{}
	rootCmd := rootcmd.New(&rootCfg)

	// add subcommands to root command
	rootCmd.Subcommands = []*ffcli.Command{
		catalogcmd.New(rootCfg),
	}

	return rootCmd
}
