package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	cfg := Config{
		ProjectName: "test-project",
		ModulePath:  "github.com/test/test-project",
		BinaryName:  "test-project",
		Description: "A test project",
		Author:      "Test Author",
		Year:        "2026",
		GitUser:     "test",
		RepoURL:     "https://github.com/test/test-project",
		CIGitHub:    true,
		CIGitLab:    false,
		EnvPrefix:   "TEST_PROJECT",
		License:     "mit",
	}

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Verify expected files exist
	expectedFiles := []string{
		"cmd/test-project/main.go",
		"cmd/test-project/root.go",
		"cmd/test-project/version.go",
		"cmd/test-project/ui.go",
		"internal/version/version.go",
		"internal/version/version_test.go",
		"Makefile",
		".goreleaser.yaml",
		"go.mod",
		"README.md",
		"CLAUDE.md",
		"LICENSE",
		".gitignore",
		"install.sh",
		"uninstall.sh",
		"scripts/release.sh",
		".github/workflows/ci.yml",
		"CHANGELOG.yaml",
		".chlog.yaml",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(destDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", f)
		}
	}

	// Verify GitLab CI was NOT generated (CIGitLab=false)
	gitlabPath := filepath.Join(destDir, ".gitlab-ci.yml")
	if _, err := os.Stat(gitlabPath); err == nil {
		t.Error(".gitlab-ci.yml should not exist when CIGitLab=false")
	}
}

func TestGenerateGitLab(t *testing.T) {
	cfg := Config{
		ProjectName: "my-cli",
		ModulePath:  "gitlab.com/test/my-cli",
		BinaryName:  "my-cli",
		Description: "Test",
		Author:      "Test",
		Year:        "2026",
		GitUser:     "test",
		RepoURL:     "https://gitlab.com/test/my-cli",
		CIGitHub:    false,
		CIGitLab:    true,
		EnvPrefix:   "MY_CLI",
		License:     "apache2",
	}

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// GitLab CI should exist
	if _, err := os.Stat(filepath.Join(destDir, ".gitlab-ci.yml")); os.IsNotExist(err) {
		t.Error(".gitlab-ci.yml should exist when CIGitLab=true")
	}

	// GitHub CI should NOT exist
	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); err == nil {
		t.Error(".github/workflows/ci.yml should not exist when CIGitHub=false")
	}

	// Apache license should exist
	content, err := os.ReadFile(filepath.Join(destDir, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if len(content) == 0 {
		t.Error("LICENSE should not be empty")
	}
}

func TestGenerateNoLicense(t *testing.T) {
	cfg := Config{
		ProjectName: "my-cli",
		ModulePath:  "github.com/test/my-cli",
		BinaryName:  "my-cli",
		Description: "Test",
		Author:      "Test",
		Year:        "2026",
		GitUser:     "test",
		RepoURL:     "https://github.com/test/my-cli",
		CIGitHub:    true,
		CIGitLab:    false,
		EnvPrefix:   "MY_CLI",
		License:     "none",
	}

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// No LICENSE file should exist
	if _, err := os.Stat(filepath.Join(destDir, "LICENSE")); err == nil {
		t.Error("LICENSE should not exist when License=none")
	}
}

func TestResolveOutputPath(t *testing.T) {
	tests := map[string]struct {
		relPath    string
		binaryName string
		license    string
		want       string
	}{
		"strip tmpl":         {"Makefile.tmpl", "foo", "mit", "Makefile"},
		"binary name":        {"cmd/{{BinaryName}}/main.go.tmpl", "my-cli", "mit", "cmd/my-cli/main.go"},
		"github prefix":      {"github/workflows/ci.yml.tmpl", "foo", "mit", ".github/workflows/ci.yml"},
		"gitlab prefix":      {"gitlab/gitlab-ci.yml.tmpl", "foo", "mit", ".gitlab-ci.yml"},
		"gitignore":          {"gitignore.tmpl", "foo", "mit", ".gitignore"},
		"license mit":        {"LICENSE_mit.tmpl", "foo", "mit", "LICENSE"},
		"license apache":     {"LICENSE_apache2.tmpl", "foo", "apache2", "LICENSE"},
		"changelog yaml":     {"chlog.yaml.tmpl", "foo", "mit", "CHANGELOG.yaml"},
		"chlog config":       {"chlog-config.yaml.tmpl", "foo", "mit", ".chlog.yaml"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := resolveOutputPath(tt.relPath, tt.binaryName, tt.license)
			if got != tt.want {
				t.Errorf("resolveOutputPath(%q, %q, %q) = %q, want %q",
					tt.relPath, tt.binaryName, tt.license, got, tt.want)
			}
		})
	}
}

func TestGenerateShellScriptsExecutable(t *testing.T) {
	cfg := Config{
		ProjectName: "test-project",
		ModulePath:  "github.com/test/test-project",
		BinaryName:  "test-project",
		Description: "Test",
		Author:      "Test",
		Year:        "2026",
		GitUser:     "test",
		RepoURL:     "https://github.com/test/test-project",
		CIGitHub:    true,
		CIGitLab:    false,
		EnvPrefix:   "TEST_PROJECT",
		License:     "mit",
	}

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	scripts := []string{"install.sh", "uninstall.sh", "scripts/release.sh"}
	for _, s := range scripts {
		info, err := os.Stat(filepath.Join(destDir, s))
		if err != nil {
			t.Errorf("stat %s: %v", s, err)
			continue
		}
		if info.Mode()&0o111 == 0 {
			t.Errorf("%s should be executable, got mode %v", s, info.Mode())
		}
	}
}

func TestGenerateFileContents(t *testing.T) {
	cfg := Config{
		ProjectName: "my-tool",
		ModulePath:  "github.com/alice/my-tool",
		BinaryName:  "my-tool",
		Description: "A cool tool",
		Author:      "Alice",
		Year:        "2026",
		GitUser:     "alice",
		RepoURL:     "https://github.com/alice/my-tool",
		CIGitHub:    true,
		CIGitLab:    false,
		EnvPrefix:   "MY_TOOL",
		License:     "mit",
	}

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Check go.mod has correct module path
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if got := string(gomod); !contains(got, "module github.com/alice/my-tool") {
		t.Errorf("go.mod should contain module path, got:\n%s", got)
	}

	// Check root.go imports correct module
	rootGo, err := os.ReadFile(filepath.Join(destDir, "cmd/my-tool/root.go"))
	if err != nil {
		t.Fatalf("reading root.go: %v", err)
	}
	if !contains(string(rootGo), "github.com/alice/my-tool") {
		t.Errorf("root.go should reference module path")
	}

	// Check Makefile has correct module path
	makefile, err := os.ReadFile(filepath.Join(destDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if !contains(string(makefile), "MODULE_PATH=github.com/alice/my-tool") {
		t.Errorf("Makefile should contain MODULE_PATH")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
