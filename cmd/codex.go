package cmd

import (
	"fmt"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/container"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var codexBackend string
var codexImage string
var codexResume bool

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

		filteredArgs, backend, image, resume, err := parseCodexRuntimeArgs(args)
		if err != nil {
			return err
		}
		if err := preflightBackend(backend); err != nil {
			return err
		}
		filteredArgs = codexResumeArgs(filteredArgs, resume)

		if backend != BackendVM {
			codexArgs := append([]string{"--dangerously-bypass-approvals-and-sandbox"}, filteredArgs...)
			command := dockerToolBootstrap("ripgrep git ca-certificates") +
				" && npm install -g @openai/codex --no-audit >/dev/null && codex " + shellQuoteArgs(codexArgs)
			fmt.Printf("==> Launching Codex CLI in %s\n", engineFor(backend))
			return container.Run(container.Options{
				ProjectRoot: root,
				Config:      cfg,
				Image:       image,
				Command:     command,
				Engine:      engineFor(backend),
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

func dockerToolBootstrap(packages string) string {
	return "export DEBIAN_FRONTEND=noninteractive; apt-get update >/dev/null && apt-get install -y " + packages + " >/dev/null"
}

func codexResumeArgs(args []string, resume bool) []string {
	if !resume {
		return args
	}
	return append([]string{"resume", "--last"}, args...)
}

func parseCodexRuntimeArgs(args []string) ([]string, Backend, string, bool, error) {
	backendStr := codexBackend
	image := codexImage
	resume := codexResume
	filtered := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--backend" && i+1 < len(args):
			backendStr = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--backend="):
			backendStr = strings.TrimPrefix(args[i], "--backend=")
		case args[i] == "--image" && i+1 < len(args):
			image = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--image="):
			image = strings.TrimPrefix(args[i], "--image=")
		case args[i] == "--resume" || args[i] == "-r":
			resume = true
		default:
			filtered = append(filtered, args[i])
		}
	}

	backend, err := parseBackend(backendStr)
	if err != nil {
		return nil, BackendVM, "", false, err
	}
	return filtered, backend, image, resume, nil
}

func init() {
	codexCmd.Flags().StringVar(&codexBackend, "backend", "vm", "Runtime backend: vm, docker, or container")
	codexCmd.Flags().StringVar(&codexImage, "image", container.DefaultImage, "Container image for docker/container backends")
	codexCmd.Flags().BoolVarP(&codexResume, "resume", "r", false, "Resume the most recent conversation")
	rootCmd.AddCommand(codexCmd)
}
