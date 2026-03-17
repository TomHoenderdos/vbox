package profile

func init() {
	register(&Profile{
		Name:        "base",
		Description: "Always included. Installs asdf, Claude Code, and dev essentials.",
		Provision: func(projectDir string) string {
			return `
    apt-get update
    apt-get install -y wget gnupg2 git curl unzip build-essential python3

    # Install ASDF version manager (skip if already installed)
    if [ ! -d /home/vagrant/.asdf ]; then
      su - vagrant -c 'git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0'
    fi

    # Remove npm-installed claude if present (conflicts with native, blocks auto-update)
    npm uninstall -g @anthropic-ai/claude-code 2>/dev/null || true

    # Install Claude Code (native) — must run as vagrant user for correct ownership
    su - vagrant -c '
      mkdir -p ~/.local/bin
      curl -fsSL https://claude.ai/install.sh | bash
    '

    # Ensure ~/.local is fully owned by vagrant (native install needs write access for auto-updates)
    chown -R vagrant:vagrant /home/vagrant/.local

    # Install GitHub CLI
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list
    apt-get update
    apt-get install -y gh

    # Trust GitHub SSH host key
    su - vagrant -c 'mkdir -p ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null'

    # Configure bashrc (idempotent)
    BASHRC="/home/vagrant/.bashrc"
    grep -qF '.asdf/asdf.sh' "$BASHRC" || su - vagrant -c 'echo ". $HOME/.asdf/asdf.sh" >> ~/.bashrc'
    grep -qF '.asdf/completions' "$BASHRC" || su - vagrant -c 'echo ". $HOME/.asdf/completions/asdf.bash" >> ~/.bashrc'
    grep -qF '.local/bin' "$BASHRC" || echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$BASHRC"
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
settings.pop('enabledPlugins', None)
with open(path, 'w') as f:
    json.dump(settings, f, indent=2)
"
    chown -R vagrant:vagrant /home/vagrant/.claude
`
		},
	})
}
