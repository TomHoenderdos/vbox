package cmd

import "testing"

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

func TestCodexResumeArgs(t *testing.T) {
	got := codexResumeArgs([]string{"continue this"}, true)
	want := []string{"resume", "--last", "continue this"}
	if len(got) != len(want) {
		t.Fatalf("codexResumeArgs() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("codexResumeArgs() = %v, want %v", got, want)
		}
	}
}

func TestCodexHelpDoesNotRequireProject(t *testing.T) {
	if err := codexCmd.RunE(codexCmd, []string{"--help"}); err != nil {
		t.Fatalf("codex help returned error: %v", err)
	}
}
