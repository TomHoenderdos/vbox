# CLAUDE.md

## Project

vbox - Vagrant-based isolated dev environments with Claude Code pre-installed.

Currently a bash script being rewritten to Go. See PLAN.md for the full architecture.

## Build & test

```bash
go build -o vbox .           # build binary
go vet ./...                 # lint
./vbox help                  # smoke test
./vbox profile list          # test profile loading (needs ~/.vbox/profiles/)
```

## Key conventions

- Profiles are bash scripts in ~/.vbox/profiles/*.sh — Go shells out, never parses them
- Config is key=value in .vbox.conf — must stay bash-sourceable
- Use syscall.Exec (not os/exec) for ssh/code/logs -f so terminal works cleanly
- Use SysProcAttr.Setpgid for rsync-auto process group management
- bubbletea TUI only for `init` (wizard) and `ps` (dashboard), all other commands are plain CLI
- Keep it lean — no unnecessary abstractions

## Implementation task

Read PLAN.md and implement all Go source files. Build order:

1. `internal/config/config.go` — types and config I/O
2. `internal/profile/profile.go` — shell out to profile scripts
3. `internal/vagrant/vagrant.go` — vagrant command wrappers
4. `internal/vagrant/vagrantfile.go` — Vagrantfile generation
5. `internal/vagrant/rsync.go` — rsync-auto management
6. `cmd/root.go` — cobra root
7. `cmd/up.go`, `cmd/down.go`, `cmd/ssh.go`, `cmd/code.go`, `cmd/exec_cmd.go`, `cmd/logs.go`, `cmd/sync.go`, `cmd/usb.go`, `cmd/regen.go`, `cmd/profile.go` — all commands
8. `internal/tui/init_wizard.go` — bubbletea init wizard
9. `internal/tui/ps_dashboard.go` — bubbletea ps dashboard
10. `cmd/init.go`, `cmd/ps.go` — commands that use TUI
11. `main.go` — entry point

After each package, run `go build ./...` to verify compilation.

The existing bash `vbox` script is in the repo root for reference. The profiles/ directory contains all .sh profile scripts.
