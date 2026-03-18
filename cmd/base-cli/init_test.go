package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// resetInitFlags resets all init command flags to defaults.
// Necessary because cobra persists flag values on the global rootCmd.
func resetInitFlags() {
	flagDescription = ""
	flagAuthor = ""
	flagLicense = "mit"
	flagCI = "both"
	flagLayout = "both"
	flagDir = ""
	flagNoGitInit = false
	flagNoGoreleaser = false
	flagNoCommunity = false
	flagNoChangelog = false
	flagAgentMD = "both"
}

func TestRunInitNonInteractive(t *testing.T) {
	resetInitFlags()
	destDir := t.TempDir()
	projectDir := filepath.Join(destDir, "my-proj")

	rootCmd.SetArgs([]string{
		"init", "my-proj", "github.com/test/my-proj",
		"--description", "Test project",
		"--author", "Tester",
		"--license", "mit",
		"--ci", "github",
		"--layout", "both",
		"--dir", projectDir,
		"--no-git-init",
	})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	expectedFiles := []string{
		"go.mod",
		"Makefile",
		"README.md",
		"cmd/my-proj/main.go",
		"pkg/myproj/doc.go",
	}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(projectDir, f)); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", f)
		}
	}

	gomod, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if !strings.Contains(string(gomod), "module github.com/test/my-proj") {
		t.Errorf("go.mod has wrong module path: %s", gomod)
	}
}

func TestRunInitInvalidLicense(t *testing.T) {
	resetInitFlags()
	rootCmd.SetArgs([]string{
		"init", "test", "github.com/test/test",
		"--description", "x",
		"--license", "bsd",
		"--dir", t.TempDir(),
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for invalid license, got nil")
	}
}

func TestRunInitInvalidLayout(t *testing.T) {
	resetInitFlags()
	rootCmd.SetArgs([]string{
		"init", "test", "github.com/test/test",
		"--description", "x",
		"--layout", "invalid",
		"--dir", t.TempDir(),
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for invalid layout, got nil")
	}
}

func TestRunInitInvalidCI(t *testing.T) {
	resetInitFlags()
	rootCmd.SetArgs([]string{
		"init", "test", "github.com/test/test",
		"--description", "x",
		"--ci", "jenkins",
		"--dir", t.TempDir(),
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for invalid CI provider, got nil")
	}
}

func TestRunInitExistingNonEmptyDir(t *testing.T) {
	resetInitFlags()
	destDir := t.TempDir()
	projectDir := filepath.Join(destDir, "exists")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{
		"init", "exists", "github.com/test/exists",
		"--description", "x",
		"--dir", projectDir,
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for non-empty directory, got nil")
	}
}

func TestRunInitLibLayout(t *testing.T) {
	resetInitFlags()
	destDir := t.TempDir()
	projectDir := filepath.Join(destDir, "my-lib")

	rootCmd.SetArgs([]string{
		"init", "my-lib", "github.com/test/my-lib",
		"--description", "A library",
		"--license", "mit",
		"--ci", "github",
		"--layout", "lib",
		"--dir", projectDir,
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init (lib) failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectDir, "pkg/mylib/doc.go")); os.IsNotExist(err) {
		t.Error("pkg/mylib/doc.go should exist for lib layout")
	}

	if _, err := os.Stat(filepath.Join(projectDir, "cmd")); err == nil {
		t.Error("cmd/ should not exist for lib layout")
	}
}

func TestRunInitDerivedModuleNonInteractive(t *testing.T) {
	resetInitFlags()

	// Use a pipe (not a char device) so isTerminal() returns false.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	dir := t.TempDir()
	rootCmd.SetArgs([]string{
		"init", "test",
		"--description", "x",
		"--dir", dir,
		"--no-git-init",
	})

	// Module should be auto-derived from host + git user + project name.
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("expected module to be auto-derived in non-interactive mode, got error: %v", err)
	}
}

func TestRunInitInvalidAgentMD(t *testing.T) {
	resetInitFlags()
	rootCmd.SetArgs([]string{
		"init", "test", "github.com/test/test",
		"--description", "x",
		"--agent-md", "invalid",
		"--dir", t.TempDir(),
		"--no-git-init",
	})

	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for invalid agent-md, got nil")
	}
}
