package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vbox",
	Short: "Vagrant-based isolated dev environments with Claude Code",
	Long:  "vbox - Create and manage Vagrant VMs pre-configured with Claude Code.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
