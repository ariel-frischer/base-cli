---
name: base-cli
description: >
  Go CLI + library project scaffold generator. Use when generating new Go projects,
  understanding base-cli commands, or using the scaffold engine programmatically.
  Covers init, version, uninstall, completion commands, layout options, and template system.
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
base-cli init <name> --ci github        # CI: github|gitlab|both (default: github)
base-cli init <name> --dir <path>       # Output directory (default: ./<name>)
base-cli init <name> --no-git-init      # Skip git init
base-cli version                        # Show version info
base-cli version --plain                # Plain version string
base-cli uninstall                      # Remove base-cli from system
base-cli uninstall --yes                # Skip confirmation
base-cli completion bash                # Shell completion: bash|zsh|fish|powershell
```

Interactive prompts for `--module` and `--description` if not provided (requires TTY).

## Layout Options

| Layout | Generated | Use case |
|--------|-----------|----------|
| `both` (default) | `cmd/` + `pkg/` + `internal/version/` | CLI tool that's also importable as a library |
| `cli` | `cmd/` + `internal/` | Standalone CLI tool |
| `lib` | `pkg/` | Pure library (no cobra, no build targets) |

## Generated Project Structure

Every project includes: Makefile, goreleaser config, CI pipeline, shell installer/uninstaller, release script, CHANGELOG.yaml, .chlog.yaml, README.md, CLAUDE.md, LICENSE, .gitignore.

CLI layouts add: Cobra CLI with root, version, and completion commands, version info via ldflags.

Library layouts add: public `pkg/` importable by other Go projects.

## Go Library

The scaffold engine is importable:

```go
import "github.com/ariel-frischer/base-cli/pkg/scaffold"
```

## Key Design Details

- Templates use `[% %]` delimiters (avoids goreleaser `{{ }}`, Make `$`, bash `[[ ]]` conflicts)
- `embed.FS` bundles templates at compile time — zero runtime dependencies
- `go mod tidy` runs after generation (best-effort, prints note if offline)
