#!/usr/bin/env bash
# Elixir profile: installs Erlang + Elixir via asdf.

# guest:default_host:label
profile_ports() {
  echo "4000:4000:Phoenix"
}

profile_provision() {
  local erlang_version elixir_version
  for tvf in "${PROJECT_DIR:-.}/.tool-versions" "$HOME/.tool-versions"; do
    if [[ -f "$tvf" ]]; then
      [[ -z "$erlang_version" ]] && erlang_version=$(awk '/^erlang/ {print $2}' "$tvf")
      [[ -z "$elixir_version" ]] && elixir_version=$(awk '/^elixir/ {print $2}' "$tvf")
    fi
  done
  erlang_version="${erlang_version:-28.3.3}"
  elixir_version="${elixir_version:-1.19.5-otp-28}"

cat <<PROVISION
    # Erlang build dependencies
    apt-get install -y autoconf m4 libncurses5-dev \\
      libwxgtk3.2-dev libwxgtk-webview3.2-dev libgl1-mesa-dev \\
      libglu1-mesa-dev libpng-dev libssh-dev xsltproc fop libxml2-utils \\
      libncurses-dev openjdk-21-jdk

    # Install Erlang ${erlang_version} and Elixir ${elixir_version}
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add erlang'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf plugin add elixir'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install erlang ${erlang_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf install elixir ${elixir_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global erlang ${erlang_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && asdf global elixir ${elixir_version}'
    su - vagrant -c 'source ~/.asdf/asdf.sh && mix local.hex --force'
    su - vagrant -c 'source ~/.asdf/asdf.sh && mix local.rebar --force'
PROVISION
}
