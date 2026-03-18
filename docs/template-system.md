# Template System Reference

base-cli uses Go's `text/template` with custom delimiters and functions to generate project files. This document covers the template engine internals — useful if you want to contribute new templates or understand how file generation works.

## Delimiters

Templates use `[%` and `%]` instead of the standard `{{ }}`:

```
[% .ProjectName %]
[% if .HasCLI %]...CLI content...[% end %]
```

This avoids conflicts with:
- goreleaser's `{{ }}` syntax
- Makefile's `$` variables
- Bash's `[[ ]]` conditionals

### Whitespace Control

- `[%-` strips whitespace before the tag
- `-%]` strips whitespace after the tag

```
[%- if .HasCLI %]
content here
[%- end %]
```

## Template Variables

All fields from the `scaffold.Config` struct are available in templates:

### String Variables

| Variable | Example Value | Description |
|----------|---------------|-------------|
| `.ProjectName` | `my-tool` | Project name as provided |
| `.ModulePath` | `github.com/alice/my-tool` | Full Go module path |
| `.BinaryName` | `my-tool` | Binary name (same as ProjectName) |
| `.Description` | `A handy CLI tool` | One-line description |
| `.Author` | `Alice` | Author name |
| `.Year` | `2026` | Copyright year |
| `.GitUser` | `alice` | GitHub/GitLab username |
| `.RepoURL` | `https://github.com/alice/my-tool` | Repository URL |
| `.EnvPrefix` | `MY_TOOL` | Env var prefix (uppercased, hyphens → underscores) |
| `.License` | `mit` | License type: `mit`, `apache2`, `none` |
| `.Layout` | `both` | Layout type: `both`, `cli`, `lib` |
| `.LibPackage` | `mytool` | Go-safe package name (hyphens stripped) |

### Boolean Variables

| Variable | Description |
|----------|-------------|
| `.HasCLI` | true for `both` and `cli` layouts |
| `.HasLib` | true for `both` and `lib` layouts |
| `.CIGitHub` | Generate GitHub Actions workflows |
| `.CIGitLab` | Generate GitLab CI config |
| `.Goreleaser` | Include goreleaser config |
| `.Community` | Include community files |
| `.Changelog` | Include changelog files and CI gate |
| `.Config` | Include `internal/config` package and `config` subcommands |

## Template Functions

In addition to Go's built-in template functions:

| Function | Usage | Result |
|----------|-------|--------|
| `upper` | `[% upper .ProjectName %]` | `MY-TOOL` |
| `lower` | `[% lower .ProjectName %]` | `my-tool` |
| `eq` | `[% if eq .License "mit" %]` | Equality check |
| `ne` | `[% if ne .License "none" %]` | Inequality check |

## Conditional Logic

### Simple conditionals

```
[% if .HasCLI %]
// This content only appears in cli and both layouts
[% end %]
```

### If/else

```
[% if eq .License "mit" %]
MIT License
[% else if eq .License "apache2" %]
Apache 2.0 License
[% else %]
No license
[% end %]
```

### Negation

```
[% if ne .License "none" %]
LICENSE file included
[% end %]
```

## File Naming Conventions

Template files live in `pkg/scaffold/templates/` and are embedded at compile time via `//go:embed`.

### Extension stripping

All `.tmpl` extensions are removed in the output:

```
Makefile.tmpl → Makefile
main.go.tmpl  → main.go
```

### Dynamic path segments

Directory and file names can contain placeholders:

| Placeholder | Replaced With | Example |
|-------------|---------------|---------|
| `{{BinaryName}}` | Config.BinaryName | `cmd/{{BinaryName}}/main.go.tmpl` → `cmd/my-tool/main.go` |
| `{{LibPackage}}` | Config.LibPackage | `pkg/{{LibPackage}}/doc.go.tmpl` → `pkg/mytool/doc.go` |

### Special filename mappings

Some template filenames are transformed to their final output names:

| Template File | Output File | Reason |
|---------------|-------------|--------|
| `gitignore.tmpl` | `.gitignore` | Can't have dotfiles in embed.FS root |
| `gitkeep` | `.gitkeep` | Same reason |
| `goreleaser.yaml.tmpl` | `.goreleaser.yaml` | Dotfile |
| `chlog.yaml.tmpl` | `CHANGELOG.yaml` | Name transformation |
| `chlog-config.yaml.tmpl` | `.chlog.yaml` | Dotfile + rename |
| `LICENSE_mit.tmpl` | `LICENSE` | License variant selection |
| `LICENSE_apache2.tmpl` | `LICENSE` | License variant selection |

### Directory prefix mappings

| Template Directory | Output Directory | Notes |
|--------------------|-----------------|-------|
| `github/` | `.github/` | Dotdir prefix added |
| `gitlab/` | root | `.gitlab-ci.yml` lands at project root |
| `skills/` | `.skills/` | Dotdir prefix added |

## Conditional File Generation

Files are skipped entirely based on config values. The scaffold engine checks conditions before rendering:

- **Layout filtering**: `cmd/`, `internal/` skipped for `lib` layout; `pkg/` skipped for `cli` layout
- **CI filtering**: `.github/` skipped when `CIGitHub` is false; `.gitlab-ci.yml` skipped when `CIGitLab` is false
- **License filtering**: Only the matching `LICENSE_<type>.tmpl` is rendered; both skipped when `License` is `"none"`
- **Goreleaser filtering**: `.goreleaser.yaml`, release workflow, and `scripts/release.sh` skipped when `Goreleaser` is false
- **Community filtering**: Issue templates, PR template, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` skipped when `Community` is false
- **Changelog filtering**: `CHANGELOG.yaml`, `CHANGELOG.md`, `.chlog.yaml` skipped when changelog is disabled
- **Config filtering**: `internal/config/` and `cmd/<tool>/config.go` skipped when `Config` is false (always false for `lib` layout)

## Adding a New Template

1. Create a `.tmpl` file in the appropriate subdirectory of `pkg/scaffold/templates/`
2. Use `[% %]` delimiters for all template logic
3. Add conditional generation logic in `scaffold.go` if the file should only appear for certain configs
4. The file is automatically embedded — no registration step needed

Example template:

```
// Package [% .LibPackage %] provides [% .Description %].
package [% .LibPackage %]

[% if .HasCLI %]
// DefaultBinary is the name of the CLI binary.
const DefaultBinary = "[% .BinaryName %]"
[% end %]
```
