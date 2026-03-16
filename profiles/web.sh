#!/usr/bin/env bash
# Web profile: Nginx, Apache utils, and HTTP test tools.

profile_ports() {
  echo "80:8080:HTTP"
  echo "443:8443:HTTPS"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y nginx apache2-utils httpie
    systemctl enable nginx
PROVISION
}
