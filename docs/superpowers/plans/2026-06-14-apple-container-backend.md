# Apple Container Backend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Apple Container (`container` CLI) as a runtime backend for `code` (Claude) and `codex`, alongside Vagrant VM and Docker, exposed through a unified `--backend vm|docker|container` flag.

**Architecture:** Rename `internal/docker` → engine-agnostic `internal/container` with an `Engine` field selecting the `docker` or `container` binary (Docker-compatible flags). A new shared `cmd/backend.go` parses/validates the `--backend` enum, maps it to an engine, and preflights the binary. `code.go` and `codex.go` dispatch: `vm` → existing Vagrant path; `docker`/`container` → `container.Run`. The old `--docker`/`--docker-image` flags are removed.

**Tech Stack:** Go, cobra, os/exec.

---

## File Structure

- `internal/container/container.go` — renamed from `internal/docker/docker.go`; package `container`; `Options.Engine`; binary selection.
- `internal/container/container_test.go` — renamed from `docker_test.go`; package `container`; adds engine-binary test.
- `cmd/backend.go` — NEW. `Backend` type, `parseBackend`, `engineFor`, `preflightBackend`.
- `cmd/backend_test.go` — NEW. Tests for parse/map/preflight.
- `cmd/code.go` — replace `--docker`/`--docker-image` with `--backend`/`--image`; dispatch.
- `cmd/codex.go` — update manual arg parser to `--backend`/`--image`; dispatch.
- `cmd/codex_test.go` — update for new parser signature.
- `cmd/root.go` — broaden help text.
- `README.md` — document `--backend`/`--image`, Apple Container caveats; drop `--docker`.

---

## Task 1: Rename package and add Engine field

**Files:**
- Rename: `internal/docker/docker.go` → `internal/container/container.go`
- Rename: `internal/docker/docker_test.go` → `internal/container/container_test.go`
- Modify: both files' package + add engine selection

- [ ] **Step 1: git mv the package directory**

```bash
git mv internal/docker internal/container
git mv internal/container/docker.go internal/container/container.go
git mv internal/container/docker_test.go internal/container/container_test.go
```

- [ ] **Step 2: Change package name and test package**

In `internal/container/container.go` change `package docker` → `package container`.
In `internal/container/container_test.go` change `package docker` → `package container`.

- [ ] **Step 3: Write failing test for engine binary selection**

Append to `internal/container/container_test.go`:

```go
func TestEngineBinary(t *testing.T) {
	cases := map[string]string{
		"":          "docker",
		"docker":    "docker",
		"container": "container",
	}
	for engine, want := range cases {
		if got := engineBinary(engine); got != want {
			t.Fatalf("engineBinary(%q) = %q, want %q", engine, got, want)
		}
	}
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `go test ./internal/container/ -run TestEngineBinary -v`
Expected: FAIL — `undefined: engineBinary`

- [ ] **Step 5: Add Engine field and engineBinary, use it in Run**

In `internal/container/container.go`, add the `Engine` field to `Options`:

```go
type Options struct {
	ProjectRoot string
	Config      *config.Config
	Image       string
	Command     string
	Engine      string // "docker" (default) or "container"
}
```

Add the helper:

```go
func engineBinary(engine string) string {
	if engine == "container" {
		return "container"
	}
	return "docker"
}
```

In `Run`, replace `cmd := exec.Command("docker", args...)` with:

```go
	cmd := exec.Command(engineBinary(opts.Engine), args...)
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/container/ -v`
Expected: PASS (`TestPortArgs`, `TestEngineBinary`)

- [ ] **Step 7: Update import paths in callers (compile-only fix)**

In `cmd/code.go` and `cmd/codex.go` change:
`"github.com/TomHoenderdos/vbox/internal/docker"` → `"github.com/TomHoenderdos/vbox/internal/container"`
and every `docker.Run`/`docker.Options`/`docker.DefaultImage` → `container.Run`/`container.Options`/`container.DefaultImage`.

- [ ] **Step 8: Verify build**

Run: `go build ./...`
Expected: builds (flags still old `--docker`; fixed in later tasks). If unused-var errors appear, they are addressed in Tasks 3–4.

- [ ] **Step 9: Commit**

```bash
git add internal/container cmd/code.go cmd/codex.go
git commit -m "Rename internal/docker to internal/container with Engine selection"
```

---

## Task 2: Backend enum helper

**Files:**
- Create: `cmd/backend.go`
- Test: `cmd/backend_test.go`

- [ ] **Step 1: Write the failing test**

Create `cmd/backend_test.go`:

```go
package cmd

import "testing"

func TestParseBackend(t *testing.T) {
	valid := map[string]Backend{
		"":          BackendVM,
		"vm":        BackendVM,
		"docker":    BackendDocker,
		"container": BackendContainer,
	}
	for in, want := range valid {
		got, err := parseBackend(in)
		if err != nil {
			t.Fatalf("parseBackend(%q) error: %v", in, err)
		}
		if got != want {
			t.Fatalf("parseBackend(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := parseBackend("xyz"); err == nil {
		t.Fatal("parseBackend(\"xyz\") error = nil, want error")
	}
}

func TestEngineFor(t *testing.T) {
	if engineFor(BackendDocker) != "docker" {
		t.Fatalf("engineFor(BackendDocker) = %q, want docker", engineFor(BackendDocker))
	}
	if engineFor(BackendContainer) != "container" {
		t.Fatalf("engineFor(BackendContainer) = %q, want container", engineFor(BackendContainer))
	}
}

func TestPreflightBackendVMAlwaysOK(t *testing.T) {
	if err := preflightBackend(BackendVM); err != nil {
		t.Fatalf("preflightBackend(BackendVM) = %v, want nil", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/ -run 'TestParseBackend|TestEngineFor|TestPreflight' -v`
Expected: FAIL — `undefined: Backend` / `parseBackend`

- [ ] **Step 3: Write the implementation**

Create `cmd/backend.go`:

```go
package cmd

import (
	"fmt"
	"os/exec"
)

type Backend int

const (
	BackendVM Backend = iota
	BackendDocker
	BackendContainer
)

// parseBackend maps a flag value to a Backend. Empty defaults to vm.
func parseBackend(s string) (Backend, error) {
	switch s {
	case "", "vm":
		return BackendVM, nil
	case "docker":
		return BackendDocker, nil
	case "container":
		return BackendContainer, nil
	default:
		return BackendVM, fmt.Errorf("invalid --backend %q (valid: vm, docker, container)", s)
	}
}

// engineFor returns the container engine binary name for a backend.
func engineFor(b Backend) string {
	if b == BackendContainer {
		return "container"
	}
	return "docker"
}

// preflightBackend verifies the required binary is installed for container backends.
func preflightBackend(b Backend) error {
	if b == BackendVM {
		return nil
	}
	bin := engineFor(b)
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("%s is not installed or not on PATH", bin)
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/ -run 'TestParseBackend|TestEngineFor|TestPreflight' -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/backend.go cmd/backend_test.go
git commit -m "Add backend enum parsing and preflight helper"
```

---

## Task 3: Wire `code` command to `--backend`/`--image`

**Files:**
- Modify: `cmd/code.go`

- [ ] **Step 1: Replace flag variables**

In `cmd/code.go` replace:

```go
var codeResume bool
var codeDocker bool
var codeDockerImage string
```

with:

```go
var codeResume bool
var codeBackend string
var codeImage string
```

- [ ] **Step 2: Replace the docker dispatch block with backend dispatch**

In `RunE`, after `root, cfg, err := config.FindAndLoad()` (keep that and its error check), replace the `if codeDocker { ... }` block with:

```go
		backend, err := parseBackend(codeBackend)
		if err != nil {
			return err
		}
		if err := preflightBackend(backend); err != nil {
			return err
		}

		if backend != BackendVM {
			claudeArgs := []string{"--dangerously-skip-permissions"}
			if codeResume {
				claudeArgs = append(claudeArgs, "--resume")
			}
			command := dockerToolBootstrap("ripgrep git ca-certificates") +
				" && npm install -g @anthropic-ai/claude-code --no-audit >/dev/null && claude " + shellQuoteArgs(claudeArgs)
			fmt.Printf("==> Launching Claude Code in %s\n", engineFor(backend))
			return container.Run(container.Options{
				ProjectRoot: root,
				Config:      cfg,
				Image:       codeImage,
				Command:     command,
				Engine:      engineFor(backend),
			})
		}
```

Leave the existing VM path (credentials sync + `vagrant.ExecReplace`) unchanged below this block.

- [ ] **Step 3: Replace flag registration in init()**

In `func init()` replace:

```go
	codeCmd.Flags().BoolVar(&codeDocker, "docker", false, "Launch in Docker instead of the VM")
	codeCmd.Flags().StringVar(&codeDockerImage, "docker-image", docker.DefaultImage, "Docker image to use with --docker")
```

with:

```go
	codeCmd.Flags().StringVar(&codeBackend, "backend", "vm", "Runtime backend: vm, docker, or container")
	codeCmd.Flags().StringVar(&codeImage, "image", container.DefaultImage, "Container image for docker/container backends")
```

- [ ] **Step 4: Verify build**

Run: `go build ./...`
Expected: `cmd/code.go` compiles. (`cmd/codex.go` may still fail — fixed in Task 4.)

- [ ] **Step 5: Commit**

```bash
git add cmd/code.go
git commit -m "Wire code command to unified --backend flag"
```

---

## Task 4: Wire `codex` command and update parser

**Files:**
- Modify: `cmd/codex.go`
- Modify: `cmd/codex_test.go`

- [ ] **Step 1: Update the parser test (failing)**

Replace the body of `TestParseCodexRuntimeArgs` in `cmd/codex_test.go` with:

```go
func TestParseCodexRuntimeArgs(t *testing.T) {
	originalBackend := codexBackend
	originalImage := codexImage
	originalResume := codexResume
	defer func() {
		codexBackend = originalBackend
		codexImage = originalImage
		codexResume = originalResume
	}()

	codexBackend = "vm"
	codexImage = "node:22-bookworm"
	codexResume = false

	args, backend, image, resume, err := parseCodexRuntimeArgs([]string{
		"--backend",
		"container",
		"--image",
		"custom:latest",
		"--resume",
		"--model",
		"gpt-5.2",
	})
	if err != nil {
		t.Fatalf("parseCodexRuntimeArgs error: %v", err)
	}
	if backend != BackendContainer {
		t.Fatalf("backend = %v, want BackendContainer", backend)
	}
	if image != "custom:latest" {
		t.Fatalf("image = %q, want custom:latest", image)
	}
	if !resume {
		t.Fatal("resume = false, want true")
	}
	want := []string{"--model", "gpt-5.2"}
	if len(args) != len(want) {
		t.Fatalf("args = %v, want %v", args, want)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("args = %v, want %v", args, want)
		}
	}
}
```

(Leave `TestCodexResumeArgs` and `TestCodexHelpDoesNotRequireProject` unchanged.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/ -run TestParseCodexRuntimeArgs -v`
Expected: FAIL — `undefined: codexBackend` and signature mismatch

- [ ] **Step 3: Update flag vars and parser in codex.go**

Replace:

```go
var codexDocker bool
var codexDockerImage string
var codexResume bool
```

with:

```go
var codexBackend string
var codexImage string
var codexResume bool
```

Replace `parseCodexRuntimeArgs` with:

```go
func parseCodexRuntimeArgs(args []string) ([]string, Backend, string, bool, error) {
	backendStr := codexBackend
	image := codexImage
	resume := codexResume
	filtered := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--backend" && i+1 < len(args):
			backendStr = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--backend="):
			backendStr = strings.TrimPrefix(args[i], "--backend=")
		case args[i] == "--image" && i+1 < len(args):
			image = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--image="):
			image = strings.TrimPrefix(args[i], "--image=")
		case args[i] == "--resume" || args[i] == "-r":
			resume = true
		default:
			filtered = append(filtered, args[i])
		}
	}

	backend, err := parseBackend(backendStr)
	if err != nil {
		return nil, BackendVM, "", false, err
	}
	return filtered, backend, image, resume, nil
}
```

- [ ] **Step 4: Update RunE dispatch in codex.go**

Replace the body after the `--help` check and `config.FindAndLoad()` (keep both) — i.e. the `filteredArgs, dockerMode, dockerImage, resume := ...` line through the end of the docker `if` block — with:

```go
		filteredArgs, backend, image, resume, err := parseCodexRuntimeArgs(args)
		if err != nil {
			return err
		}
		if err := preflightBackend(backend); err != nil {
			return err
		}
		filteredArgs = codexResumeArgs(filteredArgs, resume)

		if backend != BackendVM {
			codexArgs := append([]string{"--dangerously-bypass-approvals-and-sandbox"}, filteredArgs...)
			command := dockerToolBootstrap("ripgrep git ca-certificates") +
				" && npm install -g @openai/codex --no-audit >/dev/null && codex " + shellQuoteArgs(codexArgs)
			fmt.Printf("==> Launching Codex CLI in %s\n", engineFor(backend))
			return container.Run(container.Options{
				ProjectRoot: root,
				Config:      cfg,
				Image:       image,
				Command:     command,
				Engine:      engineFor(backend),
			})
		}
```

Leave the VM script path (`script := "mkdir -p ~/.codex ...`) unchanged below.

- [ ] **Step 5: Update flag registration in codex.go init()**

Replace:

```go
	codexCmd.Flags().BoolVar(&codexDocker, "docker", false, "Launch in Docker instead of the VM")
	codexCmd.Flags().StringVar(&codexDockerImage, "docker-image", docker.DefaultImage, "Docker image to use with --docker")
	codexCmd.Flags().BoolVarP(&codexResume, "resume", "r", false, "Resume the most recent conversation")
```

with:

```go
	codexCmd.Flags().StringVar(&codexBackend, "backend", "vm", "Runtime backend: vm, docker, or container")
	codexCmd.Flags().StringVar(&codexImage, "image", container.DefaultImage, "Container image for docker/container backends")
	codexCmd.Flags().BoolVarP(&codexResume, "resume", "r", false, "Resume the most recent conversation")
```

- [ ] **Step 6: Run tests and build**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all PASS

- [ ] **Step 7: Commit**

```bash
git add cmd/codex.go cmd/codex_test.go
git commit -m "Wire codex command to unified --backend flag"
```

---

## Task 5: Docs and help text

**Files:**
- Modify: `cmd/root.go`
- Modify: `README.md`

- [ ] **Step 1: Broaden root help text**

In `cmd/root.go` replace:

```go
	Short: "Vagrant-based isolated dev environments with Claude Code and Codex CLI",
	Long:  "vbox - Create and manage Vagrant VMs pre-configured with Claude Code and Codex CLI.",
```

with:

```go
	Short: "Isolated dev environments with Claude Code and Codex CLI",
	Long:  "vbox - Run Claude Code and Codex CLI in isolated backends: a Vagrant VM, Docker, or Apple Container.",
```

- [ ] **Step 2: Update README**

In `README.md`, find any `--docker` / `--docker-image` documentation and replace with `--backend vm|docker|container` and `--image`. Add a note:

```markdown
### Backends

`code` and `codex` accept `--backend vm|docker|container` (default `vm`) and
`--image <ref>` (default `node:22-bookworm`, used by docker/container).

- `vm` — Vagrant + Parallels (full VM, supports profiles/ports/USB).
- `docker` — Docker container.
- `container` — Apple Container (macOS only). Note: the macOS local-network
  firewall can block published ports, and macOS 26.1 has a known port-forwarding
  bug (apple/container#919).
```

- [ ] **Step 3: Final verification**

Run: `go build -o vbox . && ./vbox code --help && ./vbox codex --help && ./vbox code --backend xyz 2>&1`
Expected: both helps show `--backend` and `--image`; `--backend xyz` prints `invalid --backend "xyz" (valid: vm, docker, container)`.

- [ ] **Step 4: Commit**

```bash
git add cmd/root.go README.md
git commit -m "Document --backend flag and Apple Container backend"
```

---

## Self-Review Notes

- **Spec coverage:** rename pkg (T1), Engine field (T1), backend.go enum+preflight (T2), code.go wiring (T3), codex.go wiring+parser (T4), drop `--docker` (T3/T4), docs (T5). All spec sections mapped.
- **Type consistency:** `Backend`, `parseBackend`, `engineFor`, `preflightBackend`, `engineBinary`, `Options.Engine` used consistently across tasks. `parseCodexRuntimeArgs` new signature `([]string, Backend, string, bool, error)` matches both its test (T4 S1) and caller (T4 S4).
- **Caveat:** Task 1 Step 8 may surface unused-var errors for `codeDocker`/`codexDocker` until Tasks 3–4 land — noted in-step; full clean build asserted at Task 4 Step 6.
