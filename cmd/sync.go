package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync files between host and VM (push/pull)",
	Long: `Sync files between host and VM.

  vbox sync        Push files from host to VM (default)
  vbox sync push   Push files from host to VM
  vbox sync pull   Pull files from VM to host (overwrites local files!)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSyncPush()
	},
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push files from host to VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSyncPush()
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull files from VM to host (overwrites local files!)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSyncPull()
	},
}

func runSyncPush() error {
	root, _, err := config.FindAndLoad()
	if err != nil {
		return err
	}

	fmt.Println("==> Pushing files from host to VM...")
	if err := vagrant.Run(root, "rsync"); err != nil {
		return err
	}
	fmt.Println("==> Files pushed to VM")
	return nil
}

func runSyncPull() error {
	root, _, err := config.FindAndLoad()
	if err != nil {
		return err
	}

	fmt.Println("WARNING: This will overwrite local files with VM contents.")
	fmt.Print("Continue? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("==> Cancelled")
		return nil
	}

	fmt.Println("==> Pulling files from VM to host...")
	if err := vagrant.Run(root, "rsync-back"); err != nil {
		return err
	}
	fmt.Println("==> Files pulled from VM")
	return nil
}

func init() {
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
	rootCmd.AddCommand(syncCmd)
}
