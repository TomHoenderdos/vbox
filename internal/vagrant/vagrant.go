package vagrant

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run executes a vagrant command with inherited stdio.
func Run(dir string, args ...string) error {
	cmd := exec.Command("vagrant", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunSilent executes a vagrant command and captures its output.
func RunSilent(dir string, args ...string) (string, error) {
	cmd := exec.Command("vagrant", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}

// Status returns the VM state (e.g. "running", "poweroff", "not_created").
func Status(dir string) (string, error) {
	out, err := RunSilent(dir, "status", "--machine-readable")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) >= 4 && fields[2] == "state" {
			return fields[3], nil
		}
	}
	return "unknown", nil
}

// VMID extracts the Parallels VM ID from vagrant machine-readable status.
func VMID(dir string) (string, error) {
	out, err := RunSilent(dir, "status", "--machine-readable")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) >= 4 && fields[2] == "id" {
			return fields[3], nil
		}
	}
	return "", fmt.Errorf("could not determine VM ID")
}
