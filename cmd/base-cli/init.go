package main

import (
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
	flagNoConfig     bool
	flagTodo         bool
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
	initCmd.Flags().BoolVar(&flagNoConfig, "no-config", false, "Skip config package and config subcommands (internal/config + cmd config)")
	initCmd.Flags().BoolVar(&flagTodo, "todo", false, "Include TODO.md with MVP/stretch goals/tech debt sections")
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

	// Resolve git user: config > auto-detect
	gitUser := userCfg.GitUser
	if gitUser == "" {
		gitUser = detectGitUser()
	}

	// Resolve host: config > auto-detect
	host := userCfg.Host
	if host == "" {
		host = detectHost()
	}

	// Resolve module path: explicit arg > auto-derive > prompt (TTY)
	var modulePath string
	defaultModule := fmt.Sprintf("%s/%s/%s", host, gitUser, projectName)
	if len(args) >= 2 {
		modulePath = args[1]
	} else if host != "" && gitUser != "" {
		modulePath = defaultModule
	} else if isTTY {
		modulePath = prompt("Go module path", defaultModule)
	} else {
		modulePath = defaultModule
	}

	// Resolve description
	description := flagDescription
	if description == "" {
		if isTTY {
			description = prompt("Project description", "A CLI tool")
		} else {
			description = "A CLI tool"
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
	repoURL := fmt.Sprintf("https://%s/%s/%s", host, gitUser, projectName)
	if ciGitLab && !ciGitHub && host == "github.com" {
		repoURL = fmt.Sprintf("https://gitlab.com/%s/%s", gitUser, projectName)
	}

	cfg := scaffold.Config{
		ProjectName:   projectName,
		ModulePath:    modulePath,
		BinaryName:    projectName,
		Description:   description,
		Author:        author,
		Year:          fmt.Sprintf("%d", time.Now().Year()),
		GitUser:       gitUser,
		RepoURL:       repoURL,
		CIGitHub:      ciGitHub,
		CIGitLab:      ciGitLab,
		EnvPrefix:     envPrefix,
		License:       flagLicense,
		Layout:        flagLayout,
		HasCLI:        hasCLI,
		HasLib:        hasLib,
		LibPackage:    libPackage,
		Goreleaser:    !flagNoGoreleaser,
		Community:     !flagNoCommunity,
		Changelog:     !flagNoChangelog,
		Config:        !flagNoConfig && hasCLI,
		Todo:          flagTodo,
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
	setBoolIfUnchanged("no-config", cfg.NoConfig)
	setBoolIfUnchanged("todo", cfg.Todo)
}
