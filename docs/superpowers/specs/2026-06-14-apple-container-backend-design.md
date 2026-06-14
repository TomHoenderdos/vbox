# Apple Container backend + unified `--backend` flag

Date: 2026-06-14

## Goal

Add Apple Container (`container` CLI) as a runtime backend for the `code` (Claude)
and `codex` commands, alongside the existing Vagrant VM and Docker backends.
Replace the current `--docker` boolean with a unified `--backend vm|docker|container`
flag so either agent (Claude or Codex) can run in any backend.

## Background

- Vagrant + Parallels is the primary backend (`vagrant.ExecReplace`).
- Docker was added as an alternative via `--docker` / `--docker-image` on both
  `code` and `codex`, routed through `internal/docker`.
- Apple's `container` CLI is Docker-syntax-compatible for the flags vbox uses:
  `-v src:dst`, `-w`, `-i`, `-t`, `--rm`, and `-p host:container` (publishes to
  localhost; flags must precede the image name — `internal/docker` already does this).
- Caveats (documented, not blocking): macOS-only; the macOS local-network firewall
  can block published ports; a known port-forwarding bug exists on macOS 26.1.

## Decisions

- **Scope:** Add Apple Container as a third backend (keep Vagrant + Docker).
- **Agent UX:** No new agent selector. Both `code` and `codex` gain the new
  backend choice; running "either claude or codex" is served by both commands
  supporting all backends.
- **Flag UX:** Unified `--backend vm|docker|container` enum (Approach C).
- **`--docker`:** Removed cleanly (no deprecated alias).
- **Package:** Rename `internal/docker` → `internal/container` (engine-agnostic).

## Design

### internal/container (renamed from internal/docker)

- `git mv internal/docker internal/container`; package `container`.
- `Options` gains `Engine string` — `"docker"` or `"container"`.
- `Run` selects the binary from `opts.Engine` (default `"docker"` if empty);
  all other logic (`configMounts`, `portArgs`, image default, bootstrap, stdio,
  `-it` detection) is unchanged.
- `DefaultImage = "node:22-bookworm"` stays.
- Update import path `internal/docker` → `internal/container` in `code.go`,
  `codex.go`, and the renamed test file.

### cmd/backend.go (new)

Shared backend handling for both commands:

- `Backend` type with constants `BackendVM`, `BackendDocker`, `BackendContainer`.
- `parseBackend(s string) (Backend, error)` — validates; bad value errors with
  the list of valid options. Empty/unset → `BackendVM`.
- `engineFor(b Backend) string` — maps `BackendDocker`→`"docker"`,
  `BackendContainer`→`"container"` for `container.Options.Engine`.
- `preflightBackend(b Backend) error` — for docker/container, verify the binary
  is on `PATH` (`exec.LookPath`); friendly "X not installed" error otherwise.

### cmd/code.go

- Replace `--docker` (bool) and `--docker-image` with:
  - `--backend` string (default `"vm"`).
  - `--image` string (default `container.DefaultImage`), used by docker/container.
- Parse + preflight backend. Dispatch:
  - `vm` → existing credential-sync + `vagrant.ExecReplace` path (unchanged).
  - `docker` / `container` → `container.Run` with `Engine: engineFor(b)`,
    `Image: image`, existing Claude bootstrap command.
- `--resume` unchanged.

### cmd/codex.go

- Uses `DisableFlagParsing` + manual `parseCodexRuntimeArgs`. Update it to read
  `--backend <v>` / `--backend=<v>` and `--image <v>` / `--image=<v>` instead of
  `--docker` / `--docker-image`. Keep `--resume` / `-r`.
- Same dispatch as `code.go`: vm path unchanged; docker/container via
  `container.Run` with the codex bootstrap command.

### Docs / help

- `cmd/root.go`: broaden Short/Long beyond "Vagrant" to mention backends.
- `README.md`: document `--backend` and `--image`; note Apple Container caveats
  (macOS-only, firewall, 26.1 bug); remove `--docker` references.

## Testing

- Rename `internal/docker/docker_test.go` → `internal/container/...`; update
  package + add a case asserting the binary chosen matches `Engine`.
- Update `cmd/codex_test.go` for the new arg parsing (`--backend`, `--image`,
  invalid backend value error).
- `go build ./...` and `go vet ./...` clean. Smoke: `./vbox code --help`,
  `./vbox codex --help` show the new flags; invalid `--backend xyz` errors.

## Out of scope

- No agent-picker command, no `.vbox.conf` default-agent setting.
- No change to Vagrant/Parallels machinery, profiles, rsync, or USB.
- Apple Container networking quirks are documented, not worked around.
