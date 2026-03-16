//go:build windows

package vagrant

import (
	"os"
	"os/exec"
)

func startRsyncAutoProcess(dir string) (*exec.Cmd, error) {
	cmd := exec.Command("vagrant", "rsync-auto")
	cmd.Dir = dir
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
		if proc, err := os.FindProcess(pid); err == nil {
			_ = proc.Kill()
		}
		removePidFile(dir)
	}
	return nil
}
