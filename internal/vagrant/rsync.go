package vagrant

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const pidFile = "rsync-auto.pid"

func pidFilePath(dir string) string {
	return filepath.Join(dir, ".vagrant", pidFile)
}

// StartRsyncAuto kills any existing rsync-auto process and starts a new one
// in its own process group.
func StartRsyncAuto(dir string) error {
	StopRsyncAuto(dir)

	cmd := exec.Command("vagrant", "rsync-auto")
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start rsync-auto: %w", err)
	}

	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFilePath(dir), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("write pid file: %w", err)
	}

	fmt.Printf("==> File sync running in background (pid %d)\n", pid)
	return nil
}

// StopRsyncAuto stops any running rsync-auto process for the project.
func StopRsyncAuto(dir string) error {
	pf := pidFilePath(dir)
	data, err := os.ReadFile(pf)
	if err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil {
			syscall.Kill(-pid, syscall.SIGTERM)
		}
		os.Remove(pf)
	}

	// Fallback: pgrep for strays
	out, err := exec.Command("pgrep", "-f", "rsync-auto.*"+dir).Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if pid, err := strconv.Atoi(line); err == nil {
				syscall.Kill(pid, syscall.SIGTERM)
			}
		}
	}
	return nil
}
