package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if cfg.Author != "" {
		t.Errorf("expected empty author, got %q", cfg.Author)
	}
}

func TestLoadAndSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	want := &Config{
		Author:       "Alice",
		License:      "apache2",
		CI:           "github",
		Layout:       "cli",
		AgentMD:      "none",
		NoGitInit:    BoolPtr(true),
		NoGoreleaser: BoolPtr(true),
		NoCommunity:  BoolPtr(false),
		NoChangelog:  BoolPtr(true),
	}

	if err := Save(want, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.Author != want.Author {
		t.Errorf("Author: got %q, want %q", got.Author, want.Author)
	}
	if got.License != want.License {
		t.Errorf("License: got %q, want %q", got.License, want.License)
	}
	if got.CI != want.CI {
		t.Errorf("CI: got %q, want %q", got.CI, want.CI)
	}
	if got.Layout != want.Layout {
		t.Errorf("Layout: got %q, want %q", got.Layout, want.Layout)
	}
	if got.AgentMD != want.AgentMD {
		t.Errorf("AgentMD: got %q, want %q", got.AgentMD, want.AgentMD)
	}
	if !BoolVal(got.NoGitInit, false) {
		t.Error("NoGitInit: expected true")
	}
	if !BoolVal(got.NoGoreleaser, false) {
		t.Error("NoGoreleaser: expected true")
	}
	if BoolVal(got.NoCommunity, true) {
		t.Error("NoCommunity: expected false")
	}
	if !BoolVal(got.NoChangelog, false) {
		t.Error("NoChangelog: expected true")
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "config.yaml")

	if err := Save(&Config{Author: "Bob"}, path); err != nil {
		t.Fatalf("Save to nested path: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte(":::invalid"), 0o644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestBoolVal(t *testing.T) {
	tests := map[string]struct {
		ptr      *bool
		fallback bool
		want     bool
	}{
		"nil_true":    {nil, true, true},
		"nil_false":   {nil, false, false},
		"true_false":  {BoolPtr(true), false, true},
		"false_true":  {BoolPtr(false), true, false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := BoolVal(tt.ptr, tt.fallback); got != tt.want {
				t.Errorf("BoolVal = %v, want %v", got, tt.want)
			}
		})
	}
}
