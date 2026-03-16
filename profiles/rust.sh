#!/usr/bin/env bash
# Rust profile: installs Rust via rustup.

profile_ports() { :; }

profile_provision() {
cat <<'PROVISION'
    if [ ! -d /home/vagrant/.cargo ]; then
      su - vagrant -c 'curl --proto "=https" --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y'
    fi
    grep -qF '.cargo/env' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "source \$HOME/.cargo/env" >> ~/.bashrc'
PROVISION
}
