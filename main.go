package main

import (
	"context"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/commands/generatecmd"
	"github.com/sensu/catalog-api/internal/commands/rootcmd"
)

func main() {
	// setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// create root command
	rootCmd, rootConfig, err := rootcmd.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize root command")
	}

	// create generate command
	generatecmd, err := generatecmd.New(rootConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize generate subcommand")
	}

	// add subcommands to root command
	rootCmd.Subcommands = []*ffcli.Command{
		generatecmd,
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse command arguments")
	}

	if err := rootCmd.Run(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed execution of command")
	}
}
