# Changelog

All notable changes to base-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Generated CLI projects now include fuller YAML config commands: get, set, toggle, keys, root --config, and <ENV_PREFIX>_CONFIG path override.
- Makefile aliases: `make bin` for build and `make install-global` for go-install in base-cli and generated CLI projects

## [0.1.0] - 2026-03-18

### Added

- todo flag to opt-in to TODO.md generation (off by default)
- SKILL.md for AI agent integration
- Installation instructions in README
- `--no-changelog` flag to skip changelog files (CHANGELOG.yaml, CHANGELOG.md, .chlog.yaml) and CI changelog gate
- `--agent-md` flag to control AI agent doc generation (both, claude, agents, none)
- AGENTS.md template with module path, architecture, dependencies, and file conventions
- User-level config system at `~/.config/base-cli/config.yaml` — set defaults for all init flags
- `base-cli config` command with `init`, `show`, `set`, `edit`, `path` subcommands
- Config values applied as defaults when CLI flags are not explicitly passed
- Shorthand -d flag for --description on init command
- assets/demo.tape for VHS terminal GIF generation
- Config scaffold feature: scaffolded CLI projects now include `internal/config` package (Load/Save/DefaultPath, typed Config struct) and `config` subcommands (init/show/path/edit) by default; disable with `--no-config` flag or `no_config` in `~/.config/base-cli/config.yaml`
- Config fields 'host' and 'git_user' for auto-deriving module path without requiring full URL
- `base-cli setup` interactive first-time configuration wizard
- GitLab user detection via glab CLI alongside existing GitHub detection

### Fixed

- no_config missing from config show, set, and init template
- Make --description flag optional; non-interactive mode now defaults to 'A CLI tool' instead of erroring

### Changed

- Remove CLAUDE.md and TODO.md from scaffolded .gitignore so they are tracked in git
- module path is now a positional arg: `base-cli init <name> <module>` — no --module flag needed
- Non-interactive mode now auto-derives module path from host/git_user/project-name instead of erroring
- Auto-detect host (github.com/gitlab.com) based on available CLI authentication
- Skip Go module path prompt when host and git_user are available (config or auto-detected)

## [0.0.1] - 2026-03-16

### Added

- Initial project scaffolding
- `base-cli init` command to generate Go CLI projects
- Scaffold engine with embed.FS and template walking
- Interactive prompts for module path and description
- Support for MIT, Apache-2.0, and no license
- GitHub Actions and GitLab CI template generation
- Shell installer with checksum verification
- Uninstaller script for generated projects
- `base-cli uninstall` for self-removal
- `base-cli version` with pretty ASCII box display

[Unreleased]: https://github.com/ariel-frischer/base-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ariel-frischer/base-cli/compare/v0.0.1...v0.1.0
