package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ariel-frischer/base-cli/internal/config"
	"github.com/spf13/cobra"
)

var configForce bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage user-level defaults (~/.config/base-cli/config.yaml)",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default config file",
	RunE:  runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show resolved configuration",
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a config value. Valid keys:
  host            Module host: github.com, gitlab.com, etc.
  git_user        Git username (default: auto-detected from gh/git)
  author          Default author name
  license         Default license: mit, apache2, none
  ci              Default CI: github, gitlab, both
  layout          Default layout: both, cli, lib
  agent_md        Default agent docs: both, claude, agents, none
  no_git_init     Skip git init: true/false
  no_goreleaser   Skip goreleaser: true/false
  no_community    Skip community files: true/false
  no_changelog    Skip changelog files: true/false
  no_config       Skip config package: true/false
  todo            Include TODO.md: true/false`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config in $EDITOR",
	RunE:  runConfigEdit,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.DefaultPath())
	},
}

func init() {
	configInitCmd.Flags().BoolVar(&configForce, "force", false, "overwrite existing config")
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configPathCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	path := config.DefaultPath()
	if _, err := os.Stat(path); err == nil && !configForce {
		return fmt.Errorf("%s already exists (use --force to overwrite)", path)
	}

	if err := os.MkdirAll(config.DefaultDir(), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	content := `# base-cli configuration — user-level defaults for "base-cli init"
# All fields are optional. CLI flags always override these values.

# host: github.com       # Module host — also used for repo URL
# git_user: yourname     # Git username (default: auto-detected from gh/git)
# author: Your Name
# license: mit           # mit, apache2, none
# ci: both               # github, gitlab, both
# layout: both           # both, cli, lib
# agent_md: both         # both, claude, agents, none
# no_git_init: false
# no_goreleaser: false
# no_community: false
# no_changelog: false
# no_config: false
# todo: false
`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	success("Created %s", fileRef(path))
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		return err
	}

	printCfgRow("host", stringOrDefault(cfg.Host, "github.com"), cfg.Host != "")
	printCfgRow("git_user", stringOrDefault(cfg.GitUser, "(auto-detect)"), cfg.GitUser != "")
	printCfgRow("author", cfg.Author, cfg.Author != "")
	printCfgRow("license", stringOrDefault(cfg.License, "mit"), cfg.License != "")
	printCfgRow("ci", stringOrDefault(cfg.CI, "both"), cfg.CI != "")
	printCfgRow("layout", stringOrDefault(cfg.Layout, "both"), cfg.Layout != "")
	printCfgRow("agent_md", stringOrDefault(cfg.AgentMD, "both"), cfg.AgentMD != "")
	printCfgRow("no_git_init", fmt.Sprintf("%v", config.BoolVal(cfg.NoGitInit, false)), cfg.NoGitInit != nil)
	printCfgRow("no_goreleaser", fmt.Sprintf("%v", config.BoolVal(cfg.NoGoreleaser, false)), cfg.NoGoreleaser != nil)
	printCfgRow("no_community", fmt.Sprintf("%v", config.BoolVal(cfg.NoCommunity, false)), cfg.NoCommunity != nil)
	printCfgRow("no_changelog", fmt.Sprintf("%v", config.BoolVal(cfg.NoChangelog, false)), cfg.NoChangelog != nil)
	printCfgRow("no_config", fmt.Sprintf("%v", config.BoolVal(cfg.NoConfig, false)), cfg.NoConfig != nil)
	printCfgRow("todo", fmt.Sprintf("%v", config.BoolVal(cfg.Todo, false)), cfg.Todo != nil)
	return nil
}

func printCfgRow(key, value string, isCustom bool) {
	source := "default"
	if isCustom {
		source = "config"
	}
	fmt.Printf("%-18s %-20s (%s)\n", highlight(key), value, source)
}

func stringOrDefault(val, def string) string {
	if val != "" {
		return val
	}
	return def
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	path := config.DefaultPath()

	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	switch key {
	case "host":
		cfg.Host = value
	case "git_user":
		cfg.GitUser = value
	case "author":
		cfg.Author = value
	case "license":
		switch value {
		case "mit", "apache2", "none":
		default:
			return fmt.Errorf("invalid license %q: must be mit, apache2, or none", value)
		}
		cfg.License = value
	case "ci":
		switch value {
		case "github", "gitlab", "both":
		default:
			return fmt.Errorf("invalid ci %q: must be github, gitlab, or both", value)
		}
		cfg.CI = value
	case "layout":
		switch value {
		case "both", "cli", "lib":
		default:
			return fmt.Errorf("invalid layout %q: must be both, cli, or lib", value)
		}
		cfg.Layout = value
	case "agent_md":
		switch value {
		case "both", "claude", "agents", "none":
		default:
			return fmt.Errorf("invalid agent_md %q: must be both, claude, agents, or none", value)
		}
		cfg.AgentMD = value
	case "no_git_init":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("no_git_init expects true/false, got %q", value)
		}
		cfg.NoGitInit = config.BoolPtr(b)
	case "no_goreleaser":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("no_goreleaser expects true/false, got %q", value)
		}
		cfg.NoGoreleaser = config.BoolPtr(b)
	case "no_community":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("no_community expects true/false, got %q", value)
		}
		cfg.NoCommunity = config.BoolPtr(b)
	case "no_changelog":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("no_changelog expects true/false, got %q", value)
		}
		cfg.NoChangelog = config.BoolPtr(b)
	case "no_config":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("no_config expects true/false, got %q", value)
		}
		cfg.NoConfig = config.BoolPtr(b)
	case "todo":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("todo expects true/false, got %q", value)
		}
		cfg.Todo = config.BoolPtr(b)
	default:
		validKeys := strings.Join([]string{
			"host", "git_user", "author", "license", "ci", "layout", "agent_md",
			"no_git_init", "no_goreleaser", "no_community", "no_changelog",
			"no_config", "todo",
		}, ", ")
		return fmt.Errorf("unknown key %q\nvalid keys: %s", key, validKeys)
	}

	if err := config.Save(cfg, path); err != nil {
		return err
	}
	success("Set %s = %s", highlight(key), highlight(value))
	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	path := config.DefaultPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		warn("%s not found — run `base-cli config init` first", fileRef(path))
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
