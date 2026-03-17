package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "influxdb",
		Description: "InfluxDB profile: installs InfluxDB 2.x time-series database.",
		Ports:       []config.Port{{Guest: 8086, Host: 8086, Label: "InfluxDB"}},
		Provision: func(projectDir string) string {
			return `
    curl -fsSL https://repos.influxdata.com/influxdata-archive_compat.key | gpg --dearmor -o /usr/share/keyrings/influxdb-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/influxdb-archive-keyring.gpg] https://repos.influxdata.com/ubuntu stable main" > /etc/apt/sources.list.d/influxdb.list
    apt-get update
    apt-get install -y influxdb2

    systemctl enable influxdb
    systemctl start influxdb
`
		},
	})
}
