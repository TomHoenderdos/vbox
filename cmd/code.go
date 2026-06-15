package cmd

import (
	"fmt"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/container"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/spf13/cobra"
)

var codeResume bool
var codeBackend string
var codeImage string

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Launch Claude Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, cfg, err := config.FindAndLoad()
		if err != nil {
			return err
		}

		backend, err := parseBackend(codeBackend)
		if err != nil {
			return err
		}
		if err := preflightBackend(backend); err != nil {
			return err
		}

		if backend != BackendVM {
			claudeArgs := []string{"--dangerously-skip-permissions"}
			if codeResume {
				claudeArgs = append(claudeArgs, "--resume")
			}
			// Container runs as root; IS_SANDBOX=1 lets claude accept
			// --dangerously-skip-permissions as root (the container is the sandbox).
			command := dockerToolBootstrap("ripgrep git ca-certificates") +
				" && npm install -g @anthropic-ai/claude-code --no-audit >/dev/null && IS_SANDBOX=1 claude " + shellQuoteArgs(claudeArgs)
			fmt.Printf("==> Launching Claude Code in %s\n", engineFor(backend))
			return container.Run(container.Options{
				ProjectRoot: root,
				Config:      cfg,
				Image:       codeImage,
				Command:     command,
				Engine:      engineFor(backend),
			})
		}

		// Get credentials and config from host
		creds := getClaudeCredentials()
		claudeJson := getClaudeConfig()

		// Build script that writes all config, credentials, and launches claude
		// All in one ssh session so rsync can't overwrite in between
		setupLine := "mkdir -p ~/.local/bin ~/.claude ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null; "
		if creds != "" {
			escaped := strings.ReplaceAll(creds, "'", "'\\''")
			setupLine += fmt.Sprintf("echo '%s' > ~/.claude/.credentials.json && ", escaped)
		}
		if claudeJson != "" {
			escaped := strings.ReplaceAll(claudeJson, "'", "'\\''")
			setupLine += fmt.Sprintf("echo '%s' > ~/.claude.json && ", escaped)
		}

		script := setupLine + `python3 -c "
import json, os
path = os.path.expanduser('~/.claude/settings.json')
d = {}
if os.path.exists(path):
    with open(path) as f: d = json.load(f)
d['skipDangerousModePermissionPrompt'] = True
dirs = d.get('trustedDirectories', [])
if '/vagrant' not in dirs: dirs.append('/vagrant')
d['trustedDirectories'] = dirs
with open(path, 'w') as f: json.dump(d, f, indent=2)
cj = os.path.expanduser('~/.claude.json')
if os.path.exists(cj):
    with open(cj) as f: d2 = json.load(f)
    if d2.get('installMethod') == 'npm-global':
        d2['installMethod'] = 'native'
        with open(cj, 'w') as f: json.dump(d2, f, indent=2)
" && cd /vagrant && claude --dangerously-skip-permissions`
		if codeResume {
			script += " --resume"
		}

		fmt.Println("==> Claude Code credentials synced to VM")
		return vagrant.ExecReplace(root, "ssh", "-c", script)
	},
}

func init() {
	codeCmd.Flags().BoolVarP(&codeResume, "resume", "r", false, "Resume the most recent conversation")
	codeCmd.Flags().StringVar(&codeBackend, "backend", "vm", "Runtime backend: vm, docker, or container")
	codeCmd.Flags().StringVar(&codeImage, "image", container.DefaultImage, "Container image for docker/container backends")
	rootCmd.AddCommand(codeCmd)
}
