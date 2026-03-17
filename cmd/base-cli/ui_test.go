package main

import (
	"testing"

	"github.com/fatih/color"
)

func init() {
	// Disable colors in tests for predictable output.
	color.NoColor = true
}

func TestHighlight(t *testing.T) {
	got := highlight("hello")
	if got != "hello" {
		t.Errorf("highlight(%q) = %q, want %q", "hello", got, "hello")
	}
}

func TestFileRef(t *testing.T) {
	got := fileRef("path/to/file")
	if got != "path/to/file" {
		t.Errorf("fileRef(%q) = %q, want %q", "path/to/file", got, "path/to/file")
	}
}
