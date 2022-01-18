package main

import (
	"os"
	"path"

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
			Str("service", "test").
			Msgf("Failed to determine the working directory: %s", err)
	}
	log.Info().Str("path", workingDir).Msg("Using base directory")

	// create a new integration manager which is used to determine versions from
	// git tags, unmarshal resources, and generate the api
	im, err := newIntegrationManager(workingDir)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "test").
			Msgf("Failed to create integration manager: %s", err)
	}

	// get a list of namespaces & the integrations that belong to them from the
	// list of git tags
	nsIntegrations, err := im.GetNamespacedIntegrations()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "test").
			Msgf("Failed to get the list of integrations from git tags: %s", err)
	}

	// loop through the list of namespaces & integrations, unmarshal the configs
	// & resource files, and then generate the static api
	for namespace, versionedIntegrations := range nsIntegrations {
		for name, versions := range versionedIntegrations {
			for _, integration := range versions {
				log.Debug().
					Str("namespace", namespace).
					Str("name", name).
					Str("version", integration.BaseVersion()).
					Msg("Processing integration")

				// attempt to get files from git
				integrationPath := path.Join("integrations", namespace, name)
				im.ProcessIntegration(integration, integrationPath)
			}
		}
	}

	noop(nsIntegrations)
}
