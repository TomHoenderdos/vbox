package cmd

import (
	"fmt"
	"os/exec"
)

type Backend int

const (
	BackendVM Backend = iota
	BackendDocker
	BackendContainer
)

// parseBackend maps a flag value to a Backend. Empty defaults to vm.
func parseBackend(s string) (Backend, error) {
	switch s {
	case "", "vm":
		return BackendVM, nil
	case "docker":
		return BackendDocker, nil
	case "container":
		return BackendContainer, nil
	default:
		return BackendVM, fmt.Errorf("invalid --backend %q (valid: vm, docker, container)", s)
	}
}

// engineFor returns the container engine binary name for a backend.
func engineFor(b Backend) string {
	if b == BackendContainer {
		return "container"
	}
	return "docker"
}

// preflightBackend verifies the required binary is installed for container backends.
func preflightBackend(b Backend) error {
	if b == BackendVM {
		return nil
	}
	bin := engineFor(b)
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("%s is not installed or not on PATH", bin)
	}
	return nil
}
