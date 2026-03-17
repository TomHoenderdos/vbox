package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "node",
		Description: "Node.js profile: installs Node via asdf.",
		Ports:       []config.Port{{Guest: 3000, Host: 3000, Label: "Node/Express"}},
		Excludes:    []string{"node_modules/"},
		Provision: func(projectDir string) string {
			nodeVersion := versionOr(projectDir, "nodejs", "22.14.0")
			return fmt.Sprintf(`
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add nodejs'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install nodejs %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global nodejs %s'
`, nodeVersion, nodeVersion)
		},
	})
}
