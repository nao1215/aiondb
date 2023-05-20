// Package cmd is a package that contains subcommands for the aion CLI command.
package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aion",
		Short: "aion",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newBugReportCmd())

	return cmd
}

// Execute run leadtime process.
func Execute() int {
	rootCmd := newRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		return 1
	}
	return 0
}
