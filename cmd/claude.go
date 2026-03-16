package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/vagrant"
)

// getClaudeConfig reads ~/.claude.json from the host.
func getClaudeConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(home + "/.claude.json")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// getClaudeCredentials extracts Claude Code credentials from macOS Keychain.
// Returns empty string on non-macOS or if no credentials found.
func getClaudeCredentials() string {
	if runtime.GOOS != "darwin" {
		return ""
	}

	out, err := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w").Output()
	if err != nil || len(out) == 0 {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// syncClaudeCredentials syncs credentials to VM (used by vbox up, before the ssh session)
func syncClaudeCredentials(root string) {
	creds := getClaudeCredentials()
	if creds == "" {
		return
	}

	err := vagrant.RunSilentInput(root, creds, "ssh", "-c", "mkdir -p ~/.claude && cat > ~/.claude/.credentials.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not sync Claude credentials: %v\n", err)
		return
	}

	fmt.Println("==> Claude Code credentials synced to VM")
}
