#!/usr/bin/env bash
# Base profile: always included. Installs asdf, Claude Code, and dev essentials.

profile_ports() { :; }

profile_provision() {
cat <<'PROVISION'
    apt-get update
    apt-get install -y wget gnupg2 git curl unzip build-essential python3

    # Install ASDF version manager (skip if already installed)
    if [ ! -d /home/vagrant/.asdf ]; then
      su - vagrant -c 'git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0'
    fi

    # Install Node.js and Claude Code
    apt-get install -y nodejs npm
    npm install -g @anthropic-ai/claude-code --no-audit

    # Install GitHub CLI
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list
    apt-get update
    apt-get install -y gh

    # Configure bashrc (idempotent)
    BASHRC="/home/vagrant/.bashrc"
    grep -qF '.asdf/asdf.sh' "$BASHRC" || su - vagrant -c 'echo ". \$HOME/.asdf/asdf.sh" >> ~/.bashrc'
    grep -qF '.asdf/completions' "$BASHRC" || su - vagrant -c 'echo ". \$HOME/.asdf/completions/asdf.bash" >> ~/.bashrc'
    grep -qF 'cd /vagrant' "$BASHRC" || echo 'cd /vagrant' >> "$BASHRC"
    grep -qF 'alias claude=' "$BASHRC" || echo 'alias claude="claude --dangerously-skip-permissions"' >> "$BASHRC"


    # Claude Code: merge vbox defaults into existing settings (synced from host)
    mkdir -p /home/vagrant/.claude
    CLAUDE_SETTINGS="/home/vagrant/.claude/settings.json"
    python3 -c "
import json, os
path = '$CLAUDE_SETTINGS'
settings = {}
if os.path.exists(path):
    with open(path) as f:
        settings = json.load(f)
settings['skipDangerousModePermissionPrompt'] = True
dirs = settings.get('trustedDirectories', [])
if '/vagrant' not in dirs:
    dirs.append('/vagrant')
settings['trustedDirectories'] = dirs
with open(path, 'w') as f:
    json.dump(settings, f, indent=2)
"
    chown -R vagrant:vagrant /home/vagrant/.claude
PROVISION
}
