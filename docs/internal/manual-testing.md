# Manual Testing Guide

**Last verified:** 2026-03-17 @ commit `6ba62cf` — all tests passing

All commands run from repo root using `go run ./cmd/base-cli/`.
Test outputs go to `/tmp/base-cli-test/`.

## Setup

```bash
rm -rf /tmp/base-cli-test && mkdir -p /tmp/base-cli-test
```

**Key behaviors to know:**
- `--dir` is the direct output path, not a parent. Files land in `<dir>/`, not `<dir>/<name>/`.
- In non-interactive mode (no TTY), `--module` and `--description` are required. Interactive mode prompts for them with defaults.
- `--agent-md` validation runs after `--module`/`--description` checks.

## 1. version

```bash
go run ./cmd/base-cli/ version
go run ./cmd/base-cli/ version --plain
go run ./cmd/base-cli/ v
go run ./cmd/base-cli/ --version
go run ./cmd/base-cli/ version --no-color
```

## 2. config

```bash
# Show config path
go run ./cmd/base-cli/ config path

# Init default config
go run ./cmd/base-cli/ config init

# Show resolved config
go run ./cmd/base-cli/ config show

# Set values
go run ./cmd/base-cli/ config set author "Test User"
go run ./cmd/base-cli/ config set license apache2
go run ./cmd/base-cli/ config set ci github
go run ./cmd/base-cli/ config set layout cli
go run ./cmd/base-cli/ config set agent_md none
go run ./cmd/base-cli/ config set no_git_init true
go run ./cmd/base-cli/ config set no_goreleaser true
go run ./cmd/base-cli/ config set no_community true
go run ./cmd/base-cli/ config set no_changelog true

# Verify changes took
go run ./cmd/base-cli/ config show

# Reset back
go run ./cmd/base-cli/ config set license mit
go run ./cmd/base-cli/ config set ci both
go run ./cmd/base-cli/ config set layout both
go run ./cmd/base-cli/ config set agent_md both
go run ./cmd/base-cli/ config set no_git_init false
go run ./cmd/base-cli/ config set no_goreleaser false
go run ./cmd/base-cli/ config set no_community false
go run ./cmd/base-cli/ config set no_changelog false

# Invalid key
go run ./cmd/base-cli/ config set bogus_key value

# Edit (opens $EDITOR)
# go run ./cmd/base-cli/ config edit
```

## 3. init — layout variations

Note: in a shell (no TTY), pass `--module` and `--description` to skip prompts.

```bash
# Default (both layout) — interactive prompts for module + description
go run ./cmd/base-cli/ init myproject --dir /tmp/base-cli-test/default
# Files land directly in /tmp/base-cli-test/default/, not in a myproject/ subdir
ls /tmp/base-cli-test/default/cmd /tmp/base-cli-test/default/pkg

# CLI only
go run ./cmd/base-cli/ init mycli --dir /tmp/base-cli-test/cli --layout cli \
  --module github.com/me/mycli --description "my cli"
ls /tmp/base-cli-test/cli/cmd/ /tmp/base-cli-test/cli/internal/
ls /tmp/base-cli-test/cli/pkg/ 2>&1 | grep -q "No such file" && echo "PASS: no pkg/" || echo "FAIL"

# Lib only
go run ./cmd/base-cli/ init mylib --dir /tmp/base-cli-test/lib --layout lib \
  --module github.com/me/mylib --description "my lib"
ls /tmp/base-cli-test/lib/pkg/
ls /tmp/base-cli-test/lib/cmd/ 2>&1 | grep -q "No such file" && echo "PASS: no cmd/" || echo "FAIL"
```

## 4. init — CI variations

```bash
# GitHub only
go run ./cmd/base-cli/ init proj-gh --dir /tmp/base-cli-test/ci --ci github
ls /tmp/base-cli-test/ci/proj-gh/.github/workflows/
test ! -f /tmp/base-cli-test/ci/proj-gh/.gitlab-ci.yml && echo "PASS: no gitlab" || echo "FAIL"

# GitLab only
go run ./cmd/base-cli/ init proj-gl --dir /tmp/base-cli-test/ci --ci gitlab
test -f /tmp/base-cli-test/ci/proj-gl/.gitlab-ci.yml && echo "PASS: has gitlab" || echo "FAIL"
test ! -d /tmp/base-cli-test/ci/proj-gl/.github && echo "PASS: no github" || echo "FAIL"

# Both (default)
go run ./cmd/base-cli/ init proj-both --dir /tmp/base-cli-test/ci --ci both
ls /tmp/base-cli-test/ci/proj-both/.github/workflows/
test -f /tmp/base-cli-test/ci/proj-both/.gitlab-ci.yml && echo "PASS" || echo "FAIL"
```

## 5. init — optional features

```bash
# No goreleaser
go run ./cmd/base-cli/ init proj-nogr --dir /tmp/base-cli-test/opts --no-goreleaser
test ! -f /tmp/base-cli-test/opts/proj-nogr/.goreleaser.yaml && echo "PASS" || echo "FAIL"

# No community files
go run ./cmd/base-cli/ init proj-nocom --dir /tmp/base-cli-test/opts --no-community
test ! -d /tmp/base-cli-test/opts/proj-nocom/.github/ISSUE_TEMPLATE && echo "PASS" || echo "FAIL"

# No changelog
go run ./cmd/base-cli/ init proj-nocl --dir /tmp/base-cli-test/opts --no-changelog
test ! -f /tmp/base-cli-test/opts/proj-nocl/CHANGELOG.yaml && echo "PASS" || echo "FAIL"
test ! -f /tmp/base-cli-test/opts/proj-nocl/CHANGELOG.md && echo "PASS" || echo "FAIL"

# No git init
go run ./cmd/base-cli/ init proj-nogit --dir /tmp/base-cli-test/opts --no-git-init
test ! -d /tmp/base-cli-test/opts/proj-nogit/.git && echo "PASS" || echo "FAIL"

# All opts disabled
go run ./cmd/base-cli/ init proj-minimal --dir /tmp/base-cli-test/opts \
  --no-goreleaser --no-community --no-changelog --no-git-init --agent-md none
ls /tmp/base-cli-test/opts/proj-minimal/
```

## 6. init — agent-md variations

```bash
go run ./cmd/base-cli/ init proj-agent-both --dir /tmp/base-cli-test/agent --agent-md both
test -f /tmp/base-cli-test/agent/proj-agent-both/CLAUDE.md && echo "PASS" || echo "FAIL"
test -f /tmp/base-cli-test/agent/proj-agent-both/AGENTS.md && echo "PASS" || echo "FAIL"

go run ./cmd/base-cli/ init proj-agent-claude --dir /tmp/base-cli-test/agent --agent-md claude
test -f /tmp/base-cli-test/agent/proj-agent-claude/CLAUDE.md && echo "PASS" || echo "FAIL"
test ! -f /tmp/base-cli-test/agent/proj-agent-claude/AGENTS.md && echo "PASS" || echo "FAIL"

go run ./cmd/base-cli/ init proj-agent-agents --dir /tmp/base-cli-test/agent --agent-md agents
test ! -f /tmp/base-cli-test/agent/proj-agent-agents/CLAUDE.md && echo "PASS" || echo "FAIL"
test -f /tmp/base-cli-test/agent/proj-agent-agents/AGENTS.md && echo "PASS" || echo "FAIL"

go run ./cmd/base-cli/ init proj-agent-none --dir /tmp/base-cli-test/agent --agent-md none
test ! -f /tmp/base-cli-test/agent/proj-agent-none/CLAUDE.md && echo "PASS" || echo "FAIL"
test ! -f /tmp/base-cli-test/agent/proj-agent-none/AGENTS.md && echo "PASS" || echo "FAIL"

# Invalid value
go run ./cmd/base-cli/ init proj-agent-bad --dir /tmp/base-cli-test/agent --agent-md invalid
```

## 7. init — other flags

```bash
# Custom module, author, description, license
go run ./cmd/base-cli/ init proj-custom --dir /tmp/base-cli-test/custom \
  --module github.com/testuser/proj-custom \
  --author "Test Author" \
  --description "A test project" \
  --license apache2
grep -q "Test Author" /tmp/base-cli-test/custom/proj-custom/LICENSE && echo "PASS" || echo "FAIL"
grep -q "apache" /tmp/base-cli-test/custom/proj-custom/LICENSE && echo "PASS: apache" || echo "FAIL"

# No license
go run ./cmd/base-cli/ init proj-nolic --dir /tmp/base-cli-test/custom --license none
test ! -f /tmp/base-cli-test/custom/proj-nolic/LICENSE && echo "PASS" || echo "FAIL"
```

## 8. init — generated project builds

```bash
cd /tmp/base-cli-test/default/myproject && go build ./... && echo "PASS: builds" || echo "FAIL"
cd /tmp/base-cli-test/cli/mycli && go build ./... && echo "PASS: builds" || echo "FAIL"
cd /tmp/base-cli-test/lib/mylib && go build ./... && echo "PASS: builds" || echo "FAIL"
```

## 9. init — edge cases

```bash
# Dir already exists and is not empty (should error)
go run ./cmd/base-cli/ init myproject --dir /tmp/base-cli-test/default \
  --module github.com/me/myproject --description "test"
# Expected: "directory ... already exists and is not empty"

# No project name (should error)
go run ./cmd/base-cli/ init
# Expected: "accepts 1 arg(s), received 0"

# Non-interactive missing --module (should error)
go run ./cmd/base-cli/ init myproject --dir /tmp/base-cli-test/edge --description "test"
# Expected: "--module is required in non-interactive mode"

# --no-color flag
go run ./cmd/base-cli/ init proj-nocolor --dir /tmp/base-cli-test/edge --no-color \
  --module github.com/me/proj-nocolor --description "test"
```

## 10. completion

```bash
go run ./cmd/base-cli/ completion bash > /dev/null && echo "PASS: bash" || echo "FAIL"
go run ./cmd/base-cli/ completion zsh > /dev/null && echo "PASS: zsh" || echo "FAIL"
go run ./cmd/base-cli/ completion fish > /dev/null && echo "PASS: fish" || echo "FAIL"
go run ./cmd/base-cli/ completion powershell > /dev/null && echo "PASS: powershell" || echo "FAIL"
```

## 11. uninstall (CAUTION)

```bash
# Dry check — just see the help, don't actually run without --yes
go run ./cmd/base-cli/ uninstall --help
# To actually test: go run ./cmd/base-cli/ uninstall (will prompt)
```

## Cleanup

```bash
rm -rf /tmp/base-cli-test
```
