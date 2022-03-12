package catalogcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog/log"
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
	fs.BoolVar(&c.snapshot, "snapshot", defaultSnapshot, "generate a catalog api for the current catalog branch")
	fs.BoolVar(&c.watch, "watch", defaultWatchMode, "enter watch mode, which rebuilds on file change")
}

func (c *Config) execGenerate(ctx context.Context, _ []string) error {
	if c.watch {
		return c.execGenerateWatcher(ctx)
	}
	return c.execRunGenerate(ctx)
}

func (c *Config) execGenerateWatcher(ctx context.Context) error {
	// produce initial build
	outdir, err := c.generate(ctx)
	if err != nil {
		return err
	}
	log.Info().Msgf("outdir: %s", outdir)

	// configure channel for exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	// setup file-system notifications
	watcher, err := c.createWatcher(ctx)
	if err != nil {
		return err
	}
	defer watcher.Close()

	// wait on fs events or exit signal
	for {
		lastOutdir := outdir
		select {
		case err := <-dedupeWatchEvents(ctx, watcher):
			if err != nil {
				log.Error().Err(err)
				continue
			}

			log.Debug().Msg("update detected, rebuilding...")
			outdir, err = c.generate(ctx)
			if err != nil {
				log.Error().Err(err)
			}

			log.Debug().Msgf("deleting last tmpdir: %s", lastOutdir)
			if err := os.RemoveAll(lastOutdir); err != nil {
				log.Error().Err(err)
			}

			log.Info().Msgf("new outdir: %s", outdir)
		case <-exit:
			return nil
		}
	}
}

func (c *Config) execRunGenerate(ctx context.Context) error {
	outdir, err := c.generate(ctx)
	if err != nil {
		return err
	}

	// print outputs for github actions
	// TODO(jk): enable this with a command flag,
	// e.g. --with-github-action-outputs
	fmt.Printf("::set-output name=release-dir::%s\n", outdir)
	return nil
}
