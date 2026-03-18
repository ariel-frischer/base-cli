package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ariel-frischer/base-cli/internal/config"
)

// pipeInput replaces os.Stdin with a pipe fed by input and resets the shared
// stdinReader. Returns a cleanup function to restore the original stdin.
func pipeInput(t *testing.T, input string) func() {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = io.Copy(w, strings.NewReader(input))
		_ = w.Close()
	}()

	oldStdin := os.Stdin
	os.Stdin = r
	stdinReader = nil // reset shared reader so it picks up new stdin
	return func() {
		os.Stdin = oldStdin
		stdinReader = nil
	}
}

func TestSetupRequiresTerminal(t *testing.T) {
	cleanup := pipeInput(t, "")
	defer cleanup()

	rootCmd.SetArgs([]string{"setup"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error when stdin is not a terminal")
	}
}

func TestSetupSavesConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// 6 fields (accept defaults) + confirm.
	cleanup := pipeInput(t, "\n\n\n\n\n\nY\n")
	defer cleanup()

	cfg := &config.Config{}
	if err := runSetupWithPath(path, cfg); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("loading saved config: %v", err)
	}
	if loaded.Host == "" {
		t.Error("expected host to be set")
	}
}

func TestSetupPreservesExistingValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	existing := &config.Config{
		Host:    "gitlab.com",
		GitUser: "testuser",
		Author:  "Test Author",
		License: "apache2",
		CI:      "gitlab",
		Layout:  "cli",
	}
	if err := config.Save(existing, path); err != nil {
		t.Fatal(err)
	}

	// Accept all defaults (6 fields), then confirm.
	cleanup := pipeInput(t, "\n\n\n\n\n\nY\n")
	defer cleanup()

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := runSetupWithPath(path, loaded); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	saved, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct{ got, want string }{
		"host":     {saved.Host, "gitlab.com"},
		"git_user": {saved.GitUser, "testuser"},
		"license":  {saved.License, "apache2"},
		"ci":       {saved.CI, "gitlab"},
		"layout":   {saved.Layout, "cli"},
	}
	for name, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s: got %q, want %q", name, tc.got, tc.want)
		}
	}
}

func TestSetupCancelDoesNotSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// 6 fields (accept defaults) then decline save.
	cleanup := pipeInput(t, "\n\n\n\n\n\nn\n")
	defer cleanup()

	cfg := &config.Config{}
	if err := runSetupWithPath(path, cfg); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("config file should not exist after cancelling setup")
	}
}

func TestSetupOverridesValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// Provide explicit values for all 6 fields + confirm.
	cleanup := pipeInput(t, "gitlab.com\nmyuser\nJane Doe\napache2\ngitlab\nlib\nY\n")
	defer cleanup()

	cfg := &config.Config{}
	if err := runSetupWithPath(path, cfg); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	saved, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct{ got, want string }{
		"host":     {saved.Host, "gitlab.com"},
		"git_user": {saved.GitUser, "myuser"},
		"author":   {saved.Author, "Jane Doe"},
		"license":  {saved.License, "apache2"},
		"ci":       {saved.CI, "gitlab"},
		"layout":   {saved.Layout, "lib"},
	}
	for name, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s: got %q, want %q", name, tc.got, tc.want)
		}
	}
}
