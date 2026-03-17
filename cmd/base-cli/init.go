package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ariel-frischer/base-cli/internal/config"
	"github.com/ariel-frischer/base-cli/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	flagDescription  string
	flagAuthor       string
	flagLicense      string
	flagCI           string
	flagLayout       string
	flagDir          string
	flagNoGitInit    bool
	flagNoGoreleaser bool
	flagNoCommunity  bool
	flagNoChangelog  bool
	flagAgentMD      string
)

var initCmd = &cobra.Command{
	Use:   "init <project-name> [module]",
	Short: "Generate a new Go project",
	Long:  "Scaffold a complete, ready-to-build Go project with best practices.\n\nLayout options:\n  both   CLI + library (cmd/ + pkg/) — default\n  cli    CLI only (cmd/ + internal/)\n  lib    Library only (pkg/)",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringVarP(&flagDescription, "description", "d", "", "One-line project description")
	initCmd.Flags().StringVar(&flagAuthor, "author", "", "Author name (default: git config user.name)")
	initCmd.Flags().StringVar(&flagLicense, "license", "mit", "License type: mit, apache2, none")
	initCmd.Flags().StringVar(&flagCI, "ci", "both", "CI provider: github, gitlab, both")
	initCmd.Flags().StringVar(&flagLayout, "layout", "both", "Project layout: both (cli+lib), cli, lib")
	initCmd.Flags().StringVar(&flagDir, "dir", "", "Output directory (default: ./<name>)")
	initCmd.Flags().BoolVar(&flagNoGitInit, "no-git-init", false, "Skip git init")
	initCmd.Flags().BoolVar(&flagNoGoreleaser, "no-goreleaser", false, "Skip goreleaser config and release workflow")
	initCmd.Flags().BoolVar(&flagNoCommunity, "no-community", false, "Skip community files (issue templates, PR template, CONTRIBUTING, CODE_OF_CONDUCT)")
	initCmd.Flags().BoolVar(&flagNoChangelog, "no-changelog", false, "Skip changelog files (CHANGELOG.yaml, CHANGELOG.md, .chlog.yaml, CI changelog gate)")
	initCmd.Flags().StringVar(&flagAgentMD, "agent-md", "both", "AI agent docs: both, claude, agents, none")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	isTTY := isTerminal()

	// Load user-level config defaults.
	userCfg, err := config.Load(config.DefaultPath())
	if err != nil {
		warn("failed to load config: %v", err)
		userCfg = &config.Config{}
	}
	applyConfigDefaults(cmd, userCfg)

	// Resolve author
	author := flagAuthor
	if author == "" {
		author = gitConfigValue("user.name")
	}

	// Resolve git user from git config
	gitUser := gitConfigValue("user.name")
	// Try to extract from remote URL or use login name
	if u := gitHubUser(); u != "" {
		gitUser = u
	}

	// Resolve module path
	var modulePath string
	if len(args) >= 2 {
		modulePath = args[1]
	} else if isTTY {
		defaultModule := fmt.Sprintf("github.com/%s/%s", gitUser, projectName)
		modulePath = prompt("Go module path", defaultModule)
	} else {
		return fmt.Errorf("module path is required in non-interactive mode: base-cli init <name> <module>")
	}

	// Resolve description
	description := flagDescription
	if description == "" {
		if isTTY {
			description = prompt("Project description", "A CLI tool")
		} else {
			return fmt.Errorf("--description is required in non-interactive mode")
		}
	}

	// Resolve output directory
	destDir := flagDir
	if destDir == "" {
		destDir = filepath.Join(".", projectName)
	}

	// Validate license
	switch flagLicense {
	case "mit", "apache2", "none":
	default:
		return fmt.Errorf("invalid license %q: must be mit, apache2, or none", flagLicense)
	}

	// Validate layout
	hasCLI, hasLib := true, true
	switch flagLayout {
	case "both":
	case "cli":
		hasLib = false
	case "lib":
		hasCLI = false
	default:
		return fmt.Errorf("invalid layout %q: must be both, cli, or lib", flagLayout)
	}

	// Validate CI
	ciGitHub, ciGitLab := false, false
	switch flagCI {
	case "github":
		ciGitHub = true
	case "gitlab":
		ciGitLab = true
	case "both":
		ciGitHub = true
		ciGitLab = true
	default:
		return fmt.Errorf("invalid CI provider %q: must be github, gitlab, or both", flagCI)
	}

	// Validate agent-md
	agentMDClaude, agentMDAgents := false, false
	switch flagAgentMD {
	case "both":
		agentMDClaude, agentMDAgents = true, true
	case "claude":
		agentMDClaude = true
	case "agents":
		agentMDAgents = true
	case "none":
	default:
		return fmt.Errorf("invalid agent-md %q: must be both, claude, agents, or none", flagAgentMD)
	}

	// Build env prefix: MY-CLI -> MY_CLI
	envPrefix := strings.ToUpper(strings.ReplaceAll(projectName, "-", "_"))

	// Build Go-safe package name: my-tool -> mytool
	libPackage := strings.ReplaceAll(projectName, "-", "")

	// Determine repo URL
	repoURL := fmt.Sprintf("https://github.com/%s/%s", gitUser, projectName)
	if ciGitLab && !ciGitHub {
		repoURL = fmt.Sprintf("https://gitlab.com/%s/%s", gitUser, projectName)
	}

	cfg := scaffold.Config{
		ProjectName: projectName,
		ModulePath:  modulePath,
		BinaryName:  projectName,
		Description: description,
		Author:      author,
		Year:        fmt.Sprintf("%d", time.Now().Year()),
		GitUser:     gitUser,
		RepoURL:     repoURL,
		CIGitHub:    ciGitHub,
		CIGitLab:    ciGitLab,
		EnvPrefix:   envPrefix,
		License:     flagLicense,
		Layout:      flagLayout,
		HasCLI:      hasCLI,
		HasLib:      hasLib,
		LibPackage:  libPackage,
		Goreleaser:  !flagNoGoreleaser,
		Community:   !flagNoCommunity,
		Changelog:     !flagNoChangelog,
		AgentMDClaude: agentMDClaude,
		AgentMDAgents: agentMDAgents,
	}

	// Check if directory already exists and is non-empty
	if entries, err := os.ReadDir(destDir); err == nil && len(entries) > 0 {
		return fmt.Errorf("directory %s already exists and is not empty", destDir)
	}

	fmt.Printf("\nScaffolding %s into %s...\n\n", highlight(projectName), fileRef(destDir))

	if err := scaffold.Generate(cfg, destDir); err != nil {
		return fmt.Errorf("generating scaffold: %w", err)
	}

	// Run go mod tidy (best effort)
	if err := runInDir(destDir, "go", "mod", "tidy"); err != nil {
		warn("go mod tidy failed (you may be offline): %v", err)
	}

	// Git init
	if !flagNoGitInit {
		if err := runInDir(destDir, "git", "init"); err != nil {
			warn("git init failed: %v", err)
		} else if err := runInDir(destDir, "git", "add", "."); err != nil {
			warn("git add failed: %v", err)
		} else if err := runInDir(destDir, "git", "commit", "-m", "Initial scaffold from base-cli"); err != nil {
			warn("git commit failed: %v", err)
		}
	}

	fmt.Println()
	success("Project %s created successfully!", projectName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", destDir)
	if hasCLI {
		fmt.Println("  make build")
		fmt.Printf("  ./bin/%s version\n", projectName)
	} else {
		fmt.Println("  make test")
	}
	fmt.Println()

	return nil
}

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func prompt(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return defaultVal
	}
	return answer
}

func gitConfigValue(key string) string {
	out, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitHubUser() string {
	// Try to extract from GitHub CLI
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	// Fall back to git config
	email := gitConfigValue("user.email")
	if parts := strings.Split(email, "@"); len(parts) > 0 {
		// Common pattern: user@github.com or user+noreply@github.com
		name := strings.Split(parts[0], "+")[0]
		if name != "" {
			return name
		}
	}
	return gitConfigValue("user.name")
}

func runInDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// applyConfigDefaults sets flag values from user config when the flag was not
// explicitly provided on the command line. CLI flags always take precedence.
func applyConfigDefaults(cmd *cobra.Command, cfg *config.Config) {
	setIfUnchanged := func(name, val string) {
		if val != "" && !cmd.Flags().Changed(name) {
			_ = cmd.Flags().Set(name, val)
		}
	}

	setIfUnchanged("author", cfg.Author)
	setIfUnchanged("license", cfg.License)
	setIfUnchanged("ci", cfg.CI)
	setIfUnchanged("layout", cfg.Layout)
	setIfUnchanged("agent-md", cfg.AgentMD)

	setBoolIfUnchanged := func(name string, ptr *bool) {
		if ptr != nil && !cmd.Flags().Changed(name) {
			_ = cmd.Flags().Set(name, fmt.Sprintf("%v", *ptr))
		}
	}

	setBoolIfUnchanged("no-git-init", cfg.NoGitInit)
	setBoolIfUnchanged("no-goreleaser", cfg.NoGoreleaser)
	setBoolIfUnchanged("no-community", cfg.NoCommunity)
	setBoolIfUnchanged("no-changelog", cfg.NoChangelog)
}
