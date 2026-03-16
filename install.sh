#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="$HOME/.local/bin"
VBOX_DIR="$HOME/.vbox"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Installing vbox..."

# Install main script
mkdir -p "$INSTALL_DIR"
cp "$SCRIPT_DIR/vbox" "$INSTALL_DIR/vbox"
chmod +x "$INSTALL_DIR/vbox"

# Install profiles
mkdir -p "$VBOX_DIR/profiles"
cp "$SCRIPT_DIR/profiles/"*.sh "$VBOX_DIR/profiles/"
chmod +x "$VBOX_DIR/profiles/"*.sh

echo "vbox installed successfully!"
echo ""
echo "Make sure $INSTALL_DIR is in your PATH."
echo "Run 'vbox help' to get started."
