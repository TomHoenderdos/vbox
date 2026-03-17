package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "php",
		Description: "PHP profile: installs PHP, extensions, and Composer.",
		Ports:       []config.Port{{Guest: 8000, Host: 8000, Label: "PHP"}},
		Excludes:    []string{"vendor/"},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y php php-cli php-fpm php-mysql php-pgsql php-sqlite3 \
      php-curl php-gd php-mbstring php-xml php-zip composer
`
		},
	})
}
