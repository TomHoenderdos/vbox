package vagrant

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const pidFile = "rsync-auto.pid"

func pidFilePath(dir string) string {
	return filepath.Join(dir, ".vagrant", pidFile)
}

func readPid(dir string) (int, error) {
	data, err := os.ReadFile(pidFilePath(dir))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func writePid(dir string, pid int) error {
	return os.WriteFile(pidFilePath(dir), []byte(strconv.Itoa(pid)), 0644)
}

func removePidFile(dir string) {
	os.Remove(pidFilePath(dir))
}

// StartRsyncAuto kills any existing rsync-auto process and starts a new one.
func StartRsyncAuto(dir string) error {
	StopRsyncAuto(dir)

	cmd, err := startRsyncAutoProcess(dir)
	if err != nil {
		return fmt.Errorf("start rsync-auto: %w", err)
	}

	pid := cmd.Process.Pid
	if err := writePid(dir, pid); err != nil {
		return fmt.Errorf("write pid file: %w", err)
	}

	fmt.Printf("==> File sync running in background (pid %d)\n", pid)
	return nil
}
