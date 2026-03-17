package profile

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
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
	ports, err := GetPorts("rust")
	if err != nil {
		t.Fatalf("GetPorts(rust) error: %v", err)
	}

	if len(ports) != 0 {
		t.Errorf("GetPorts(rust) = %v, want empty", ports)
	}
}

func TestGetUSB(t *testing.T) {
	usb, err := GetUSB("go")
	if err != nil {
		t.Fatalf("GetUSB(go) error: %v", err)
	}
	if usb {
		t.Error("go profile should not need USB")
	}
}

func TestGetUSBEmbedded(t *testing.T) {
	usb, err := GetUSB("embedded")
	if err != nil {
		t.Fatalf("GetUSB(embedded) error: %v", err)
	}
	if !usb {
		t.Error("embedded profile should need USB")
	}
}

func TestGetProvision(t *testing.T) {
	prov, err := GetProvision("base", "/tmp/testproject")
	if err != nil {
		t.Fatalf("GetProvision(base) error: %v", err)
	}

	if prov == "" {
		t.Error("GetProvision(base) returned empty provision")
	}

	if !strings.Contains(prov, "claude") {
		t.Error("base provision should reference claude")
	}
}

func TestCollectPorts(t *testing.T) {
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
	if !Exists("base") {
		t.Error("Exists(base) should be true")
	}
	if Exists("nonexistent_profile_xyz") {
		t.Error("Exists(nonexistent) should be false")
	}
}
