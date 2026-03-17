package profile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
)

// Profile defines a declarative development environment profile.
type Profile struct {
	Name        string
	Description string
	Ports       []config.Port
	Excludes    []string
	NeedsUSB    bool
	// Provision returns the shell script for provisioning.
	// projectDir is the host project directory (for .tool-versions lookup).
	Provision func(projectDir string) string
}

// Info is a summary of a profile for display purposes.
type Info struct {
	Name        string
	Description string
	Ports       []config.Port
	NeedsUSB    bool
}

// registry holds all built-in profiles keyed by name.
var registry = map[string]*Profile{}

func register(p *Profile) {
	registry[p.Name] = p
}

// Get returns a profile by name, or nil if not found.
func Get(name string) *Profile {
	return registry[name]
}

// Exists checks if a profile exists.
func Exists(name string) bool {
	return registry[name] != nil
}

// List returns all available profiles (excluding base).
func List() ([]Info, error) {
	var infos []Info
	for _, p := range registry {
		if p.Name == "base" {
			continue
		}
		infos = append(infos, Info{
			Name:        p.Name,
			Description: p.Description,
			Ports:       p.Ports,
			NeedsUSB:    p.NeedsUSB,
		})
	}
	return infos, nil
}

// GetPorts returns a profile's port definitions.
func GetPorts(name string) ([]config.Port, error) {
	p := Get(name)
	if p == nil {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	return p.Ports, nil
}

// GetUSB returns whether a profile needs USB passthrough.
func GetUSB(name string) (bool, error) {
	p := Get(name)
	if p == nil {
		return false, nil
	}
	return p.NeedsUSB, nil
}

// GetProvision returns a profile's provisioning script.
func GetProvision(name string, projectDir string) (string, error) {
	p := Get(name)
	if p == nil {
		return "", fmt.Errorf("profile %q not found", name)
	}
	return p.Provision(projectDir), nil
}

// CollectPorts gathers ports from multiple profiles, deduplicating by guest port.
func CollectPorts(profiles []string) ([]config.Port, error) {
	seen := map[int]bool{}
	var result []config.Port
	for _, name := range profiles {
		ports, _ := GetPorts(name)
		for _, p := range ports {
			if !seen[p.Guest] {
				seen[p.Guest] = true
				result = append(result, p)
			}
		}
	}
	return result, nil
}

// CollectExcludes gathers rsync excludes from multiple profiles, deduplicating.
func CollectExcludes(profiles []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, name := range profiles {
		p := Get(name)
		if p == nil {
			continue
		}
		for _, e := range p.Excludes {
			if !seen[e] {
				seen[e] = true
				result = append(result, e)
			}
		}
	}
	return result
}

// readToolVersion reads a tool version from .tool-versions files.
// It checks projectDir first, then $HOME.
func readToolVersion(projectDir, tool string) string {
	paths := []string{
		filepath.Join(projectDir, ".tool-versions"),
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".tool-versions"))
	}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 && fields[0] == tool {
				f.Close()
				return fields[1]
			}
		}
		f.Close()
	}
	return ""
}

// versionOr returns the version from .tool-versions, or the fallback default.
func versionOr(projectDir, tool, fallback string) string {
	if v := readToolVersion(projectDir, tool); v != "" {
		return v
	}
	return fallback
}
