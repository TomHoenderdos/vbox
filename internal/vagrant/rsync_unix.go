//go:build !windows

package vagrant

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func startRsyncAutoProcess(dir string) (*exec.Cmd, error) {
	cmd := exec.Command("vagrant", "rsync-auto")
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// StopRsyncAuto stops any running rsync-auto process for the project.
func StopRsyncAuto(dir string) error {
	pid, err := readPid(dir)
	if err == nil {
		_ = syscall.Kill(-pid, syscall.SIGTERM)
		removePidFile(dir)
	}

	// Fallback: pgrep for strays
	out, err := exec.Command("pgrep", "-f", "rsync-auto.*"+dir).Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if p, err := strconv.Atoi(line); err == nil {
				_ = syscall.Kill(p, syscall.SIGTERM)
			}
		}
	}
	return nil
}
