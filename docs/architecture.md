# Architecture

How the base-cli scaffold engine works internally.

## Overview

```
User runs `base-cli init my-tool`
        │
        ▼
┌─────────────────────┐
│  internal/config/   │  Loads ~/.config/base-cli/config.yaml
│  config.go          │  (silent no-op if missing)
└────────┬────────────┘
         │
         ▼
┌─────────────────┐
│  cmd/base-cli/  │  Cobra command parses flags, applies config defaults
│  init.go        │  for unchanged flags, resolves remaining values,
│                 │  builds scaffold.Config, calls scaffold.Generate()
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  pkg/scaffold/  │  Walks embed.FS, evaluates skip conditions,
│  scaffold.go    │  renders templates, writes output tree
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  templates/     │  33 .tmpl files compiled into the binary
│  (embed.FS)     │  via //go:embed at build time
└─────────────────┘
```

## Key Components

### `cmd/base-cli/init.go` — CLI Layer

Responsible for:
- Flag definitions and validation
- Loading user-level config from `~/.config/base-cli/config.yaml` and applying defaults for flags not explicitly passed (via `applyConfigDefaults`)
- Interactive prompts (module path, description) when running in a TTY
- Resolving defaults: git username from `git config`, module path from conventions
- Deriving computed fields: `EnvPrefix` (uppercase, hyphens → underscores), `LibPackage` (hyphens stripped), `HasCLI`/`HasLib` from layout
- Building a `scaffold.Config` and calling `scaffold.Generate()`
- Running `go mod tidy` after generation (best-effort)
- Optionally running `git init` + initial commit

### `cmd/base-cli/config.go` — Config Management

Subcommands for managing `~/.config/base-cli/config.yaml`:
- `config init` — Create a commented config file with all fields
- `config show` — Print resolved values with source labels (config vs default)
- `config set <key> <value>` — Set a single value with validation
- `config edit` — Open in `$EDITOR`
- `config path` — Print the config file path

### `pkg/scaffold/scaffold.go` — Engine

The core of the project. A single public function:

```go
func Generate(cfg Config, destDir string) error
```

**Pipeline:**

1. **Walk** — `fs.WalkDir` over the embedded `templates/` filesystem
2. **Skip** — For each file, evaluate whether it should be skipped based on config (layout, CI provider, license, goreleaser, community, changelog, config package)
3. **Map path** — Transform the template path to its output path:
   - Strip `.tmpl` extension
   - Replace `{{BinaryName}}` and `{{LibPackage}}` placeholders
   - Apply directory prefix mappings (`github/` → `.github/`, etc.)
   - Apply special filename mappings (`gitignore` → `.gitignore`, etc.)
4. **Render** — Parse the template with `[% %]` delimiters and custom functions (`upper`, `lower`), execute with the Config as data
5. **Write** — Create parent directories and write the rendered content

Errors are wrapped with context at each step (`fmt.Errorf("rendering %s: %w", path, err)`).

### `pkg/scaffold/templates/` — Template Filesystem

All templates are embedded at compile time:

```go
//go:embed templates/*
var templateFS embed.FS
```

This means:
- **Zero runtime dependencies** — no files to ship or find at runtime
- **Atomic** — templates are versioned with the binary
- **Adding a template** just requires creating a `.tmpl` file in the right directory; it's automatically included in the next build

### `internal/config/config.go` — User Config

Handles loading and saving user-level defaults from `~/.config/base-cli/config.yaml`:

```go
func Load(path string) (*Config, error)   // Returns empty Config{} if file missing
func Save(cfg *Config, path string) error  // Creates parent dirs as needed
func DefaultPath() string                  // ~/.config/base-cli/config.yaml
```

Design decisions:
- **Silent on missing file** — `Load()` returns empty config if the file doesn't exist. The CLI works without any config file.
- **`*bool` for toggle fields** — `NoGitInit *bool` distinguishes "not set" (nil) from "explicitly false". Only non-nil values override CLI defaults.
- **CLI flags always win** — `applyConfigDefaults()` in init.go checks `cmd.Flags().Changed(name)` before applying config values.

### `internal/version/version.go` — Version Info

Three variables set via ldflags at build time:

```go
var (
    Version = "dev"
    Commit  = "none"
    Date    = "unknown"
)
```

The Makefile passes these during `go build`:

```makefile
-ldflags "-X internal/version.Version=$(VERSION) ..."
```

## File Generation Pipeline (Detail)

```
template path: github/workflows/ci.yml.tmpl
       │
       ├─ Skip check: CIGitHub == false? → skip entire file
       │
       ├─ Strip .tmpl: github/workflows/ci.yml
       │
       ├─ Dir mapping: github/ → .github/
       │   Result: .github/workflows/ci.yml
       │
       ├─ Parse template with [% %] delimiters
       │
       ├─ Execute template with Config data
       │
       └─ Write to: <destDir>/.github/workflows/ci.yml
```

## Separation of Concerns

```
cmd/base-cli/     → User-facing CLI (flags, prompts, config mgmt, UX)
pkg/scaffold/     → Public library (reusable by other tools)
internal/config/  → User-level config loading (~/.config/base-cli/)
internal/version/ → Build metadata (not exported)
```

The `pkg/scaffold/` package has no dependency on Cobra or any CLI framework. It accepts a plain struct and writes files. This means other Go programs can import it to generate projects programmatically without pulling in CLI dependencies.

## Testing

Tests use map-based table tests (`map[string]struct{}`):

```go
tests := map[string]struct {
    cfg      scaffold.Config
    wantFile string
    wantErr  bool
}{
    "both layout generates cmd and pkg": { ... },
    "lib layout skips cmd":              { ... },
}
```

The scaffold engine is tested by generating into a temp directory and asserting on the output file tree and contents.
