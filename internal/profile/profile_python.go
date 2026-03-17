package profile

import (
	"fmt"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func init() {
	register(&Profile{
		Name:        "python",
		Description: "Python profile: installs Python via asdf.",
		Ports:       []config.Port{{Guest: 8000, Host: 8000, Label: "Python/Django"}},
		Excludes:    []string{"__pycache__/", ".venv/", "venv/"},
		Provision: func(projectDir string) string {
			pythonVersion := versionOr(projectDir, "python", "3.13.2")
			return fmt.Sprintf(`
    # Python build dependencies
    apt-get install -y libssl-dev zlib1g-dev libbz2-dev libreadline-dev \
      libsqlite3-dev libffi-dev liblzma-dev

    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add python'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install python %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global python %s'
    su - vagrant -c 'source ~/.asdf/asdf.sh && pip install --upgrade pip'
`, pythonVersion, pythonVersion)
		},
	})
}
