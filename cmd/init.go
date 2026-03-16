package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/profile"
	"github.com/TomHoenderdos/vbox/internal/tui"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	initProfiles string
	initMemory   int
	initCPUs     int
	initNoSync   bool
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new vbox project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()

		flagsSet := cmd.Flags().Changed("profile") || cmd.Flags().Changed("memory") ||
			cmd.Flags().Changed("cpus") || cmd.Flags().Changed("no-sync") || len(args) > 0
		interactive := !flagsSet && term.IsTerminal(int(os.Stdin.Fd()))

		var cfg *config.Config

		if interactive {
			var err error
			cfg, err = tui.RunInitWizard(cwd)
			if err != nil {
				return err
			}
		} else {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				name = filepath.Base(cwd)
			}

			profiles := []string{"elixir"}
			if initProfiles != "" {
				profiles = strings.Split(strings.ReplaceAll(initProfiles, " ", ""), ",")
			}

			ports, _ := profile.CollectPorts(profiles)

			cfg = &config.Config{
				Name:     name,
				Profiles: profiles,
				Ports:    ports,
				Memory:   initMemory,
				CPUs:     initCPUs,
				AutoSync: !initNoSync,
			}
		}

		projectDir := cwd
		if cfg.Name != filepath.Base(cwd) {
			projectDir = filepath.Join(config.ProjectsDir(), cfg.Name)
		}

		if _, err := os.Stat(filepath.Join(projectDir, config.ConfFile)); err == nil {
			return fmt.Errorf("already a vbox project")
		}

		os.MkdirAll(projectDir, 0755)
		fmt.Printf("==> Initializing vbox in %s...\n", projectDir)

		if err := cfg.Write(projectDir); err != nil {
			return err
		}

		if err := vagrant.GenerateVagrantfile(projectDir, cfg); err != nil {
			return err
		}

		gitDir := filepath.Join(projectDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			runGitInit(projectDir)
		} else {
			ensureGitignore(projectDir)
		}

		fmt.Printf("==> Project '%s' created at %s\n", cfg.Name, projectDir)
		fmt.Printf("==> Profiles: %s\n", strings.Join(cfg.Profiles, " "))
		for _, p := range cfg.Ports {
			fmt.Printf("==> %s: localhost:%d -> guest:%d\n", p.Label, p.Host, p.Guest)
		}
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", projectDir)
		fmt.Println("  vbox up")
		fmt.Println("  vbox ssh")

		return nil
	},
}

func runGitInit(dir string) {
	c := exec.Command("git", "init", "-q")
	c.Dir = dir
	c.Run()

	gitignore := filepath.Join(dir, ".gitignore")
	os.WriteFile(gitignore, []byte(".vagrant/\n_build/\ndeps/\nnode_modules/\n*.beam\n.vbox.conf\n"), 0644)

	exec.Command("git", "-C", dir, "add", "-A").Run()
	exec.Command("git", "-C", dir, "commit", "-q", "-m", "Initial project setup with vbox").Run()
}

func ensureGitignore(dir string) {
	path := filepath.Join(dir, ".gitignore")
	data, _ := os.ReadFile(path)
	content := string(data)

	for _, entry := range []string{".vagrant/", ".vbox.conf"} {
		if !strings.Contains(content, entry) {
			content += entry + "\n"
		}
	}
	os.WriteFile(path, []byte(content), 0644)
}

func init() {
	initCmd.Flags().StringVar(&initProfiles, "profile", "", "Comma-separated profiles (default: elixir)")
	initCmd.Flags().IntVar(&initMemory, "memory", 2048, "VM memory in MB")
	initCmd.Flags().IntVar(&initCPUs, "cpus", 2, "VM CPU count")
	initCmd.Flags().BoolVar(&initNoSync, "no-sync", false, "Disable auto file sync")
	rootCmd.AddCommand(initCmd)
}
