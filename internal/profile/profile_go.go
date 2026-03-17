package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "go",
		Description: "Go profile: installs Go via asdf.",
		Ports:       []config.Port{{Guest: 8080, Host: 8080, Label: "Go HTTP"}},
		Excludes:    []string{"vendor/"},
		Provision: func(projectDir string) string {
			goVersion := versionOr(projectDir, "golang", "1.24.1")
			return fmt.Sprintf(`
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add golang'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install golang %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global golang %s'
`, goVersion, goVersion)
		},
	})
}
