package cmd

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "One-shot rsync files to VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		if err := vagrant.Run(root, "rsync"); err != nil {
			return err
		}
		fmt.Println("==> Files synced to VM")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
