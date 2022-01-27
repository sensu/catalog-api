package rootcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog"
)

type Config struct {
	Logger   zerolog.Logger
	LogLevel string
}

func New() (*ffcli.Command, Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("catalog-api", flag.ExitOnError)
	if err := cfg.RegisterFlags(fs); err != nil {
		return nil, cfg, fmt.Errorf("error registering flags: %w", err)
	}

	return &ffcli.Command{
		Name:       "catalog-api",
		ShortUsage: "catalog-api [flags] <subcommand> [flags] [<arg>...]",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}, cfg, nil
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) error {
	defaultLogLevel := "info"

	fs.StringVar(&c.LogLevel, "log-level", defaultLogLevel, "log level of this command")

	return nil
}

func (c *Config) Exec(context.Context, []string) error {
	// The root command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}

// TODO(jk): add parsing for log level & set it
// level, err := zerolog.ParseLevel(&c.LogLevel)
// zerolog.SetGlobalLevel(l zerolog.Level)
