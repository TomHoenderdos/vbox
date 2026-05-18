package profile

func init() {
	register(&Profile{
		Name:        "base",
		Description: "Always included. Installs asdf, Claude Code, Codex CLI, and dev essentials.",
		Provision: func(projectDir string) string {
			return `
    apt-get update
    apt-get install -y wget gnupg2 git curl unzip build-essential python3 npm

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

    # Install Codex CLI
    npm install -g @openai/codex --no-audit

    # Ensure ~/.local is fully owned by vagrant (native install needs write access for auto-updates)
    chown -R vagrant:vagrant /home/vagrant/.local

    # Install GitHub CLI
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list
    apt-get update
    apt-get install -y gh

    # Trust GitHub SSH host key
    su - vagrant -c 'mkdir -p ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null'

    # Configure environment for all shell types (interactive, non-interactive, login)
    # /etc/profile.d/ is sourced by login shells; we also add to .bashrc for interactive use
    cat > /etc/profile.d/vbox-asdf.sh << 'ASDFPROFILE'
export ASDF_DIR="$HOME/.asdf"
if [ -d "$ASDF_DIR" ]; then
  . "$ASDF_DIR/asdf.sh"
fi
ASDFPROFILE

    cat > /etc/profile.d/vbox-path.sh << 'PATHPROFILE'
export PATH="$HOME/.local/bin:$PATH"
PATHPROFILE

    # Also add to .bashrc for interactive shells (completions, alias, cd)
    BASHRC="/home/vagrant/.bashrc"
    grep -qF '.asdf/asdf.sh' "$BASHRC" || su - vagrant -c 'echo ". $HOME/.asdf/asdf.sh" >> ~/.bashrc'
    grep -qF '.asdf/completions' "$BASHRC" || su - vagrant -c 'echo ". $HOME/.asdf/completions/asdf.bash" >> ~/.bashrc'
    grep -qF '.local/bin' "$BASHRC" || echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$BASHRC"
    grep -qF 'cd /vagrant' "$BASHRC" || echo 'cd /vagrant' >> "$BASHRC"
    grep -qF 'alias claude=' "$BASHRC" || echo 'alias claude="claude --dangerously-skip-permissions"' >> "$BASHRC"

    # Add to .profile so non-interactive SSH commands (used by Claude Code) can find tools
    PROFILE="/home/vagrant/.profile"
    grep -qF '.asdf/asdf.sh' "$PROFILE" || su - vagrant -c 'echo ". \$HOME/.asdf/asdf.sh" >> ~/.profile'
    grep -qF '.local/bin' "$PROFILE" || su - vagrant -c 'echo "export PATH=\$HOME/.local/bin:\$PATH" >> ~/.profile'

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
