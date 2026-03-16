package cmd

import (
	"os"

	"github.com/TomHoenderdos/vbox/internal/tui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show all vbox projects and status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return tui.RunPsDashboard()
		}
		tui.PrintPlainPS()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
