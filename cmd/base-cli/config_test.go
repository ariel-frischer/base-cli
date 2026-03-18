package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ariel-frischer/base-cli/internal/config"
	"github.com/spf13/cobra"
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

// newTestInitCmd creates a fresh cobra command with the same flags as initCmd,
// avoiding shared Changed state across tests.
func newTestInitCmd() (*cobra.Command, *string, *string, *string, *string, *bool, *bool, *bool, *bool, *string) {
	var license, ci, layout, author string
	var noGitInit, noGoreleaser, noCommunity, noChangelog bool
	var agentMD string

	cmd := &cobra.Command{Use: "init", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	cmd.Flags().StringVar(&author, "author", "", "")
	cmd.Flags().StringVar(&license, "license", "mit", "")
	cmd.Flags().StringVar(&ci, "ci", "both", "")
	cmd.Flags().StringVar(&layout, "layout", "both", "")
	cmd.Flags().StringVar(&agentMD, "agent-md", "both", "")
	cmd.Flags().BoolVar(&noGitInit, "no-git-init", false, "")
	cmd.Flags().BoolVar(&noGoreleaser, "no-goreleaser", false, "")
	cmd.Flags().BoolVar(&noCommunity, "no-community", false, "")
	cmd.Flags().BoolVar(&noChangelog, "no-changelog", false, "")

	return cmd, &license, &ci, &layout, &author, &noGitInit, &noGoreleaser, &noCommunity, &noChangelog, &agentMD
}

func TestApplyConfigDefaults(t *testing.T) {
	cmd, license, ci, layout, _, _, noGoreleaser, noCommunity, _, agentMD := newTestInitCmd()

	userCfg := &config.Config{
		License:      "apache2",
		CI:           "github",
		Layout:       "cli",
		AgentMD:      "none",
		NoGoreleaser: config.BoolPtr(true),
		NoCommunity:  config.BoolPtr(true),
	}

	applyConfigDefaults(cmd, userCfg)

	if *license != "apache2" {
		t.Errorf("License: got %q, want apache2", *license)
	}
	if *ci != "github" {
		t.Errorf("CI: got %q, want github", *ci)
	}
	if *layout != "cli" {
		t.Errorf("Layout: got %q, want cli", *layout)
	}
	if *agentMD != "none" {
		t.Errorf("AgentMD: got %q, want none", *agentMD)
	}
	if !*noGoreleaser {
		t.Error("NoGoreleaser: expected true")
	}
	if !*noCommunity {
		t.Error("NoCommunity: expected true")
	}
}

func TestApplyConfigDefaultsFlagOverride(t *testing.T) {
	cmd, license, ci, _, _, _, _, _, _, _ := newTestInitCmd()

	userCfg := &config.Config{
		License: "apache2",
		CI:      "gitlab",
	}

	// Parse with explicit --license flag so cobra marks it as Changed
	_ = cmd.ParseFlags([]string{"--license", "mit"})
	applyConfigDefaults(cmd, userCfg)

	// Explicit flag should win over config
	if *license != "mit" {
		t.Errorf("License: got %q, want mit (explicit flag should override config)", *license)
	}
	// Non-explicit should use config
	if *ci != "gitlab" {
		t.Errorf("CI: got %q, want gitlab (config default)", *ci)
	}
}

func TestConfigShowCommand(t *testing.T) {
	resetInitFlags()

	rootCmd.SetArgs([]string{"config", "show"})

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
	cmd, license, _, _, _, _, _, _, noChangelog, agentMD := newTestInitCmd()

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

	applyConfigDefaults(cmd, loaded)

	if *license != "apache2" {
		t.Errorf("License from config: got %q, want apache2", *license)
	}
	if !*noChangelog {
		t.Error("NoChangelog from config: expected true")
	}
	if *agentMD != "none" {
		t.Errorf("AgentMD from config: got %q, want none", *agentMD)
	}
}
