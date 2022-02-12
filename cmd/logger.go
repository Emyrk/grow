package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func init() {
	RootCmd.PersistentFlags().String("log-level", "debug", "Set the logging level")
}

func MustLogger(cmd *cobra.Command) zerolog.Logger {
	logger, err := getLogger(cmd)
	if err != nil {
		panic(err)
	}
	return logger
}

func getLogger(cmd *cobra.Command) (zerolog.Logger, error) {
	levelStr, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return zerolog.Logger{}, xerrors.Errorf("get log level: %w", err)
	}

	level, err := zerolog.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		return zerolog.Logger{}, xerrors.Errorf("parse log level '%s': %w", levelStr, err)
	}

	cmd.Context()

	logger := log.Level(level).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return logger, nil
}
