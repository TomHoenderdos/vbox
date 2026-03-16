package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

		fmt.Println("==> Claude Code credentials synced to VM")
		return syscall.Exec(vagrantBin, []string{"vagrant", "ssh", "-c", script}, os.Environ())
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)
}
