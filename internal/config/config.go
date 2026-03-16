package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const ConfFile = ".vbox.conf"

type Port struct {
	Guest int
	Host  int
	Label string
}

type Config struct {
	Name     string
	Profiles []string
	Ports    []Port
	Memory   int
	CPUs     int
	AutoSync bool
}

// HomeDir returns ~/.vbox
func HomeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vbox")
}

// ProjectsDir returns ~/Projects
func ProjectsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Projects")
}

// FindProjectRoot walks up from cwd looking for .vbox.conf
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ConfFile)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a vbox project (no %s found)", ConfFile)
		}
		dir = parent
	}
}

// Load parses a .vbox.conf key=value file from the given directory.
func Load(dir string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(dir, ConfFile))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Memory:   2048,
		CPUs:     2,
		AutoSync: true,
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		// Strip quotes from value
		value = strings.Trim(value, "\"")

		switch key {
		case "VBOX_NAME":
			cfg.Name = value
		case "VBOX_PROFILES":
			if value != "" {
				cfg.Profiles = strings.Fields(value)
			}
		case "VBOX_PORTS":
			cfg.Ports = ParsePorts(value)
		case "VBOX_MEMORY":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.Memory = v
			}
		case "VBOX_CPUS":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.CPUs = v
			}
		case "VBOX_AUTO_SYNC":
			cfg.AutoSync = value == "true"
		}
	}

	return cfg, nil
}

// FindAndLoad finds the project root and loads the config.
func FindAndLoad() (string, *Config, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", nil, err
	}
	cfg, err := Load(root)
	if err != nil {
		return "", nil, err
	}
	return root, cfg, nil
}

// Write writes .vbox.conf in bash-sourceable format.
func (c *Config) Write(dir string) error {
	content := fmt.Sprintf(`# vbox project config
VBOX_NAME="%s"
VBOX_PROFILES="%s"
VBOX_PORTS="%s"
VBOX_MEMORY=%d
VBOX_CPUS=%d
VBOX_AUTO_SYNC=%t
`, c.Name, strings.Join(c.Profiles, " "), c.PortsString(), c.Memory, c.CPUs, c.AutoSync)

	return os.WriteFile(filepath.Join(dir, ConfFile), []byte(content), 0644)
}

// PortsString serializes ports to pipe-separated "guest:host:label" format.
func (c *Config) PortsString() string {
	parts := make([]string, len(c.Ports))
	for i, p := range c.Ports {
		parts[i] = fmt.Sprintf("%d:%d:%s", p.Guest, p.Host, p.Label)
	}
	return strings.Join(parts, "|")
}

// ParsePorts parses a pipe-separated port string like "4000:4000:Phoenix|5432:15432:PostgreSQL".
func ParsePorts(s string) []Port {
	if s == "" {
		return nil
	}
	var ports []Port
	for _, entry := range strings.Split(s, "|") {
		parts := strings.SplitN(entry, ":", 3)
		if len(parts) < 3 {
			continue
		}
		guest, err1 := strconv.Atoi(parts[0])
		host, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}
		ports = append(ports, Port{Guest: guest, Host: host, Label: parts[2]})
	}
	return ports
}
