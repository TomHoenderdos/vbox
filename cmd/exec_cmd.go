package cmd

import (
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec <command>",
	Short: "Run a command in the VM",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		command := strings.Join(args, " ")
		return vagrant.Run(root, "ssh", "-c", command)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
