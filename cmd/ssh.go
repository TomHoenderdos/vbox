package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/TomHoenderdos/vbox/internal/config"
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

		vagrantBin, err := exec.LookPath("vagrant")
		if err != nil {
			return err
		}

		if err := os.Chdir(root); err != nil {
			return err
		}

		return syscall.Exec(vagrantBin, []string{"vagrant", "ssh"}, os.Environ())
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
