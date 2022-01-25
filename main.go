package main

import (
	"fmt"
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

	// create a new catalog manager which is used to determine versions from git
	// tags, unmarshal resources, and generate the api
	m, err := newCatalogManager(workingDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create catalog manager")
	}

	// process the catalog & all its integrations
	releasePath, err := m.ProcessCatalog()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to process integrations")
	}

	fmt.Printf("::set-output name=release_path::%s\n", releasePath)
}
