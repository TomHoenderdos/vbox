#!/usr/bin/env bash
# Ruby profile: installs Ruby via asdf.

profile_ports() {
  echo "3000:3000:Rails"
}

profile_provision() {
  local ruby_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      ruby_version=$(awk '/^ruby/ {print $2}' "$tvf")
      [[ -n "$ruby_version" ]] && break
    fi
  done
  ruby_version="${ruby_version:-3.3.6}"

cat <<PROVISION
    # Ruby build dependencies
    apt-get install -y libssl-dev libreadline-dev zlib1g-dev libyaml-dev libffi-dev

    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add ruby'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install ruby ${ruby_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global ruby ${ruby_version}'
PROVISION
}
