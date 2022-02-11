package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func init() {
	RootCmd.AddCommand(versioncmd)
}

var RootCmd = &cobra.Command{
	Use:   "grow",
	Short: "The game of growth",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := getLogger(cmd)
		if err != nil {
			return xerrors.Errorf("get logger: %w", err)
		}
		return nil
	},
}
