package cmd

import (
	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var regenCmd = &cobra.Command{
	Use:   "regen",
	Short: "Regenerate Vagrantfile from config",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}
		return vagrant.GenerateVagrantfile(root, cfg)
	},
}

func init() {
	rootCmd.AddCommand(regenCmd)
}
