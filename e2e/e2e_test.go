//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// binaryPath holds the compiled binary for all tests in this suite.
var binaryPath string

func TestMain(m *testing.M) {
	bin, err := compileBinary()
	if err != nil {
		panic("failed to build base-cli: " + err.Error())
	}
	binaryPath = bin
	os.Exit(m.Run())
}

func compileBinary() (string, error) {
	tmp, err := os.MkdirTemp("", "base-cli-e2e-bin-*")
	if err != nil {
		return "", err
	}
	bin := filepath.Join(tmp, "base-cli")
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/base-cli/")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmp)
		return "", fmt.Errorf("build failed: %w\n%s", err, out)
	}
	return bin, nil
}

// run executes the binary with args and returns combined output + error.
func run(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// scaffold runs init with --no-git-init and required non-interactive flags,
// returning the output dir. Module path is passed as a positional arg.
func scaffold(t *testing.T, name string, extra ...string) string {
	t.Helper()
	dir := t.TempDir()
	args := []string{
		"init", name, "github.com/e2e/" + name,
		"--dir", dir,
		"--description", "e2e test project",
		"--no-git-init",
	}
	args = append(args, extra...)
	out, err := run(args...)
	if err != nil {
		t.Fatalf("scaffold %q failed:\n%s", name, out)
	}
	return dir
}

// builds runs `go build ./...` in the given directory.
func builds(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./... failed in %s:\n%s", dir, out)
	}
}

func exists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file/dir to exist: %s", path)
	}
}

func notExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected file/dir to NOT exist: %s", path)
	}
}

// --- version ---

func TestVersionPlain(t *testing.T) {
	out, err := run("version", "--plain")
	if err != nil {
		t.Fatalf("version --plain: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestVersionAlias(t *testing.T) {
	out, err := run("v")
	if err != nil {
		t.Fatalf("version alias: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

// --- completion ---

func TestCompletion(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		t.Run(shell, func(t *testing.T) {
			out, err := run("completion", shell)
			if err != nil {
				t.Fatalf("completion %s: %v\n%s", shell, err, out)
			}
			if len(out) == 0 {
				t.Errorf("completion %s: empty output", shell)
			}
		})
	}
}

// --- layouts ---

func TestLayoutBothBuilds(t *testing.T) {
	dir := scaffold(t, "proj-both")
	exists(t, filepath.Join(dir, "cmd/proj-both/main.go"))
	exists(t, filepath.Join(dir, "pkg/projboth/doc.go"))
	notExists(t, filepath.Join(dir, ".git"))
	builds(t, dir)
}

func TestLayoutCLIBuilds(t *testing.T) {
	dir := scaffold(t, "proj-cli", "--layout", "cli")
	exists(t, filepath.Join(dir, "cmd/proj-cli/main.go"))
	exists(t, filepath.Join(dir, "internal/version/version.go"))
	notExists(t, filepath.Join(dir, "pkg"))
	builds(t, dir)
}

func TestLayoutLibBuilds(t *testing.T) {
	dir := scaffold(t, "proj-lib", "--layout", "lib")
	exists(t, filepath.Join(dir, "pkg/projlib/doc.go"))
	notExists(t, filepath.Join(dir, "cmd"))
	notExists(t, filepath.Join(dir, "internal"))
	builds(t, dir)
}

// --- CI ---

func TestCIGitHubOnly(t *testing.T) {
	dir := scaffold(t, "proj-gh", "--ci", "github")
	exists(t, filepath.Join(dir, ".github/workflows/ci.yml"))
	notExists(t, filepath.Join(dir, ".gitlab-ci.yml"))
}

func TestCIGitLabOnly(t *testing.T) {
	dir := scaffold(t, "proj-gl", "--ci", "gitlab")
	exists(t, filepath.Join(dir, ".gitlab-ci.yml"))
	notExists(t, filepath.Join(dir, ".github/workflows/ci.yml"))
}

func TestCIBoth(t *testing.T) {
	dir := scaffold(t, "proj-ci-both", "--ci", "both")
	exists(t, filepath.Join(dir, ".github/workflows/ci.yml"))
	exists(t, filepath.Join(dir, ".gitlab-ci.yml"))
}

// --- optional feature flags ---

func TestNoGoreleaser(t *testing.T) {
	dir := scaffold(t, "proj-nogr", "--no-goreleaser")
	notExists(t, filepath.Join(dir, ".goreleaser.yaml"))
	notExists(t, filepath.Join(dir, ".github/workflows/release.yml"))
	notExists(t, filepath.Join(dir, "scripts/release.sh"))
	exists(t, filepath.Join(dir, ".github/workflows/ci.yml"))
}

func TestNoCommunity(t *testing.T) {
	dir := scaffold(t, "proj-nocom", "--no-community")
	notExists(t, filepath.Join(dir, ".github/ISSUE_TEMPLATE"))
	notExists(t, filepath.Join(dir, "CONTRIBUTING.md"))
	notExists(t, filepath.Join(dir, "CODE_OF_CONDUCT.md"))
}

func TestNoChangelog(t *testing.T) {
	dir := scaffold(t, "proj-nocl", "--no-changelog")
	notExists(t, filepath.Join(dir, "CHANGELOG.yaml"))
	notExists(t, filepath.Join(dir, "CHANGELOG.md"))
	notExists(t, filepath.Join(dir, ".chlog.yaml"))
}

func TestConfigDefault(t *testing.T) {
	dir := scaffold(t, "proj-cfg")
	exists(t, filepath.Join(dir, "internal/config/config.go"))
	exists(t, filepath.Join(dir, "internal/config/config_test.go"))
	exists(t, filepath.Join(dir, "cmd/proj-cfg/config.go"))
	builds(t, dir)
}

func TestNoConfig(t *testing.T) {
	dir := scaffold(t, "proj-nocfg", "--no-config")
	notExists(t, filepath.Join(dir, "internal/config"))
	notExists(t, filepath.Join(dir, "cmd/proj-nocfg/config.go"))
	// internal/version should still exist
	exists(t, filepath.Join(dir, "internal/version/version.go"))
	builds(t, dir)
}

func TestNoConfigLibLayout(t *testing.T) {
	// lib layout never generates config regardless of flag
	dir := scaffold(t, "proj-libcfg", "--layout", "lib")
	notExists(t, filepath.Join(dir, "internal"))
	notExists(t, filepath.Join(dir, "cmd"))
	builds(t, dir)
}

func TestNoGitInit(t *testing.T) {
	// scaffold() already passes --no-git-init; just verify
	dir := scaffold(t, "proj-nogit")
	notExists(t, filepath.Join(dir, ".git"))
}

// --- agent-md ---

func TestAgentMDBoth(t *testing.T) {
	dir := scaffold(t, "proj-agent-both", "--agent-md", "both")
	exists(t, filepath.Join(dir, "CLAUDE.md"))
	exists(t, filepath.Join(dir, "AGENTS.md"))
}

func TestAgentMDClaude(t *testing.T) {
	dir := scaffold(t, "proj-agent-claude", "--agent-md", "claude")
	exists(t, filepath.Join(dir, "CLAUDE.md"))
	notExists(t, filepath.Join(dir, "AGENTS.md"))
}

func TestAgentMDAgents(t *testing.T) {
	dir := scaffold(t, "proj-agent-agents", "--agent-md", "agents")
	notExists(t, filepath.Join(dir, "CLAUDE.md"))
	exists(t, filepath.Join(dir, "AGENTS.md"))
}

func TestAgentMDNone(t *testing.T) {
	dir := scaffold(t, "proj-agent-none", "--agent-md", "none")
	notExists(t, filepath.Join(dir, "CLAUDE.md"))
	notExists(t, filepath.Join(dir, "AGENTS.md"))
}

// --- licenses ---

func TestLicenseMIT(t *testing.T) {
	dir := scaffold(t, "proj-mit", "--license", "mit")
	exists(t, filepath.Join(dir, "LICENSE"))
}

func TestLicenseApache2(t *testing.T) {
	dir := scaffold(t, "proj-apache", "--license", "apache2", "--author", "Test Author")
	content, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if string(content) == "" {
		t.Error("LICENSE should not be empty")
	}
}

func TestLicenseNone(t *testing.T) {
	dir := scaffold(t, "proj-nolic", "--license", "none")
	notExists(t, filepath.Join(dir, "LICENSE"))
}

// --- error cases ---

func TestInitNoArgs(t *testing.T) {
	_, err := run("init")
	if err == nil {
		t.Error("expected error when no project name given")
	}
}

func TestInitExistingNonEmptyDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("x"), 0o644)
	_, err := run("init", "proj", "github.com/e2e/proj",
		"--dir", dir,
		"--description", "test",
		"--no-git-init",
	)
	if err == nil {
		t.Error("expected error for non-empty directory")
	}
}

func TestInitInvalidAgentMD(t *testing.T) {
	_, err := run("init", "proj", "github.com/e2e/proj",
		"--dir", t.TempDir(),
		"--description", "test",
		"--agent-md", "invalid",
		"--no-git-init",
	)
	if err == nil {
		t.Error("expected error for invalid --agent-md")
	}
}

func TestInitInvalidLicense(t *testing.T) {
	_, err := run("init", "proj", "github.com/e2e/proj",
		"--dir", t.TempDir(),
		"--description", "test",
		"--license", "bsd",
		"--no-git-init",
	)
	if err == nil {
		t.Error("expected error for invalid --license")
	}
}

func TestInitInvalidCI(t *testing.T) {
	_, err := run("init", "proj", "github.com/e2e/proj",
		"--dir", t.TempDir(),
		"--description", "test",
		"--ci", "jenkins",
		"--no-git-init",
	)
	if err == nil {
		t.Error("expected error for invalid --ci")
	}
}

// --- minimal scaffold ---

func TestMinimalScaffoldBuilds(t *testing.T) {
	dir := scaffold(t, "proj-minimal",
		"--no-goreleaser",
		"--no-community",
		"--no-changelog",
		"--no-config",
		"--agent-md", "none",
	)
	notExists(t, filepath.Join(dir, ".goreleaser.yaml"))
	notExists(t, filepath.Join(dir, "CHANGELOG.yaml"))
	notExists(t, filepath.Join(dir, "CLAUDE.md"))
	notExists(t, filepath.Join(dir, "internal/config"))
	builds(t, dir)
}
