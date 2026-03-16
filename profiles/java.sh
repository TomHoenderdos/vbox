#!/usr/bin/env bash
# Java profile: installs Java via asdf.

profile_ports() {
  echo "8080:8080:Java/Spring"
}

profile_provision() {
  local java_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      java_version=$(awk '/^java/ {print $2}' "$tvf")
      [[ -n "$java_version" ]] && break
    fi
  done
  java_version="${java_version:-openjdk-21.0.2}"

cat <<PROVISION
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add java 2>/dev/null; true'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install java ${java_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global java ${java_version}'
    grep -qF 'set-java-home' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo ". ~/.asdf/plugins/java/set-java-home.bash" >> ~/.bashrc'
PROVISION
}
