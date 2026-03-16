#!/usr/bin/env bash
# Security profile: network scanners, packet tools, and pen-testing utilities.

profile_ports() { :; }

profile_provision() {
cat <<'PROVISION'
    apt-get install -y nmap tcpdump wireshark-common netcat-openbsd john hashcat hydra
PROVISION
}
