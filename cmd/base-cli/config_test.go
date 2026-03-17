package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ariel-frischer/base-cli/internal/config"
)

func TestConfigInitCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &config.Config{}
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("config file should exist after save")
	}
}

func TestConfigSetAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &config.Config{
		License:      "apache2",
		CI:           "gitlab",
		NoGoreleaser: config.BoolPtr(true),
	}

	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.License != "apache2" {
		t.Errorf("License: got %q, want apache2", loaded.License)
	}
	if loaded.CI != "gitlab" {
		t.Errorf("CI: got %q, want gitlab", loaded.CI)
	}
	if !config.BoolVal(loaded.NoGoreleaser, false) {
		t.Error("NoGoreleaser: expected true")
	}
}

func TestApplyConfigDefaults(t *testing.T) {
	resetInitFlags()

	userCfg := &config.Config{
		License:      "apache2",
		CI:           "github",
		Layout:       "cli",
		AgentMD:      "none",
		NoGoreleaser: config.BoolPtr(true),
		NoCommunity:  config.BoolPtr(true),
	}

	// Reset cobra's Changed state by re-parsing with no flags
	initCmd.ParseFlags([]string{})

	applyConfigDefaults(initCmd, userCfg)

	if flagLicense != "apache2" {
		t.Errorf("License: got %q, want apache2", flagLicense)
	}
	if flagCI != "github" {
		t.Errorf("CI: got %q, want github", flagCI)
	}
	if flagLayout != "cli" {
		t.Errorf("Layout: got %q, want cli", flagLayout)
	}
	if flagAgentMD != "none" {
		t.Errorf("AgentMD: got %q, want none", flagAgentMD)
	}
	if !flagNoGoreleaser {
		t.Error("NoGoreleaser: expected true")
	}
	if !flagNoCommunity {
		t.Error("NoCommunity: expected true")
	}
}

func TestApplyConfigDefaultsFlagOverride(t *testing.T) {
	resetInitFlags()

	userCfg := &config.Config{
		License: "apache2",
		CI:      "gitlab",
	}

	// Parse with explicit --license flag so cobra marks it as Changed
	initCmd.ParseFlags([]string{"--license", "mit"})
	applyConfigDefaults(initCmd, userCfg)

	// Explicit flag should win
	if flagLicense != "mit" {
		t.Errorf("License: got %q, want mit (explicit flag should override config)", flagLicense)
	}
	// Non-explicit should use config
	if flagCI != "gitlab" {
		t.Errorf("CI: got %q, want gitlab (config default)", flagCI)
	}
}

func TestConfigShowCommand(t *testing.T) {
	resetInitFlags()

	rootCmd.SetArgs([]string{"config", "show"})

	// Should not error even without a config file
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("config show failed: %v", err)
	}
}

func TestConfigPathCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"config", "path"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("config path failed: %v", err)
	}
}

func TestInitWithConfigFile(t *testing.T) {
	resetInitFlags()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	cfg := &config.Config{
		License:     "apache2",
		NoChangelog: config.BoolPtr(true),
		AgentMD:     "none",
	}
	if err := config.Save(cfg, cfgPath); err != nil {
		t.Fatalf("Save config: %v", err)
	}

	loaded, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("Load config: %v", err)
	}

	// Reset cobra's Changed state
	initCmd.ParseFlags([]string{})
	applyConfigDefaults(initCmd, loaded)

	if flagLicense != "apache2" {
		t.Errorf("License from config: got %q, want apache2", flagLicense)
	}
	if !flagNoChangelog {
		t.Error("NoChangelog from config: expected true")
	}
	if flagAgentMD != "none" {
		t.Errorf("AgentMD from config: got %q, want none", flagAgentMD)
	}
}
