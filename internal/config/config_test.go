package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Port
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single port",
			input: "4000:4000:Phoenix",
			want:  []Port{{Guest: 4000, Host: 4000, Label: "Phoenix"}},
		},
		{
			name:  "multiple ports",
			input: "4000:4000:Phoenix|5432:15432:PostgreSQL",
			want: []Port{
				{Guest: 4000, Host: 4000, Label: "Phoenix"},
				{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
			},
		},
		{
			name:  "label with spaces",
			input: "8080:8080:Go HTTP",
			want:  []Port{{Guest: 8080, Host: 8080, Label: "Go HTTP"}},
		},
		{
			name:  "invalid guest port skipped",
			input: "abc:8080:Bad|4000:4000:Good",
			want:  []Port{{Guest: 4000, Host: 4000, Label: "Good"}},
		},
		{
			name:  "incomplete entry skipped",
			input: "4000:4000",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePorts(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePorts(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPortsString(t *testing.T) {
	cfg := &Config{
		Ports: []Port{
			{Guest: 4000, Host: 4000, Label: "Phoenix"},
			{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
		},
	}
	got := cfg.PortsString()
	want := "4000:4000:Phoenix|5432:15432:PostgreSQL"
	if got != want {
		t.Errorf("PortsString() = %q, want %q", got, want)
	}
}

func TestPortsStringEmpty(t *testing.T) {
	cfg := &Config{}
	got := cfg.PortsString()
	if got != "" {
		t.Errorf("PortsString() = %q, want empty", got)
	}
}

func TestPortsRoundTrip(t *testing.T) {
	original := []Port{
		{Guest: 4000, Host: 4000, Label: "Phoenix"},
		{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
		{Guest: 8080, Host: 8080, Label: "Go HTTP"},
	}
	cfg := &Config{Ports: original}
	serialized := cfg.PortsString()
	parsed := ParsePorts(serialized)

	if !reflect.DeepEqual(parsed, original) {
		t.Errorf("round trip failed: got %v, want %v", parsed, original)
	}
}

func TestLoadAndWrite(t *testing.T) {
	dir := t.TempDir()

	original := &Config{
		Name:     "myproject",
		Profiles: []string{"elixir", "postgres"},
		Ports: []Port{
			{Guest: 4000, Host: 4000, Label: "Phoenix"},
			{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
		},
		Memory:   4096,
		CPUs:     4,
		AutoSync: true,
	}

	if err := original.Write(dir); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ConfFile)); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Name != original.Name {
		t.Errorf("Name = %q, want %q", loaded.Name, original.Name)
	}
	if !reflect.DeepEqual(loaded.Profiles, original.Profiles) {
		t.Errorf("Profiles = %v, want %v", loaded.Profiles, original.Profiles)
	}
	if !reflect.DeepEqual(loaded.Ports, original.Ports) {
		t.Errorf("Ports = %v, want %v", loaded.Ports, original.Ports)
	}
	if loaded.Memory != original.Memory {
		t.Errorf("Memory = %d, want %d", loaded.Memory, original.Memory)
	}
	if loaded.CPUs != original.CPUs {
		t.Errorf("CPUs = %d, want %d", loaded.CPUs, original.CPUs)
	}
	if loaded.AutoSync != original.AutoSync {
		t.Errorf("AutoSync = %v, want %v", loaded.AutoSync, original.AutoSync)
	}
}

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ConfFile), []byte(`VBOX_NAME="test"
VBOX_PROFILES="go"
`), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Memory != 2048 {
		t.Errorf("default Memory = %d, want 2048", cfg.Memory)
	}
	if cfg.CPUs != 2 {
		t.Errorf("default CPUs = %d, want 2", cfg.CPUs)
	}
	if !cfg.AutoSync {
		t.Error("default AutoSync should be true")
	}
}

func TestLoadAutoSyncFalse(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ConfFile), []byte(`VBOX_NAME="test"
VBOX_AUTO_SYNC=false
`), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.AutoSync {
		t.Error("AutoSync should be false")
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Error("Load() should error on missing file")
	}
}

func TestLoadSkipsComments(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ConfFile), []byte(`# vbox project config
VBOX_NAME="test"
# this is a comment
VBOX_PROFILES="go"
`), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Name != "test" {
		t.Errorf("Name = %q, want %q", cfg.Name, "test")
	}
}

func TestWriteBashSourceable(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{
		Name:     "myapp",
		Profiles: []string{"go", "postgres"},
		Ports:    []Port{{Guest: 8080, Host: 8080, Label: "Go HTTP"}},
		Memory:   2048,
		CPUs:     2,
		AutoSync: true,
	}

	if err := cfg.Write(dir); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, ConfFile))
	content := string(data)

	for _, expected := range []string{
		`VBOX_NAME="myapp"`,
		`VBOX_PROFILES="go postgres"`,
		`VBOX_PORTS="8080:8080:Go HTTP"`,
		`VBOX_MEMORY=2048`,
		`VBOX_CPUS=2`,
		`VBOX_AUTO_SYNC=true`,
	} {
		if !strings.Contains(content, expected) {
			t.Errorf("config missing %q in:\n%s", expected, content)
		}
	}
}

func TestFindProjectRoot(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, ConfFile), []byte(`VBOX_NAME="test"`), 0644)

	sub := filepath.Join(root, "src", "deep")
	os.MkdirAll(sub, 0755)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(sub)

	found, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot() error: %v", err)
	}
	if found != root {
		t.Errorf("FindProjectRoot() = %q, want %q", found, root)
	}
}

func TestFindProjectRootNotFound(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	_, err := FindProjectRoot()
	if err == nil {
		t.Error("FindProjectRoot() should error when no .vbox.conf exists")
	}
}
