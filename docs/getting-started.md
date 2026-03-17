# Getting Started with base-cli

A step-by-step tutorial covering installation, project generation, layout options, the library API, CI/CD, and every customization flag.

## Install

**Shell installer** (recommended):

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/install.sh | sh
```

**Go install**:

```bash
go install github.com/ariel-frischer/base-cli/cmd/base-cli@latest
```

Verify it works:

```bash
base-cli version
```

## Your First Project

```bash
base-cli init my-tool
```

If you're in an interactive terminal, base-cli prompts for the Go module path and description. Otherwise, pass them explicitly:

```bash
base-cli init my-tool \
  --module github.com/alice/my-tool \
  --description "A handy CLI tool"
```

This generates a complete Go project with a Cobra CLI, a public library package, CI pipelines, a goreleaser config, an installer script, and more. Build and run it immediately:

```bash
cd my-tool
make build
./bin/my-tool version
```

## Choosing a Layout

The `--layout` flag controls what kind of project gets generated.

### `both` (default) — CLI + Library

```bash
base-cli init my-tool --layout both
```

Generates `cmd/my-tool/` (Cobra CLI), `pkg/mytool/` (importable library), and `internal/version/`. This is the right choice when your tool should also be usable as a Go library by other projects.

### `cli` — Standalone CLI

```bash
base-cli init my-tool --layout cli
```

Generates `cmd/my-tool/` and `internal/` only. No `pkg/` directory. Use this for tools that don't need to expose a library API.

### `lib` — Pure Library

```bash
base-cli init my-lib --layout lib
```

Generates `pkg/mylib/` only. No Cobra, no `cmd/`, no build targets, no installer. The Makefile only includes `test`, `lint`, and `format`. The `go.mod` has no CLI dependencies.

### What each layout includes

| Feature | `both` | `cli` | `lib` |
|---------|--------|-------|-------|
| `cmd/` (Cobra CLI) | Yes | Yes | No |
| `pkg/` (library) | Yes | No | Yes |
| `internal/version/` | Yes | Yes | No |
| `install.sh` / `uninstall.sh` | Yes | Yes | No |
| goreleaser config | Yes | Yes | No |
| `make build` target | Yes | Yes | No |
| `make test/lint/format` | Yes | Yes | Yes |
| `scripts/release.sh` | Yes | Yes | No |

## CI Configuration

The `--ci` flag controls which CI pipelines are generated.

```bash
base-cli init my-tool --ci github    # GitHub Actions only
base-cli init my-tool --ci gitlab    # GitLab CI only
base-cli init my-tool --ci both      # Both (default)
```

**GitHub Actions** generates:
- `.github/workflows/ci.yml` — lint, test, build, changelog check
- `.github/workflows/release.yml` — goreleaser-based release automation (when goreleaser is enabled)

**GitLab CI** generates:
- `.gitlab-ci.yml` — lint, test, build stages, plus a release job (when goreleaser is enabled)

The changelog CI gate uses `continue-on-error` (GitHub) / `allow_failure` (GitLab) so it won't block your pipeline — it just nudges you to keep the changelog updated.

## Optional Features

### Goreleaser

Enabled by default. Generates `.goreleaser.yaml`, a release workflow, and `scripts/release.sh` (pre-flight checks, tagging, changelog management).

```bash
base-cli init my-tool --no-goreleaser   # Skip all release tooling
```

### Community Files

Enabled by default. Generates issue templates, a PR template, `CONTRIBUTING.md`, and `CODE_OF_CONDUCT.md`.

```bash
base-cli init my-tool --no-community    # Skip community files
```

### Changelog

Enabled by default. Generates `CHANGELOG.yaml`, `CHANGELOG.md`, and `.chlog.yaml` for use with [chlog](https://github.com/ariel-frischer/chlog). Also adds a changelog check step to CI.

```bash
base-cli init my-tool --no-changelog    # Skip changelog files and CI gate
```

## All Flags Reference

```
base-cli init <project-name> [flags]

  --module <path>         Go module path (default: github.com/<git-user>/<name>)
  --description <text>    One-line project description
  --author <name>         Author name (default: git config user.name)
  --license mit|apache2|none           (default: mit)
  --ci github|gitlab|both             (default: both)
  --layout both|cli|lib               (default: both)
  --dir <path>            Output directory (default: ./<name>)
  --no-git-init           Skip git init and initial commit
  --no-goreleaser         Skip goreleaser config and release workflow
  --no-community          Skip community files (issue templates, PR template, etc.)
  --no-changelog          Skip changelog files and CI changelog gate
  --no-color              Disable colored output (global flag)
```

## Using the Library API

The scaffold engine is a public Go library. You can import it and generate projects programmatically:

```go
package main

import (
    "log"
    "github.com/ariel-frischer/base-cli/pkg/scaffold"
)

func main() {
    cfg := scaffold.Config{
        ProjectName: "my-tool",
        ModulePath:  "github.com/alice/my-tool",
        BinaryName:  "my-tool",
        Description: "A cool CLI tool",
        Author:      "Alice",
        Year:        "2026",
        GitUser:     "alice",
        RepoURL:     "https://github.com/alice/my-tool",
        CIGitHub:    true,
        CIGitLab:    false,
        EnvPrefix:   "MY_TOOL",
        License:     "mit",
        Layout:      "both",
        HasCLI:      true,
        HasLib:      true,
        LibPackage:  "mytool",
        Goreleaser:  true,
        Community:   true,
        Changelog:   true,
    }

    if err := scaffold.Generate(cfg, "./my-tool"); err != nil {
        log.Fatalf("scaffold failed: %v", err)
    }
}
```

`scaffold.Generate` walks the embedded template filesystem, evaluates conditionals, renders templates with your config, and writes the output tree to the destination directory. It returns a single error (wrapped with context for each file that fails).

### Config Fields

| Field | Type | Description |
|-------|------|-------------|
| `ProjectName` | string | Project name (used in directory names, README, etc.) |
| `ModulePath` | string | Full Go module path |
| `BinaryName` | string | Binary name (same as ProjectName) |
| `Description` | string | One-line project description |
| `Author` | string | Author name |
| `Year` | string | Copyright year |
| `GitUser` | string | GitHub/GitLab username |
| `RepoURL` | string | Full repository URL |
| `CIGitHub` | bool | Generate GitHub Actions workflows |
| `CIGitLab` | bool | Generate GitLab CI config |
| `EnvPrefix` | string | Env var prefix (hyphens → underscores, uppercase) |
| `License` | string | `"mit"`, `"apache2"`, or `"none"` |
| `Layout` | string | `"both"`, `"cli"`, or `"lib"` |
| `HasCLI` | bool | true for `both` and `cli` layouts |
| `HasLib` | bool | true for `both` and `lib` layouts |
| `LibPackage` | string | Go-safe package name (hyphens stripped) |
| `Goreleaser` | bool | Include goreleaser config and release tooling |
| `Community` | bool | Include community files |
| `Changelog` | bool | Include changelog files and CI gate ([chlog](https://github.com/ariel-frischer/chlog)) |

## Generated Installer

Every CLI project includes an `install.sh` that users can run:

```bash
curl -fsSL https://raw.githubusercontent.com/<user>/<project>/main/install.sh | sh
```

The installer:
- Detects OS (Linux, macOS, Windows WSL) and architecture (amd64, arm64)
- Downloads the latest release from GitHub (or a specific version via env var)
- Verifies SHA256 checksums when available
- Backs up existing installations (keeps last 3)
- Installs to `~/.local/bin` by default

Customize with environment variables:
- `<ENV_PREFIX>_INSTALL_DIR` — Custom install directory
- `<ENV_PREFIX>_VERSION` — Pin a specific version (e.g., `v1.2.3`)

## Release Workflow

When goreleaser is enabled, the generated project includes a full release pipeline:

1. **`scripts/release.sh`** — Run this locally to cut a release. It:
   - Runs tests and linters
   - Builds the binary
   - Prompts for the next version
   - Updates the changelog
   - Creates a git tag and pushes it

2. **CI release job** — Triggered by the new tag, it runs goreleaser to build multi-platform binaries and create a GitHub/GitLab release.

## Other Commands

```bash
base-cli version              # Version, commit, build date
base-cli version --plain      # Machine-readable version string
base-cli uninstall            # Remove base-cli from your system
base-cli uninstall --yes      # Skip confirmation prompt
base-cli completion bash      # Shell completion (bash|zsh|fish|powershell)
```
