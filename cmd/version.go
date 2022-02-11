package cmd

import (
	"fmt"

	"github.com/emyrk/grow/internal/version"
	"github.com/spf13/cobra"
)

var versioncmd = &cobra.Command{
	Use:   "version",
	Short: "Display app version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", version.Version)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Commit : %s\n", version.CommitSHA1)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Date   : %s\n", version.CompiledDate)

		return nil
	},
}
