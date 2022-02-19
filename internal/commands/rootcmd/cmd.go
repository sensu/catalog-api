package rootcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rs/zerolog"

	cmderrors "github.com/sensu/catalog-api/internal/commands/errors"
)

const DefaultLogLevel = "info"

type ExecFn func(context.Context, []string) error

type Config struct {
	Logger   zerolog.Logger
	LogLevel string
}

func New(cfg *Config) *ffcli.Command {
	fs := flag.NewFlagSet("catalog-api", flag.ExitOnError)
	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "catalog-api",
		ShortUsage: "catalog-api <subcommand> [flags]",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}
}

func usage() string {
	cmd := New(&Config{})

	usageFn := cmd.UsageFunc
	if usageFn == nil {
		usageFn = ffcli.DefaultUsageFunc
	}

	return usageFn(cmd)
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.LogLevel, "log-level", DefaultLogLevel, fmt.Sprintf("log level of this command (%v)", validLogLevels()))
}

func (c *Config) PreExec(fn ExecFn) ExecFn {
	return func(ctx context.Context, args []string) error {
		level, err := zerolog.ParseLevel(c.LogLevel)
		if err != nil {
			return cmderrors.ErrHelpWithMessage{
				Message: "invalid log level",
				ErrHelp: flag.ErrHelp,
			}
		}
		zerolog.SetGlobalLevel(level)

		return fn(ctx, args)
	}
}

func (c *Config) Exec(context.Context, []string) error {
	// The root command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}

func validLogLevels() []string {
	return []string{
		"panic",
		"fatal",
		"error",
		"warn",
		"info",
		"debug",
		"trace",
	}
}
