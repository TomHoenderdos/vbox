package cmd

import (
	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Shell into the VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		return vagrant.ExecReplace(root, "ssh")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
