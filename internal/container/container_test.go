package container

import (
	"reflect"
	"testing"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func TestPortArgs(t *testing.T) {
	got := portArgs(&config.Config{
		Ports: []config.Port{
			{Guest: 4000, Host: 4000, Label: "Phoenix"},
			{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
		},
	})
	want := []string{"-p", "4000:4000", "-p", "15432:5432"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("portArgs() = %v, want %v", got, want)
	}
}

func TestEngineBinary(t *testing.T) {
	cases := map[string]string{
		"":          "docker",
		"docker":    "docker",
		"container": "container",
	}
	for engine, want := range cases {
		if got := engineBinary(engine); got != want {
			t.Fatalf("engineBinary(%q) = %q, want %q", engine, got, want)
		}
	}
}
