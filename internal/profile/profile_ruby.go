package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "ruby",
		Description: "Ruby profile: installs Ruby via asdf.",
		Ports:       []config.Port{{Guest: 3000, Host: 3000, Label: "Rails"}},
		Excludes:    []string{"vendor/bundle/"},
		Provision: func(projectDir string) string {
			rubyVersion := versionOr(projectDir, "ruby", "3.3.6")
			return fmt.Sprintf(`
    # Ruby build dependencies
    apt-get install -y libssl-dev libreadline-dev zlib1g-dev libyaml-dev libffi-dev

    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add ruby'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install ruby %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global ruby %s'
`, rubyVersion, rubyVersion)
		},
	})
}
