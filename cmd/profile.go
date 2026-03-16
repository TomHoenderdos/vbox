package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/profile"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		infos, err := profile.List()
		if err != nil {
			return err
		}

		fmt.Println("Available profiles:")
		for _, info := range infos {
			ports, _ := profile.GetPorts(info.Name)
			if len(ports) > 0 {
				var labels []string
				for _, p := range ports {
					labels = append(labels, fmt.Sprintf(":%d", p.Host))
				}
				fmt.Printf("  %-12s %-45s [%s]\n", info.Name, info.Description, strings.Join(labels, ", "))
			} else {
				fmt.Printf("  %-12s %s\n", info.Name, info.Description)
			}
		}
		return nil
	},
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a profile to current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !profile.Exists(name) {
			return fmt.Errorf("unknown profile: %s", name)
		}

		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		for _, p := range cfg.Profiles {
			if p == name {
				fmt.Printf("==> Profile '%s' already active\n", name)
				return nil
			}
		}

		newPorts, _ := profile.GetPorts(name)
		existingGuests := map[int]bool{}
		for _, p := range cfg.Ports {
			existingGuests[p.Guest] = true
		}

		reader := bufio.NewReader(os.Stdin)
		for _, p := range newPorts {
			if existingGuests[p.Guest] {
				continue
			}
			fmt.Printf("%s port (guest :%d) -> host [:%d]: ", p.Label, p.Guest, p.Host)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			host := p.Host
			if input != "" {
				fmt.Sscanf(input, "%d", &host)
			}
			cfg.Ports = append(cfg.Ports, config.Port{Guest: p.Guest, Host: host, Label: p.Label})
		}

		cfg.Profiles = append(cfg.Profiles, name)

		if err := cfg.Write(root); err != nil {
			return err
		}

		if err := vagrant.GenerateVagrantfile(root, cfg); err != nil {
			return err
		}

		fmt.Printf("==> Profile '%s' added. Run 'vbox up --provision' to apply.\n", name)
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	rootCmd.AddCommand(profileCmd)
}
