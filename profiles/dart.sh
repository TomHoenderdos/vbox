#!/usr/bin/env bash
# Dart profile: installs Dart SDK and Flutter for web/server development (no device/emulator support).

profile_ports() {
  echo "8080:8080:Flutter web"
}

profile_provision() {
cat <<'PROVISION'
    # Dart/Flutter dependencies
    apt-get install -y clang cmake ninja-build pkg-config libgtk-3-dev liblzma-dev libstdc++-12-dev

    # Install Flutter (includes Dart SDK)
    if [ ! -d /home/vagrant/.flutter ]; then
      su - vagrant -c 'git clone https://github.com/flutter/flutter.git ~/.flutter -b stable --depth 1'
      su - vagrant -c 'export PATH=$HOME/.flutter/bin:$PATH && flutter precache --web'
    fi
    grep -qF '.flutter/bin' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "export PATH=\$HOME/.flutter/bin:\$HOME/.flutter/bin/cache/dart-sdk/bin:\$PATH" >> ~/.bashrc'
PROVISION
}
