#!/usr/bin/env bash
# Python profile: installs Python via asdf.

profile_ports() {
  echo "8000:8000:Python/Django"
}

profile_provision() {
  local python_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      python_version=$(awk '/^python/ {print $2}' "$tvf")
      [[ -n "$python_version" ]] && break
    fi
  done
  python_version="${python_version:-3.13.2}"

cat <<PROVISION
    # Python build dependencies
    apt-get install -y libssl-dev zlib1g-dev libbz2-dev libreadline-dev \\
      libsqlite3-dev libffi-dev liblzma-dev

    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add python'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install python ${python_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global python ${python_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && pip install --upgrade pip'
PROVISION
}
