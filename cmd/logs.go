package cmd

import (
	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var logsFollow bool

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show VM system logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		if logsFollow {
			return vagrant.ExecReplace(root, "ssh", "-c", "sudo journalctl -f")
		}

		return vagrant.Run(root, "ssh", "-c", "sudo journalctl --no-pager -n 100")
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	rootCmd.AddCommand(logsCmd)
}
