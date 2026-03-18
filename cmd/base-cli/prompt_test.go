package main

import (
	"testing"
)

func TestPromptChoiceRejectsInvalid(t *testing.T) {
	// Feed invalid then valid input.
	cleanup := pipeInput(t, "invalid\nmit\n")
	defer cleanup()

	got := promptChoice("license", []string{"mit", "apache2", "none"}, "mit")
	if got != "mit" {
		t.Errorf("promptChoice: got %q, want mit", got)
	}
}

func TestPromptReturnsDefault(t *testing.T) {
	cleanup := pipeInput(t, "\n")
	defer cleanup()

	got := prompt("test", "default-val")
	if got != "default-val" {
		t.Errorf("prompt: got %q, want default-val", got)
	}
}

func TestPromptReturnsUserInput(t *testing.T) {
	cleanup := pipeInput(t, "custom-val\n")
	defer cleanup()

	got := prompt("test", "default-val")
	if got != "custom-val" {
		t.Errorf("prompt: got %q, want custom-val", got)
	}
}
