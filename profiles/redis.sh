#!/usr/bin/env bash
# Redis profile: installs Redis server for caching/queues.

profile_ports() {
  echo "6379:6379:Redis"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y redis-server
    systemctl enable redis-server
    systemctl start redis-server

    # Allow connections from host
    sed -i 's/^bind .*/bind 0.0.0.0/' /etc/redis/redis.conf
    systemctl restart redis-server
PROVISION
}
