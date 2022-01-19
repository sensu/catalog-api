package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// find the working directory which should be the catalog repository that
	// this tool should be run within
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to determine the working directory")
	}
	log.Info().Str("path", workingDir).Msg("Using base directory")

	// create a new integration manager which is used to determine versions from
	// git tags, unmarshal resources, and generate the api
	im, err := newIntegrationManager(workingDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create integration manager")
	}

	// load & process all integrations
	if err := im.ProcessIntegrationNamepaces(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to process integrations")
	}
}
