//go:build !windows

package vagrant

import (
	"os"
	"os/exec"
	"syscall"
)

// ExecReplace replaces the current process with vagrant (unix only).
func ExecReplace(dir string, args ...string) error {
	vagrantBin, err := exec.LookPath("vagrant")
	if err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}
	return syscall.Exec(vagrantBin, append([]string{"vagrant"}, args...), os.Environ())
}
