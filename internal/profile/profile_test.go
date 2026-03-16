package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	dir := Dir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("profiles directory not installed at", dir)
	}

	infos, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(infos) == 0 {
		t.Fatal("List() returned no profiles")
	}

	for _, info := range infos {
		if info.Name == "base" {
			t.Error("List() should exclude base profile")
		}
	}

	names := map[string]bool{}
	for _, info := range infos {
		names[info.Name] = true
	}

	for _, expected := range []string{"go", "elixir", "node", "python"} {
		if !names[expected] {
			t.Errorf("List() missing expected profile %q", expected)
		}
	}
}

func TestListDescriptions(t *testing.T) {
	dir := Dir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("profiles directory not installed at", dir)
	}

	infos, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	for _, info := range infos {
		if info.Description == "" {
			t.Errorf("profile %q has empty description", info.Name)
		}
	}
}

func TestGetPorts(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "go.sh")); os.IsNotExist(err) {
		t.Skip("go.sh profile not installed")
	}

	ports, err := GetPorts("go")
	if err != nil {
		t.Fatalf("GetPorts(go) error: %v", err)
	}

	if len(ports) == 0 {
		t.Fatal("GetPorts(go) returned no ports")
	}

	if ports[0].Guest != 8080 {
		t.Errorf("go profile guest port = %d, want 8080", ports[0].Guest)
	}
}

func TestGetPortsNoOutput(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "rust.sh")); os.IsNotExist(err) {
		t.Skip("rust.sh profile not installed")
	}

	ports, err := GetPorts("rust")
	if err != nil {
		t.Fatalf("GetPorts(rust) error: %v", err)
	}

	if len(ports) != 0 {
		t.Errorf("GetPorts(rust) = %v, want empty", ports)
	}
}

func TestGetUSB(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "go.sh")); os.IsNotExist(err) {
		t.Skip("go.sh profile not installed")
	}

	usb, err := GetUSB("go")
	if err != nil {
		t.Fatalf("GetUSB(go) error: %v", err)
	}
	if usb {
		t.Error("go profile should not need USB")
	}
}

func TestGetUSBEmbedded(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "embedded.sh")); os.IsNotExist(err) {
		t.Skip("embedded.sh profile not installed")
	}

	usb, err := GetUSB("embedded")
	if err != nil {
		t.Fatalf("GetUSB(embedded) error: %v", err)
	}
	if !usb {
		t.Error("embedded profile should need USB")
	}
}

func TestGetProvision(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "base.sh")); os.IsNotExist(err) {
		t.Skip("base.sh profile not installed")
	}

	prov, err := GetProvision("base", "/tmp/testproject")
	if err != nil {
		t.Fatalf("GetProvision(base) error: %v", err)
	}

	if prov == "" {
		t.Error("GetProvision(base) returned empty provision")
	}

	if !strings.Contains(prov, "npm") {
		t.Error("base provision should reference npm")
	}
}

func TestCollectPorts(t *testing.T) {
	if _, err := os.Stat(filepath.Join(Dir(), "go.sh")); os.IsNotExist(err) {
		t.Skip("profiles not installed")
	}

	ports, err := CollectPorts([]string{"go", "node"})
	if err != nil {
		t.Fatalf("CollectPorts() error: %v", err)
	}

	seen := map[int]int{}
	for _, p := range ports {
		seen[p.Guest]++
		if seen[p.Guest] > 1 {
			t.Errorf("duplicate guest port %d", p.Guest)
		}
	}
}

func TestExists(t *testing.T) {
	if _, err := os.Stat(Dir()); os.IsNotExist(err) {
		t.Skip("profiles directory not installed")
	}

	if !Exists("base") {
		t.Error("Exists(base) should be true")
	}
	if Exists("nonexistent_profile_xyz") {
		t.Error("Exists(nonexistent) should be false")
	}
}
