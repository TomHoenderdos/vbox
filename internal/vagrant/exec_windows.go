//go:build windows

package vagrant

import (
	"os"
	"os/exec"
)

// ExecReplace runs vagrant as a child process on Windows (no exec replacement).
func ExecReplace(dir string, args ...string) error {
	cmd := exec.Command("vagrant", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
