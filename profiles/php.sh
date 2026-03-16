#!/usr/bin/env bash
# PHP profile: installs PHP, extensions, and Composer.

profile_ports() {
  echo "8000:8000:PHP"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y php php-cli php-fpm php-mysql php-pgsql php-sqlite3 \
      php-curl php-gd php-mbstring php-xml php-zip composer
PROVISION
}
