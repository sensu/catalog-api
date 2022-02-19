package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sensu/catalog-api/internal/commands"
)

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}

func main() {
	// setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	ctx := context.Background()
	rootCmd := commands.AddCommands()
	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fatalErr(fmt.Errorf("error parsing command arguments: %s", err))
	}
	if err := rootCmd.Run(ctx); err != nil {
		fatalErr(err)
	}
}
