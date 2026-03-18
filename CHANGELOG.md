# Changelog

All notable changes to base-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [0.1.0] - 2026-03-18

### Added

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

### Changed

- module path is now a positional arg: `base-cli init <name> <module>` — no --module flag needed

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
