package cmd

import (
	"fmt"
	"os"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		if _, err := os.Stat(root + "/Vagrantfile"); os.IsNotExist(err) {
			fmt.Println("==> Vagrantfile missing, regenerating from config...")
			if err := vagrant.GenerateVagrantfile(root, cfg); err != nil {
				return err
			}
		}

		if err := vagrant.Run(root, append([]string{"up"}, args...)...); err != nil {
			return err
		}

		if cfg.AutoSync {
			if err := vagrant.StartRsyncAuto(root); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not start rsync-auto: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
