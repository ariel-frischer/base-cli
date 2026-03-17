<div align="center">

<pre>
██████╗  █████╗ ███████╗███████╗     ██████╗██╗     ██╗
██╔══██╗██╔══██╗██╔════╝██╔════╝    ██╔════╝██║     ██║
██████╔╝███████║███████╗█████╗█████╗██║     ██║     ██║
██╔══██╗██╔══██║╚════██║██╔══╝╚════╝██║     ██║     ██║
██████╔╝██║  ██║███████║███████╗    ╚██████╗███████╗██║
╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝     ╚═════╝╚══════╝╚═╝
</pre>

**Go Project Scaffold Generator**

[![CI](https://github.com/ariel-frischer/base-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/ariel-frischer/base-cli/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/ariel-frischer/base-cli)](https://github.com/ariel-frischer/base-cli/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Scaffold Go projects with a modular CLI + importable library structure out of the box.

</div>

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/install.sh | sh
```

**Go install**:

```bash
go install github.com/ariel-frischer/base-cli/cmd/base-cli@latest
```

**Library only** (use the scaffold engine programmatically):

```go
import "github.com/ariel-frischer/base-cli/pkg/scaffold"
```

## Quickstart

**1. Install the skill** so your coding agent knows how to use base-cli:

```bash
npx skills add ariel-frischer/base-cli
```

**2. Ask your agent** to scaffold whatever you need:

> "Use base-cli to generate a Go CLI project called my-tool that does X"

Your agent will pick the right flags, layout, and options for you.

**Or run it directly:**

```bash
base-cli init my-project
cd my-project
make build
./bin/my-project version
```

## What You Get

Every generated project includes:

- **Flexible layout** — Choose `--layout both` (default) for a CLI that's also importable as a Go library, `cli` for a standalone binary, or `lib` for a pure library with no CLI
- **Cobra CLI** with root, version, and completion commands (cli/both layouts)
- **Public `pkg/` library** importable by other Go projects (lib/both layouts)
- **Version info** via ldflags (version, commit, build date)
- **Makefile** with build, test, lint, format, release targets
- **goreleaser** config for multi-platform releases
- **CI pipeline** (GitHub Actions and/or GitLab CI)
- **Shell installer** (`install.sh`) with checksum verification
- **Uninstaller** (`uninstall.sh`)
- **Release script** with pre-flight checks
- **CHANGELOG.yaml** + **CHANGELOG.md** + `.chlog.yaml` — optional, powered by [chlog](https://github.com/ariel-frischer/chlog) (skip with `--no-changelog`)
- **Changelog CI gate** — validates changelog in CI when enabled (graceful skip if chlog not installed)
- **TODO.md** with MVP, stretch goals, and tech debt sections
- **testdata/** directory with sample fixture in `pkg/` (encourages fixture-based testing)
- **assets/** directory for demo content (GIFs, screenshots, etc.)
- **AI Agent Skill** — `.skills/default/SKILL.md` + README install instructions ([Agent Skills standard](https://agentskills.io))
- **README.md**, **CLAUDE.md**, **LICENSE**, **.gitignore**

## Usage

```bash
base-cli init <project-name> [module] [flags]
  -d, --description <text>  One-line project description
  --author <name>       Author name (default: git config user.name)
  --license mit|apache2|none      (default: mit)
  --ci github|gitlab|both         (default: both)
  --layout both|cli|lib           (default: both)
  --dir <path>          Output directory (default: ./<name>)
  --no-git-init         Skip git init
  --no-goreleaser       Skip goreleaser config and release workflow
  --no-community        Skip community files (issue templates, PR template, etc.)
  --no-changelog        Skip changelog files (CHANGELOG.yaml, CHANGELOG.md, .chlog.yaml)

base-cli config init [--force]    Set up ~/.config/base-cli/config.yaml
base-cli config show              Show resolved configuration
base-cli config set <key> <value> Set a config value
base-cli config edit              Open config in $EDITOR
base-cli config path              Print config file path
base-cli version [--plain]
base-cli uninstall [--yes]
base-cli completion [bash|zsh|fish|powershell]
```

Interactive prompts for `module` and `--description` if not provided (requires TTY).

## Configuration

Set user-level defaults so you don't have to pass the same flags every time:

```bash
# Create a config file at ~/.config/base-cli/config.yaml
base-cli config init

# Set defaults
base-cli config set license apache2
base-cli config set ci github
base-cli config set layout cli
base-cli config set no_goreleaser true
base-cli config set agent_md none
base-cli config set author "Your Name"

# View current config
base-cli config show

# Edit in $EDITOR
base-cli config edit

# Print config path
base-cli config path
```

CLI flags always override config values. For example, if your config has `license: apache2` but you run `base-cli init foo --license mit`, MIT wins.

### Config file

Located at `~/.config/base-cli/config.yaml`:

```yaml
# base-cli configuration — user-level defaults for "base-cli init"
# All fields are optional. CLI flags always override these values.

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

### Config keys

| Key | Type | Values | Default |
|-----|------|--------|---------|
| `author` | string | any | git config user.name |
| `license` | string | `mit`, `apache2`, `none` | `mit` |
| `ci` | string | `github`, `gitlab`, `both` | `both` |
| `layout` | string | `both`, `cli`, `lib` | `both` |
| `agent_md` | string | `both`, `claude`, `agents`, `none` | `both` |
| `no_git_init` | bool | `true`, `false` | `false` |
| `no_goreleaser` | bool | `true`, `false` | `false` |
| `no_community` | bool | `true`, `false` | `false` |
| `no_changelog` | bool | `true`, `false` | `false` |

### Layout Options

| Layout | What's generated | Use case |
|--------|-----------------|----------|
| `both` (default) | `cmd/` + `pkg/` + `internal/version/` | CLI tool that's also importable as a library |
| `cli` | `cmd/` + `internal/` | Standalone CLI tool |
| `lib` | `pkg/` | Pure library (no cobra, no build targets) |

## Example

```bash
# CLI + library (default)
base-cli init my-tool github.com/alice/my-tool \
  -d "A cool tool"

# Pure library
base-cli init my-lib github.com/alice/my-lib \
  -d "A utility library" \
  --layout lib

# CLI only
base-cli init my-cli --layout cli

# Interactive (prompts for module and description)
base-cli init my-tool
```

## Generated Project Structure

**`--layout both`** (default):

```
my-project/
  cmd/my-project/
    main.go, root.go, version.go, ui.go
  pkg/myproject/
    doc.go
    testdata/sample.yaml
  internal/version/
    version.go, version_test.go
  scripts/release.sh
  assets/.gitkeep
  .github/workflows/ci.yml     # if --ci github|both
  .gitlab-ci.yml                # if --ci gitlab|both
  .skills/default/SKILL.md
  .goreleaser.yaml
  Makefile
  install.sh
  uninstall.sh
  .gitignore
  go.mod
  README.md
  CLAUDE.md
  TODO.md
  LICENSE
  CHANGELOG.yaml                # if --no-changelog not set
  CHANGELOG.md                  # if --no-changelog not set
  .chlog.yaml                   # if --no-changelog not set
```

**`--layout lib`** generates only:

```
my-lib/
  pkg/mylib/
    doc.go
    testdata/sample.yaml
  assets/.gitkeep
  .skills/default/SKILL.md
  .github/workflows/ci.yml
  Makefile                      # test, lint, format only
  go.mod                        # no cobra dependency
  .gitignore
  README.md
  CLAUDE.md
  TODO.md
  LICENSE
  CHANGELOG.yaml                # if --no-changelog not set
  CHANGELOG.md                  # if --no-changelog not set
  .chlog.yaml                   # if --no-changelog not set
```

### AI Agent Skill

base-cli ships a [SKILL.md](.skills/default/SKILL.md) following the [Agent Skills open standard](https://agentskills.io). Install it so your coding agent knows all commands and options.

**Quick install with [`skills`](https://skills.sh) CLI** (by Vercel Labs):

```bash
npx skills add ariel-frischer/base-cli
```

<details>
<summary><strong>Manual install</strong></summary>

**Claude Code** — Skills live in `~/.claude/skills/` (global) or `.claude/skills/` (project-local).

```bash
# Global — available in all projects
mkdir -p ~/.claude/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o ~/.claude/skills/base-cli/SKILL.md

# Project-local — checked into this repo only
mkdir -p .claude/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o .claude/skills/base-cli/SKILL.md
```

Then use `/base-cli` in conversations.

**OpenCode** — reads skills from `~/.claude/skills/` (global) or `.opencode/skills/` (project-local).

```bash
# Global
mkdir -p ~/.claude/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o ~/.claude/skills/base-cli/SKILL.md

# Project-local
mkdir -p .opencode/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o .opencode/skills/base-cli/SKILL.md
```

**Codex CLI** — reads skills from `~/.codex/skills/` (global) or `.codex/skills/` (project-local).

```bash
# Global
mkdir -p ~/.codex/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o ~/.codex/skills/base-cli/SKILL.md

# Project-local
mkdir -p .codex/skills/base-cli
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/.skills/default/SKILL.md \
  -o .codex/skills/base-cli/SKILL.md
```

Or pass directly: `codex --instructions .skills/default/SKILL.md`

</details>

## Development

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run linters
make format         # Format code
```

## License

[MIT](LICENSE)
