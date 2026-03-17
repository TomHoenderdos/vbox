package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "postgres",
		Description: "PostgreSQL profile: installs and configures PostgreSQL for development.",
		Ports:       []config.Port{{Guest: 5432, Host: 15432, Label: "PostgreSQL"}},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y postgresql postgresql-contrib
    systemctl enable postgresql
    systemctl start postgresql

    # Create dev superuser (postgres/postgres)
    sudo -u postgres psql -c "CREATE USER postgres WITH PASSWORD 'postgres' SUPERUSER;" 2>/dev/null || \
    sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres' SUPERUSER;"

    # Allow connections from host via forwarded port
    PG_CONF=$(find /etc/postgresql -name postgresql.conf | head -1)
    PG_HBA=$(find /etc/postgresql -name pg_hba.conf | head -1)
    echo "listen_addresses = '*'" >> "$PG_CONF"
    echo "host all all 0.0.0.0/0 md5" >> "$PG_HBA"
    systemctl restart postgresql
`
		},
	})
}
