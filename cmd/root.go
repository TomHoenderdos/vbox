package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vbox",
	Short: "Isolated dev environments with Claude Code and Codex CLI",
	Long:  "vbox - Run Claude Code and Codex CLI in isolated backends: a Vagrant VM, Docker, or Apple Container.",
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
