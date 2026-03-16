#!/usr/bin/env bash
# PostgreSQL profile: installs and configures PostgreSQL for development.

profile_ports() {
  echo "5432:15432:PostgreSQL"
}

profile_provision() {
cat <<'PROVISION'
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
PROVISION
}
