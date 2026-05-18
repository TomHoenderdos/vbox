package docker

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
