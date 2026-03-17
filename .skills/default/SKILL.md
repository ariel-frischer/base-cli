---
name: base-cli
description: >
  Go CLI + library project scaffold generator. Use when generating new Go projects,
  understanding base-cli commands, or using the scaffold engine programmatically.
  Covers init, config, version, uninstall, completion commands, layout options, and template system.
license: MIT
compatibility:
  - Claude Code
  - Cursor
  - Codex
  - Gemini CLI
  - VS Code
metadata:
  author: ariel-frischer
  version: 0.0.1
  tags: scaffold, generator, go, cli, project-template
allowed-tools: Bash Read Write Edit
---

# base-cli

Go project scaffold generator. `base-cli init <name>` → complete, ready-to-build Go project with best practices.

## Commands

```bash
base-cli init <name>                    # Generate a new Go project
base-cli init <name> --layout both      # CLI + library (default)
base-cli init <name> --layout cli       # CLI only (cmd/ + internal/)
base-cli init <name> --layout lib       # Library only (pkg/)
base-cli init <name> --module <path>    # Custom Go module path
base-cli init <name> --description <t>  # One-line project description
base-cli init <name> --author <name>    # Author name (default: git user.name)
base-cli init <name> --license mit      # License: mit|apache2|none (default: mit)
base-cli init <name> --ci both          # CI: github|gitlab|both (default: both)
base-cli init <name> --dir <path>       # Output directory (default: ./<name>)
base-cli init <name> --agent-md both    # AI docs: both|claude|agents|none (default: both)
base-cli init <name> --no-git-init      # Skip git init
base-cli init <name> --no-goreleaser   # Skip goreleaser config and release workflow
base-cli init <name> --no-community    # Skip community files (issue templates, etc.)
base-cli init <name> --no-changelog    # Skip changelog files (CHANGELOG.yaml, .chlog.yaml)
base-cli config init                    # Create ~/.config/base-cli/config.yaml
base-cli config init --force            # Overwrite existing config
base-cli config show                    # Show resolved config (config vs default)
base-cli config set <key> <value>       # Set a config value
base-cli config edit                    # Open config in $EDITOR
base-cli config path                    # Print config file path
base-cli version                        # Show version info
base-cli version --plain                # Plain version string
base-cli uninstall                      # Remove base-cli from system
base-cli uninstall --yes                # Skip confirmation
base-cli completion bash                # Shell completion: bash|zsh|fish|powershell
```

Interactive prompts for `--module` and `--description` if not provided (requires TTY).

## Configuration

User-level defaults at `~/.config/base-cli/config.yaml`. CLI flags always override config values.

```yaml
author: Your Name
license: apache2        # mit, apache2, none
ci: github              # github, gitlab, both
layout: both            # both, cli, lib
agent_md: both          # both, claude, agents, none
no_git_init: false
no_goreleaser: false
no_community: false
no_changelog: false
```

Config keys for `base-cli config set`: `author`, `license`, `ci`, `layout`, `agent_md`, `no_git_init`, `no_goreleaser`, `no_community`, `no_changelog`.

## Layout Options

| Layout | Generated | Use case |
|--------|-----------|----------|
| `both` (default) | `cmd/` + `pkg/` + `internal/version/` | CLI tool that's also importable as a library |
| `cli` | `cmd/` + `internal/` | Standalone CLI tool |
| `lib` | `pkg/` | Pure library (no cobra, no build targets) |

## Generated Project Structure

Every project includes: Makefile, goreleaser config, CI pipelines (GitHub Actions + GitLab CI by default), shell installer/uninstaller, release script, TODO.md, assets/ directory, README.md, CLAUDE.md, LICENSE, .gitignore. Optional: CHANGELOG.yaml + CHANGELOG.md + .chlog.yaml + changelog CI gate (powered by [chlog](https://github.com/ariel-frischer/chlog), skip with `--no-changelog`).

CLI layouts add: Cobra CLI with root, version, and completion commands, version info via ldflags.

Library layouts add: public `pkg/` importable by other Go projects, with `testdata/` directory and sample fixture.

## Go Library

The scaffold engine is importable:

```go
import "github.com/ariel-frischer/base-cli/pkg/scaffold"
```

## Key Design Details

- Templates use `[% %]` delimiters (avoids goreleaser `{{ }}`, Make `$`, bash `[[ ]]` conflicts)
- `embed.FS` bundles templates at compile time — zero runtime dependencies
- `go mod tidy` runs after generation (best-effort, prints note if offline)
- User-level config at `~/.config/base-cli/config.yaml` — YAML, loaded silently (no error if missing), CLI flags always take precedence via `cmd.Flags().Changed()`
