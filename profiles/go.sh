#!/usr/bin/env bash
# Go profile: installs Go via asdf.

profile_ports() {
  echo "8080:8080:Go HTTP"
}

profile_provision() {
  local go_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      go_version=$(awk '/^golang/ {print $2}' "$tvf")
      [[ -n "$go_version" ]] && break
    fi
  done
  go_version="${go_version:-1.24.1}"

cat <<PROVISION
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add golang'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install golang ${go_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global golang ${go_version}'
PROVISION
}
