package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "elixir",
		Description: "Elixir profile: installs Erlang + Elixir via asdf.",
		Ports:       []config.Port{{Guest: 4000, Host: 4000, Label: "Phoenix"}},
		Excludes:    []string{"_build/", "deps/"},
		Provision: func(projectDir string) string {
			erlangVersion := versionOr(projectDir, "erlang", "28.3.3")
			elixirVersion := versionOr(projectDir, "elixir", "1.19.5-otp-28")
			return fmt.Sprintf(`
    # Erlang build dependencies
    apt-get install -y autoconf m4 libncurses5-dev \
      libwxgtk3.2-dev libwxgtk-webview3.2-dev libgl1-mesa-dev \
      libglu1-mesa-dev libpng-dev libssh-dev xsltproc fop libxml2-utils \
      libncurses-dev openjdk-21-jdk

    # Install Erlang %s and Elixir %s
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add erlang'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add elixir'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install erlang %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install elixir %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global erlang %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global elixir %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && mix local.hex --force'
    su - vagrant -c 'source ~/.asdf/asdf.sh && mix local.rebar --force'
`, erlangVersion, elixirVersion, erlangVersion, elixirVersion, erlangVersion, elixirVersion)
		},
	})
}
