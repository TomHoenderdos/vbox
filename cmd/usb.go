package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var usbCmd = &cobra.Command{
	Use:   "usb",
	Short: "USB device management (Parallels)",
}

var usbListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available USB devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available USB devices:")
		fmt.Println()

		c := exec.Command("prlsrvctl", "usb", "list")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			c2 := exec.Command("prlctl", "usb", "list")
			c2.Stdout = os.Stdout
			c2.Stderr = os.Stderr
			if err := c2.Run(); err != nil {
				return fmt.Errorf("could not list USB devices")
			}
		}
		return nil
	},
}

var usbAttachCmd = &cobra.Command{
	Use:   "attach <device-name>",
	Short: "Attach USB device to VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		vmID, err := vagrant.VMID(root)
		if err != nil {
			return err
		}

		device := args[0]
		c := exec.Command("prlctl", "set", vmID, "--device-set", "usb", "--connect", device)
		if err := c.Run(); err != nil {
			c2 := exec.Command("prlctl", "set", vmID, "--device-add", "usb", "--device", device, "--connect")
			if err := c2.Run(); err != nil {
				return fmt.Errorf("could not attach '%s'. Run 'vbox usb list' to see available devices", device)
			}
		}
		fmt.Printf("==> USB device '%s' attached to VM\n", device)
		return nil
	},
}

var usbDetachCmd = &cobra.Command{
	Use:   "detach <device-name>",
	Short: "Detach USB device from VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		vmID, err := vagrant.VMID(root)
		if err != nil {
			return err
		}

		device := args[0]
		c := exec.Command("prlctl", "set", vmID, "--device-set", "usb", "--disconnect", device)
		if err := c.Run(); err != nil {
			return fmt.Errorf("could not detach '%s'", device)
		}
		fmt.Printf("==> USB device '%s' detached from VM\n", device)
		return nil
	},
}

func init() {
	usbCmd.AddCommand(usbListCmd)
	usbCmd.AddCommand(usbAttachCmd)
	usbCmd.AddCommand(usbDetachCmd)
	rootCmd.AddCommand(usbCmd)
}
