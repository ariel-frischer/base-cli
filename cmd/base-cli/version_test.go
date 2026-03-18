package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/ariel-frischer/base-cli/internal/version"
)

func TestTruncateCommit(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"long hash":   {"abcdef1234567890", "abcdef12"},
		"exactly 8":   {"abcdef12", "abcdef12"},
		"short hash":  {"abc", "abc"},
		"empty":       {"", ""},
		"dev":         {"dev", "dev"},
		"8+ boundary": {"123456789", "12345678"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := truncateCommit(tt.input)
			if got != tt.want {
				t.Errorf("truncateCommit(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestVersionCommandPlain(t *testing.T) {
	// Set known version values for deterministic output.
	oldVersion := version.Version
	oldCommit := version.Commit
	oldBuild := version.BuildDate
	defer func() {
		version.Version = oldVersion
		version.Commit = oldCommit
		version.BuildDate = oldBuild
	}()
	version.Version = "1.2.3"
	version.Commit = "abc123"
	version.BuildDate = "2026-01-01"

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"version", "--plain"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version --plain failed: %v", err)
	}

	// printPlainVersion writes to stdout via fmt.Printf, not cmd.OutOrStdout(),
	// so we just verify no error. The function itself is simple enough.
}

func TestVersionCommandPretty(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Reset the --plain flag
	versionPlain = false
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version failed: %v", err)
	}
}

func TestHelpCommand(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--help"})

	// Help exits with nil error
	_ = rootCmd.Execute()

	// colorizedHelp writes to stdout via fmt, not through cobra's output.
	// Verifying no panic is the main value here.
}

func TestInitHelpCommand(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init", "--help"})

	_ = rootCmd.Execute()
}

func TestPrintPlainVersionOutput(t *testing.T) {
	oldVersion := version.Version
	oldCommit := version.Commit
	oldBuild := version.BuildDate
	defer func() {
		version.Version = oldVersion
		version.Commit = oldCommit
		version.BuildDate = oldBuild
	}()
	version.Version = "0.1.0"
	version.Commit = "deadbeef"
	version.BuildDate = "2026-03-16"

	// Capture stdout
	old := captureStdout(t, func() {
		printPlainVersion()
	})

	if !strings.Contains(old, "base-cli 0.1.0") {
		t.Errorf("printPlainVersion() missing version, got: %s", old)
	}
	if !strings.Contains(old, "commit: deadbeef") {
		t.Errorf("printPlainVersion() missing commit, got: %s", old)
	}
}

func TestPrintPrettyVersionOutput(t *testing.T) {
	oldVersion := version.Version
	defer func() { version.Version = oldVersion }()
	version.Version = "0.1.0"

	// Just verify no panic
	captureStdout(t, func() {
		printPrettyVersion()
	})
}

// captureStdout captures stdout output from fn and returns it as a string.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStdout := os.Stdout
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
