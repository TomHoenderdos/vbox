package profile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
)

type Info struct {
	Name        string
	Description string
	Ports       []config.Port
	NeedsUSB    bool
}

// Dir returns the profiles directory ~/.vbox/profiles
func Dir() string {
	return filepath.Join(config.HomeDir(), "profiles")
}

// List returns all available profiles (excluding base.sh).
func List() ([]Info, error) {
	matches, err := filepath.Glob(filepath.Join(Dir(), "*.sh"))
	if err != nil {
		return nil, err
	}

	var infos []Info
	for _, path := range matches {
		name := strings.TrimSuffix(filepath.Base(path), ".sh")
		if name == "base" {
			continue
		}

		desc := ""
		data, err := os.ReadFile(path)
		if err == nil {
			lines := strings.SplitN(string(data), "\n", 3)
			if len(lines) >= 2 {
				desc = strings.TrimPrefix(lines[1], "# ")
			}
		}

		infos = append(infos, Info{Name: name, Description: desc})
	}
	return infos, nil
}

// GetPorts shells out to a profile script to get its port definitions.
func GetPorts(name string) ([]config.Port, error) {
	script := filepath.Join(Dir(), name+".sh")
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %q && profile_ports", script))
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	var ports []config.Port
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		var guest, host int
		fmt.Sscanf(parts[0], "%d", &guest)
		fmt.Sscanf(parts[1], "%d", &host)
		if guest > 0 {
			ports = append(ports, config.Port{Guest: guest, Host: host, Label: parts[2]})
		}
	}
	return ports, nil
}

// GetUSB shells out to check if a profile needs USB passthrough.
func GetUSB(name string) (bool, error) {
	script := filepath.Join(Dir(), name+".sh")
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %q && profile_usb", script))
	out, err := cmd.Output()
	if err != nil {
		return false, nil
	}
	return strings.TrimSpace(string(out)) == "true", nil
}

// GetProvision shells out to get a profile's provisioning script.
func GetProvision(name string, projectDir string) (string, error) {
	script := filepath.Join(Dir(), name+".sh")
	cmd := exec.Command("bash", "-c", fmt.Sprintf("PROJECT_DIR=%q source %q && profile_provision", projectDir, script))
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("profile %s provision failed: %w", name, err)
	}
	return string(out), nil
}

// CollectPorts gathers ports from multiple profiles, deduplicating by guest port.
func CollectPorts(profiles []string) ([]config.Port, error) {
	seen := map[int]bool{}
	var result []config.Port

	for _, name := range profiles {
		ports, err := GetPorts(name)
		if err != nil {
			continue
		}
		for _, p := range ports {
			if !seen[p.Guest] {
				seen[p.Guest] = true
				result = append(result, p)
			}
		}
	}
	return result, nil
}

// Exists checks if a profile script exists.
func Exists(name string) bool {
	_, err := os.Stat(filepath.Join(Dir(), name+".sh"))
	return err == nil
}
