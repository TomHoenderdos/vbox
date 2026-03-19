package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var downDestroy bool

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the VM",
	Long:  "Stop the VM. Use -v to destroy it completely.",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		if downDestroy {
			fmt.Printf("This will destroy the VM and all data for '%s'. Are you sure? [y/N] ", cfg.Name)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" {
				fmt.Println("Aborted.")
				return nil
			}

			vf := filepath.Join(root, "Vagrantfile")
			vd := filepath.Join(root, ".vagrant")
			if _, err := os.Stat(vf); err == nil {
				if _, err := os.Stat(vd); err == nil {
					if err := vagrant.Run(root, "destroy", "-f"); err != nil {
						return err
					}
				}
			} else {
				os.RemoveAll(vd)
				fmt.Println("==> No VM found, cleaned up stale state")
			}
			fmt.Printf("==> VM destroyed for %s. Config kept — run 'vbox up' to recreate.\n", cfg.Name)
		} else {
			if _, err := os.Stat(filepath.Join(root, "Vagrantfile")); os.IsNotExist(err) {
				return fmt.Errorf("no Vagrantfile found. Run 'vbox down -v' to clean up, or 'vbox regen' to recreate")
			}
			if err := vagrant.Run(root, "halt"); err != nil {
				return err
			}
			fmt.Printf("==> VM stopped for %s\n", cfg.Name)
		}
		return nil
	},
}

func init() {
	downCmd.Flags().BoolVarP(&downDestroy, "volumes", "v", false, "Destroy the VM completely")
	rootCmd.AddCommand(downCmd)
}
