package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TomHoenderdos/vbox/internal/config"
	"golang.org/x/term"
)

const DefaultImage = "node:22-bookworm"

type Options struct {
	ProjectRoot string
	Config      *config.Config
	Image       string
	Command     string
	Engine      string // "docker" (default) or "container"
}

func engineBinary(engine string) string {
	if engine == "container" {
		return "container"
	}
	return "docker"
}

func Run(opts Options) error {
	image := opts.Image
	if image == "" {
		image = DefaultImage
	}

	args := []string{"run", "--rm"}
	if stdinIsTerminal() {
		args = append(args, "-it")
	}
	args = append(args, "-w", "/workspace")
	args = append(args, "-v", opts.ProjectRoot+":/workspace")
	args = append(args, configMounts()...)
	args = append(args, portArgs(opts.Config)...)
	args = append(args, image, "bash", "-lc", opts.Command)

	cmd := exec.Command(engineBinary(opts.Engine), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func configMounts() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	mounts := []struct {
		host      string
		container string
	}{
		{filepath.Join(home, ".codex"), "/root/.codex"},
		{filepath.Join(home, ".claude"), "/root/.claude"},
		{filepath.Join(home, ".claude.json"), "/root/.claude.json"},
		{filepath.Join(home, ".config", "gh"), "/root/.config/gh"},
	}

	var args []string
	for _, mount := range mounts {
		if _, err := os.Stat(mount.host); err == nil {
			args = append(args, "-v", mount.host+":"+mount.container)
		}
	}
	return args
}

func portArgs(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}

	var args []string
	for _, port := range cfg.Ports {
		args = append(args, "-p", fmt.Sprintf("%d:%d", port.Host, port.Guest))
	}
	return args
}
