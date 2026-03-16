# vbox

Isolated dev environments with [Claude Code](https://claude.ai/claude-code) built in.

Spin up a full VM, pick your stack, and start coding with Claude — credentials, plugins, and settings sync automatically from your host.

Like [ClaudeBox](https://github.com/RchGrav/claudebox), but with real VM isolation instead of Docker containers.

## Install

```bash
# macOS Apple Silicon
curl -L https://github.com/TomHoenderdos/vbox/releases/latest/download/vbox-darwin-arm64 -o vbox
chmod +x vbox && mv vbox ~/.local/bin/

# macOS Intel
curl -L https://github.com/TomHoenderdos/vbox/releases/latest/download/vbox-darwin-amd64 -o vbox
chmod +x vbox && mv vbox ~/.local/bin/

# Linux ARM64
curl -L https://github.com/TomHoenderdos/vbox/releases/latest/download/vbox-linux-arm64 -o vbox
chmod +x vbox && sudo mv vbox /usr/local/bin/
```

Requires [Vagrant](https://www.vagrantup.com/) and a VM provider ([Parallels](https://www.parallels.com/), VirtualBox, etc).

## Quick start

```bash
# Interactive wizard — walks you through profiles, ports, resources
vbox init

# Or one-liner
vbox init MyApp --profile elixir,postgres

# Start the VM and open Claude Code
vbox up
vbox code
```

`vbox code` handles everything — syncs credentials from macOS Keychain, configures the VM, and drops you into Claude Code.

## Commands

| Command | Description |
|---|---|
| `vbox init [name]` | Create a new project (interactive wizard or flags) |
| `vbox up` | Start the VM |
| `vbox down` | Stop the VM |
| `vbox down -v` | Stop and destroy the VM |
| `vbox code` | Launch Claude Code in the VM |
| `vbox ssh` | Shell into the VM |
| `vbox exec <cmd>` | Run a command in the VM |
| `vbox ps` | Interactive dashboard — manage all VMs |
| `vbox logs [-f]` | Show VM system logs |
| `vbox sync push` | Sync files host -> VM (with confirmation) |
| `vbox sync pull` | Sync files VM -> host (with confirmation) |
| `vbox usb list` | List available USB devices |
| `vbox usb attach <dev>` | Attach USB device to VM |
| `vbox profile list` | List available profiles |
| `vbox profile add <name>` | Add a profile to current project |
| `vbox regen` | Regenerate Vagrantfile from config |

## Dashboard

`vbox ps` opens an interactive TUI dashboard:

- Arrow keys to navigate
- `u` start, `d` stop, `s` ssh, `c` claude code, `D` destroy
- Live status updates after each action

## Profiles

| Profile | Description | Ports |
|---|---|---|
| `elixir` | Erlang + Elixir via asdf | :4000 |
| `rust` | Rust via rustup | - |
| `python` | Python via asdf | :8000 |
| `go` | Go via asdf | :8080 |
| `node` | Node.js via asdf | :3000 |
| `java` | Java via asdf | :8080 |
| `ruby` | Ruby via asdf | :3000 |
| `php` | PHP + Composer | :8000 |
| `dart` | Dart + Flutter web/server | :8080 |
| `c` | GCC, Clang, GDB, Valgrind, CMake | - |
| `esp` | ESP-IDF for ESP32 (USB passthrough) | :3333 |
| `embedded` | ARM toolchain, OpenOCD, PlatformIO (USB passthrough) | - |
| `postgres` | PostgreSQL server | :15432 |
| `mysql` | MySQL server | :3306 |
| `redis` | Redis server | :6379 |
| `docker` | Docker Engine inside the VM | - |
| `devops` | Kubernetes, Terraform, Ansible, AWS CLI | - |
| `security` | nmap, tcpdump, Wireshark, John, Hydra | - |
| `web` | Nginx, Apache utils, HTTPie | :8080, :8443 |

Language versions are read from `.tool-versions` (asdf). Profiles are composable — use as many as you need.

## Init options

```
--profile elixir,postgres    Comma-separated profiles (default: elixir)
--memory 2048                VM memory in MB (default: 2048)
--cpus 2                     VM CPU count (default: 2)
--no-sync                    Disable auto sync on vbox up
```

## How it works

1. `vbox init` generates a `Vagrantfile` and `.vbox.conf` from your chosen profiles
2. Each profile is a self-contained bash script defining ports and provisioning
3. `vbox up` starts the VM, does an initial file sync, and starts background VM -> host sync
4. `vbox code` syncs Claude credentials from macOS Keychain, patches VM settings, and launches Claude Code — all in one SSH session
5. Files you create in the VM automatically sync back to the host

## Adding custom profiles

Create `~/.vbox/profiles/myprofile.sh`:

```bash
#!/usr/bin/env bash
# My custom profile: one-line description shown in profile list.

profile_ports() {
  echo "9000:9000:MyService"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y my-package
PROVISION
}
```

Then: `vbox init MyApp --profile myprofile`

See [docs/CREATING_PROFILES.md](docs/CREATING_PROFILES.md) for the full guide.

## Security

| | vbox | ClaudeBox |
|---|---|---|
| Isolation | Full VM (separate kernel) | Docker container (shared kernel) |
| Host filesystem | Rsync only (explicit push/pull) | Direct volume mounts |
| SSH keys | Not mounted | Mounted into container |
| Network | Standard NAT | NET_ADMIN + NET_RAW |
| Docker socket | Not exposed | Exposed to container |
| Root on host | Not required | Required (Docker group) |

## License

MIT
