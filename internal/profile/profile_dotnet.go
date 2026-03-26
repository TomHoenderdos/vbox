package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "dotnet",
		Description: ".NET profile: installs .NET SDK via asdf.",
		Ports:       []config.Port{{Guest: 5000, Host: 5000, Label: ".NET HTTP"}, {Guest: 5001, Host: 5001, Label: ".NET HTTPS"}},
		Excludes:    []string{"bin/", "obj/"},
		Provision: func(projectDir string) string {
			dotnetVersion := versionOr(projectDir, "dotnet", "10.0.103")
			return fmt.Sprintf(`
    apt-get install -y libicu-dev libssl-dev

    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add dotnet 2>/dev/null; true'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install dotnet %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global dotnet %s'

    grep -qF 'DOTNET_ROOT' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "export DOTNET_ROOT=\$(dirname \$(asdf which dotnet))" >> ~/.bashrc'
    grep -qF 'DOTNET_CLI_TELEMETRY_OPTOUT' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "export DOTNET_CLI_TELEMETRY_OPTOUT=1" >> ~/.bashrc'
`, dotnetVersion, dotnetVersion)
		},
	})
}
