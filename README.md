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

```bash
base-cli init my-project
cd my-project
make build
./bin/my-project version
```

## What You Get

Every generated project includes:

- **Flexible layout** — CLI + library (default), CLI only, or library only
- **Cobra CLI** with root, version, and completion commands (cli/both layouts)
- **Public `pkg/` library** importable by other Go projects (lib/both layouts)
- **Version info** via ldflags (version, commit, build date)
- **Makefile** with build, test, lint, format, release targets
- **goreleaser** config for multi-platform releases
- **CI pipeline** (GitHub Actions and/or GitLab CI)
- **Shell installer** (`install.sh`) with checksum verification
- **Uninstaller** (`uninstall.sh`)
- **Release script** with pre-flight checks
- **CHANGELOG.yaml** + `.chlog.yaml` (ready for [chlog](https://github.com/ariel-frischer/chlog))
- **AI Agent Skill** — `.skills/default/SKILL.md` + README install instructions ([Agent Skills standard](https://agentskills.io))
- **README.md**, **CLAUDE.md**, **LICENSE**, **.gitignore**

## Usage

```bash
base-cli init <project-name> [flags]
  --module <path>       Go module path (default: github.com/<git-user>/<name>)
  --description <text>  One-line project description
  --author <name>       Author name (default: git config user.name)
  --license mit|apache2|none      (default: mit)
  --ci github|gitlab|both         (default: github)
  --layout both|cli|lib           (default: both)
  --dir <path>          Output directory (default: ./<name>)
  --no-git-init         Skip git init

base-cli version [--plain]
base-cli uninstall [--yes]
base-cli completion [bash|zsh|fish|powershell]
```

Interactive prompts for `--module` and `--description` if not provided (requires TTY).

### Layout Options

| Layout | What's generated | Use case |
|--------|-----------------|----------|
| `both` (default) | `cmd/` + `pkg/` + `internal/version/` | CLI tool that's also importable as a library |
| `cli` | `cmd/` + `internal/` | Standalone CLI tool |
| `lib` | `pkg/` | Pure library (no cobra, no build targets) |

## Example

```bash
# CLI + library (default)
base-cli init my-tool \
  --module github.com/alice/my-tool \
  --description "A cool tool"

# Pure library
base-cli init my-lib \
  --module github.com/alice/my-lib \
  --description "A utility library" \
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
  internal/version/
    version.go, version_test.go
  scripts/release.sh
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
  LICENSE
  CHANGELOG.yaml
  .chlog.yaml
```

**`--layout lib`** generates only:

```
my-lib/
  pkg/mylib/
    doc.go
  .skills/default/SKILL.md
  .github/workflows/ci.yml
  Makefile                      # test, lint, format only
  go.mod                        # no cobra dependency
  .gitignore
  README.md
  CLAUDE.md
  LICENSE
  CHANGELOG.yaml
  .chlog.yaml
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
