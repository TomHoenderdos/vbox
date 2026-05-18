package cmd

import (
	"fmt"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/docker"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var codexDocker bool
var codexDockerImage string

var codexCmd = &cobra.Command{
	Use:                "codex [args...]",
	Short:              "Launch Codex CLI",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
			return cmd.Help()
		}

		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		filteredArgs, dockerMode, dockerImage := parseCodexRuntimeArgs(args)
		if dockerMode {
			command := "npm install -g @openai/codex --no-audit >/dev/null && codex"
			if len(filteredArgs) > 0 {
				command += " " + shellQuoteArgs(filteredArgs)
			}
			fmt.Println("==> Launching Codex CLI in Docker")
			return docker.Run(docker.Options{
				ProjectRoot: root,
				Config:      cfg,
				Image:       dockerImage,
				Command:     command,
			})
		}

		script := "mkdir -p ~/.codex ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null; cd /vagrant && codex"
		if len(filteredArgs) > 0 {
			script += " " + shellQuoteArgs(filteredArgs)
		}

		fmt.Println("==> Launching Codex CLI in VM")
		return vagrant.ExecReplace(root, "ssh", "-c", script)
	},
}

func shellQuoteArgs(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
	}
	return strings.Join(quoted, " ")
}

func parseCodexRuntimeArgs(args []string) ([]string, bool, string) {
	dockerMode := codexDocker
	dockerImage := codexDockerImage
	filtered := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--docker":
			dockerMode = true
		case args[i] == "--docker-image" && i+1 < len(args):
			dockerImage = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--docker-image="):
			dockerImage = strings.TrimPrefix(args[i], "--docker-image=")
		default:
			filtered = append(filtered, args[i])
		}
	}

	return filtered, dockerMode, dockerImage
}

func init() {
	codexCmd.Flags().BoolVar(&codexDocker, "docker", false, "Launch in Docker instead of the VM")
	codexCmd.Flags().StringVar(&codexDockerImage, "docker-image", docker.DefaultImage, "Docker image to use with --docker")
	rootCmd.AddCommand(codexCmd)
}
