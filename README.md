<div align="center">

<pre>
██████╗  █████╗ ███████╗███████╗     ██████╗██╗     ██╗
██╔══██╗██╔══██╗██╔════╝██╔════╝    ██╔════╝██║     ██║
██████╔╝███████║███████╗█████╗█████╗██║     ██║     ██║
██╔══██╗██╔══██║╚════██║██╔══╝╚════╝██║     ██║     ██║
██████╔╝██║  ██║███████║███████╗    ╚██████╗███████╗██║
╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝     ╚═════╝╚══════╝╚═╝
</pre>

**Go CLI Project Scaffold Generator**

[![CI](https://github.com/ariel-frischer/base-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/ariel-frischer/base-cli/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/ariel-frischer/base-cli)](https://github.com/ariel-frischer/base-cli/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Generate complete, ready-to-build Go CLI projects with best practices baked in.

</div>

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/base-cli/main/install.sh | sh
```

**Go install**:

```bash
go install github.com/ariel-frischer/base-cli@latest
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

- **Cobra CLI** with root, version, and completion commands
- **Version info** via ldflags (version, commit, build date)
- **Makefile** with build, test, lint, format, release targets
- **goreleaser** config for multi-platform releases
- **CI pipeline** (GitHub Actions and/or GitLab CI)
- **Shell installer** (`install.sh`) with checksum verification
- **Uninstaller** (`uninstall.sh`)
- **Release script** with pre-flight checks
- **CHANGELOG.yaml** + `.chlog.yaml` (ready for [chlog](https://github.com/ariel-frischer/chlog))
- **README.md**, **CLAUDE.md**, **LICENSE**, **.gitignore**

## Usage

```bash
base-cli init <project-name> [flags]
  --module <path>       Go module path (default: github.com/<git-user>/<name>)
  --description <text>  One-line project description
  --author <name>       Author name (default: git config user.name)
  --license mit|apache2|none  (default: mit)
  --ci github|gitlab|both     (default: github)
  --dir <path>          Output directory (default: ./<name>)
  --no-git-init         Skip git init

base-cli version [--plain]
base-cli uninstall [--yes]
base-cli completion [bash|zsh|fish|powershell]
```

Interactive prompts for `--module` and `--description` if not provided (requires TTY).

## Example

```bash
# Non-interactive
base-cli init rentalot-cli \
  --module github.com/ariel-frischer/rentalot-cli \
  --description "CLI for managing Rentalot properties" \
  --ci both

# Interactive (prompts for module and description)
base-cli init my-tool
```

## Generated Project Structure

```
my-project/
  cmd/my-project/
    main.go, root.go, version.go, ui.go
  internal/version/
    version.go, version_test.go
  scripts/release.sh
  .github/workflows/ci.yml     # if --ci github|both
  .gitlab-ci.yml                # if --ci gitlab|both
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

## Development

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run linters
make format         # Format code
```

## License

[MIT](LICENSE)
