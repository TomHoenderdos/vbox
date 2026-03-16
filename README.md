# vbox

Vagrant-based isolated dev environments with [Claude Code](https://claude.com/claude-code) pre-installed.

Like [ClaudeBox](https://github.com/RchGrav/claudebox), but with real VM isolation instead of Docker containers — no shared kernel, no Docker socket access, no dangerous capabilities.

## Install

```bash
git clone git@github.com:TomHoenderdos/vbox.git
cd vbox
./install.sh
```

Requires [Vagrant](https://www.vagrantup.com/) and a VM provider ([Parallels](https://www.parallels.com/), VirtualBox, etc).

## Quick start

```bash
# Interactive wizard
vbox init

# Or one-liner
vbox init MyApp --profile elixir,postgres

# Start and connect
vbox up
vbox ssh
```

## Commands

| Command | Description |
|---|---|
| `vbox init [name]` | Init vbox in new or current dir |
| `vbox up` | Start the VM |
| `vbox down` | Stop the VM |
| `vbox down -v` | Stop and destroy the VM |
| `vbox ssh` | SSH into the VM |
| `vbox exec <cmd>` | Run a command in the VM |
| `vbox ps` | Show all vbox projects and status |
| `vbox logs [-f]` | Show VM system logs |
| `vbox sync` | Rsync files to VM |
| `vbox profile list` | List available profiles |
| `vbox profile add <name>` | Add a profile to current project |
| `vbox regen` | Regenerate Vagrantfile from config |

## Init options

```
--profile elixir,postgres    Comma-separated profiles (default: elixir)
--memory 2048                VM memory in MB (default: 2048)
--cpus 2                     VM CPU count (default: 2)
--no-sync                    Disable auto file sync on vbox up
```

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
| `dart` | Dart + Flutter web/server (no device/emulator) | :8080 |
| `c` | GCC, Clang, GDB, Valgrind, CMake | - |
| `esp` | ESP-IDF for ESP32/S2/S3/C3/C6 | :3333 |
| `embedded` | ARM toolchain, OpenOCD, PlatformIO | - |
| `postgres` | PostgreSQL server | :15432 |
| `mysql` | MySQL server | :3306 |
| `redis` | Redis server | :6379 |
| `docker` | Docker Engine inside the VM | - |
| `devops` | Kubernetes, Terraform, Ansible, AWS CLI | - |
| `security` | nmap, tcpdump, Wireshark, John, Hydra | - |
| `web` | Nginx, Apache utils, HTTPie | :8080, :8443 |

Language versions are automatically read from your `.tool-versions` file (asdf).

## How it works

- `vbox init` generates a `Vagrantfile` and `.vbox.conf` from your chosen profiles
- Each profile defines its own provisioning script and port forwards
- `~/.claude` is synced to the VM so Claude Code works out of the box
- `vbox up` starts the VM and runs `rsync-auto` in the background for live file sync
- All commands work from any subdirectory of your project

## Adding custom profiles

Create a file in `~/.vbox/profiles/myprofile.sh`:

```bash
#!/usr/bin/env bash
# My custom profile: does something cool.

profile_ports() {
  echo "9000:9000:MyService"
}

profile_provision() {
cat <<'PROVISION'
    apt-get install -y my-package
PROVISION
}
```

Then use it: `vbox init MyApp --profile myprofile`

## Security (vs ClaudeBox)

| | vbox | ClaudeBox |
|---|---|---|
| Isolation | Full VM (separate kernel) | Docker container (shared kernel) |
| Host filesystem | Rsync only (no live mount) | Direct volume mounts |
| SSH keys | Not mounted | Mounted into container |
| Network capabilities | Standard NAT | NET_ADMIN + NET_RAW |
| Docker socket | Not exposed | Exposed to container |
| Root on host | Not required | Required (Docker group) |

## License

MIT
