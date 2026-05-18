package cmd

import "testing"

func TestParseCodexRuntimeArgs(t *testing.T) {
	originalDocker := codexDocker
	originalImage := codexDockerImage
	defer func() {
		codexDocker = originalDocker
		codexDockerImage = originalImage
	}()

	codexDocker = false
	codexDockerImage = "node:22-bookworm"

	args, dockerMode, image := parseCodexRuntimeArgs([]string{
		"--docker",
		"--docker-image",
		"custom:latest",
		"--model",
		"gpt-5.2",
	})

	if !dockerMode {
		t.Fatal("dockerMode = false, want true")
	}
	if image != "custom:latest" {
		t.Fatalf("image = %q, want custom:latest", image)
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

func TestCodexHelpDoesNotRequireProject(t *testing.T) {
	if err := codexCmd.RunE(codexCmd, []string{"--help"}); err != nil {
		t.Fatalf("codex help returned error: %v", err)
	}
}
