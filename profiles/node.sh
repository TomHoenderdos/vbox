#!/usr/bin/env bash
# Node.js profile: installs Node via asdf.

profile_ports() {
  echo "3000:3000:Node/Express"
}

profile_provision() {
  local node_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      node_version=$(awk '/^nodejs/ {print $2}' "$tvf")
      [[ -n "$node_version" ]] && break
    fi
  done
  node_version="${node_version:-22.14.0}"

cat <<PROVISION
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add nodejs'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install nodejs ${node_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global nodejs ${node_version}'
PROVISION
}
