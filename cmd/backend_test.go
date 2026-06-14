package cmd

import "testing"

func TestParseBackend(t *testing.T) {
	valid := map[string]Backend{
		"":          BackendVM,
		"vm":        BackendVM,
		"docker":    BackendDocker,
		"container": BackendContainer,
	}
	for in, want := range valid {
		got, err := parseBackend(in)
		if err != nil {
			t.Fatalf("parseBackend(%q) error: %v", in, err)
		}
		if got != want {
			t.Fatalf("parseBackend(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := parseBackend("xyz"); err == nil {
		t.Fatal("parseBackend(\"xyz\") error = nil, want error")
	}
}

func TestEngineFor(t *testing.T) {
	if engineFor(BackendDocker) != "docker" {
		t.Fatalf("engineFor(BackendDocker) = %q, want docker", engineFor(BackendDocker))
	}
	if engineFor(BackendContainer) != "container" {
		t.Fatalf("engineFor(BackendContainer) = %q, want container", engineFor(BackendContainer))
	}
}

func TestPreflightBackendVMAlwaysOK(t *testing.T) {
	if err := preflightBackend(BackendVM); err != nil {
		t.Fatalf("preflightBackend(BackendVM) = %v, want nil", err)
	}
}
