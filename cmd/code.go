package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/spf13/cobra"
)

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Launch Claude Code in the VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		vagrantBin, err := exec.LookPath("vagrant")
		if err != nil {
			return err
		}

		if err := os.Chdir(root); err != nil {
			return err
		}

		return syscall.Exec(vagrantBin, []string{"vagrant", "ssh", "-c", "cd /vagrant && claude --dangerously-skip-permissions"}, os.Environ())
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)
}
