package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "web",
		Description: "Web profile: Nginx, Apache utils, and HTTP test tools.",
		Ports:       []config.Port{{Guest: 80, Host: 8080, Label: "HTTP"}, {Guest: 443, Host: 8443, Label: "HTTPS"}},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y nginx apache2-utils httpie
    systemctl enable nginx
`
		},
	})
}
