package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "java",
		Description: "Java profile: installs Java via asdf.",
		Ports:       []config.Port{{Guest: 8080, Host: 8080, Label: "Java/Spring"}},
		Excludes:    []string{"target/", "build/", ".gradle/"},
		Provision: func(projectDir string) string {
			javaVersion := versionOr(projectDir, "java", "openjdk-21.0.2")
			return fmt.Sprintf(`
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add java 2>/dev/null; true'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install java %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global java %s'
    grep -qF 'set-java-home' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo ". ~/.asdf/plugins/java/set-java-home.bash" >> ~/.bashrc'
`, javaVersion, javaVersion)
		},
	})
}
