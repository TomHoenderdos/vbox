package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "redis",
		Description: "Redis profile: installs Redis server for caching/queues.",
		Ports:       []config.Port{{Guest: 6379, Host: 6379, Label: "Redis"}},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y redis-server
    systemctl enable redis-server
    systemctl start redis-server

    # Allow connections from host
    sed -i 's/^bind .*/bind 0.0.0.0/' /etc/redis/redis.conf
    systemctl restart redis-server
`
		},
	})
}
