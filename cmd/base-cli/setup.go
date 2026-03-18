package main

import (
	"fmt"
	"strings"

	"github.com/ariel-frischer/base-cli/internal/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive first-time configuration wizard",
	Long:  "Walk through essential settings and save them to ~/.config/base-cli/config.yaml.\nSafe to re-run — existing values become the defaults.",
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	if !isTerminal() {
		return fmt.Errorf("setup requires an interactive terminal")
	}

	path := config.DefaultPath()
	cfg, err := config.Load(path)
	if err != nil {
		cfg = &config.Config{}
	}

	return runSetupWithPath(path, cfg)
}

// runSetupWithPath runs the interactive setup wizard, saving to path.
// Extracted for testability.
func runSetupWithPath(path string, cfg *config.Config) error {
	fmt.Println()
	fmt.Printf("  %s\n\n", highlight("base-cli setup"))

	// --- host ---
	defaultHost := cfg.Host
	if defaultHost == "" {
		defaultHost = detectHost()
	}
	cfg.Host = prompt("host", defaultHost)

	// --- git_user ---
	defaultUser := cfg.GitUser
	if defaultUser == "" {
		defaultUser = detectGitUser()
	}
	cfg.GitUser = prompt("git_user", defaultUser)

	// --- author ---
	defaultAuthor := cfg.Author
	if defaultAuthor == "" {
		defaultAuthor = gitConfigValue("user.name")
	}
	cfg.Author = prompt("author", defaultAuthor)

	// --- license ---
	defaultLicense := cfg.License
	if defaultLicense == "" {
		defaultLicense = "mit"
	}
	cfg.License = promptChoice("license", []string{"mit", "apache2", "none"}, defaultLicense)

	// --- ci ---
	defaultCI := cfg.CI
	if defaultCI == "" {
		defaultCI = "both"
	}
	cfg.CI = promptChoice("ci", []string{"github", "gitlab", "both"}, defaultCI)

	// --- layout ---
	defaultLayout := cfg.Layout
	if defaultLayout == "" {
		defaultLayout = "both"
	}
	cfg.Layout = promptChoice("layout", []string{"both", "cli", "lib"}, defaultLayout)

	// --- summary ---
	fmt.Println()
	fmt.Printf("  %s\n", highlight("Summary"))
	printSetupRow("host", cfg.Host)
	printSetupRow("git_user", cfg.GitUser)
	printSetupRow("author", cfg.Author)
	printSetupRow("license", cfg.License)
	printSetupRow("ci", cfg.CI)
	printSetupRow("layout", cfg.Layout)
	fmt.Println()

	answer := prompt("Save to "+fileRef(path)+"?", "Y")
	if !strings.HasPrefix(strings.ToLower(answer), "y") {
		warn("  Setup cancelled.")
		return nil
	}

	if err := config.Save(cfg, path); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println()
	success("  Config saved to %s", fileRef(path))
	fmt.Println()
	return nil
}

func printSetupRow(key, value string) {
	fmt.Printf("    %-14s %s\n", highlight(key), value)
}
