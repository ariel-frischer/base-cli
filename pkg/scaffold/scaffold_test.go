package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bothConfig(name string) Config {
	return Config{
		ProjectName: name,
		ModulePath:  "github.com/test/" + name,
		BinaryName:  name,
		Description: "A test project",
		Author:      "Test Author",
		Year:        "2026",
		GitUser:     "test",
		RepoURL:     "https://github.com/test/" + name,
		CIGitHub:    true,
		CIGitLab:    false,
		EnvPrefix:   "TEST_PROJECT",
		License:     "mit",
		Layout:      "both",
		HasCLI:      true,
		HasLib:      true,
		LibPackage:  strings.ReplaceAll(name, "-", ""),
	}
}

func TestGenerate(t *testing.T) {
	cfg := bothConfig("test-project")
	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	expectedFiles := []string{
		"cmd/test-project/main.go",
		"cmd/test-project/root.go",
		"cmd/test-project/version.go",
		"cmd/test-project/ui.go",
		"internal/version/version.go",
		"internal/version/version_test.go",
		"pkg/testproject/doc.go",
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
	cfg := bothConfig("my-cli")
	cfg.ModulePath = "gitlab.com/test/my-cli"
	cfg.RepoURL = "https://gitlab.com/test/my-cli"
	cfg.CIGitHub = false
	cfg.CIGitLab = true
	cfg.EnvPrefix = "MY_CLI"
	cfg.License = "apache2"

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, ".gitlab-ci.yml")); os.IsNotExist(err) {
		t.Error(".gitlab-ci.yml should exist when CIGitLab=true")
	}

	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); err == nil {
		t.Error(".github/workflows/ci.yml should not exist when CIGitHub=false")
	}

	content, err := os.ReadFile(filepath.Join(destDir, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if len(content) == 0 {
		t.Error("LICENSE should not be empty")
	}
}

func TestGenerateNoLicense(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.License = "none"

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "LICENSE")); err == nil {
		t.Error("LICENSE should not exist when License=none")
	}
}

func TestGenerateLayoutCLI(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Layout = "cli"
	cfg.HasCLI = true
	cfg.HasLib = false

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// CLI files should exist
	cliFiles := []string{
		"cmd/my-cli/main.go",
		"cmd/my-cli/root.go",
		"internal/version/version.go",
		"install.sh",
		".goreleaser.yaml",
	}
	for _, f := range cliFiles {
		if _, err := os.Stat(filepath.Join(destDir, f)); os.IsNotExist(err) {
			t.Errorf("expected CLI file %s does not exist", f)
		}
	}

	// pkg/ should NOT exist
	if _, err := os.Stat(filepath.Join(destDir, "pkg")); err == nil {
		t.Error("pkg/ should not exist for cli layout")
	}
}

func TestGenerateLayoutLib(t *testing.T) {
	cfg := bothConfig("my-lib")
	cfg.Layout = "lib"
	cfg.HasCLI = false
	cfg.HasLib = true
	cfg.LibPackage = "mylib"

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Library files should exist
	libFiles := []string{
		"pkg/mylib/doc.go",
		"go.mod",
		"Makefile",
		"README.md",
		".gitignore",
	}
	for _, f := range libFiles {
		if _, err := os.Stat(filepath.Join(destDir, f)); os.IsNotExist(err) {
			t.Errorf("expected lib file %s does not exist", f)
		}
	}

	// CLI-specific files should NOT exist
	cliOnly := []string{"cmd", "internal", "install.sh", "uninstall.sh", ".goreleaser.yaml", "scripts"}
	for _, f := range cliOnly {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist for lib layout", f)
		}
	}

	// go.mod should NOT contain cobra
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if strings.Contains(string(gomod), "cobra") {
		t.Error("lib-only go.mod should not contain cobra")
	}

	// Makefile should NOT contain build target
	makefile, err := os.ReadFile(filepath.Join(destDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if strings.Contains(string(makefile), "build:") {
		t.Error("lib-only Makefile should not contain build target")
	}
}

func TestResolveOutputPath(t *testing.T) {
	tests := map[string]struct {
		relPath    string
		binaryName string
		libPackage string
		license    string
		want       string
	}{
		"strip tmpl":     {"Makefile.tmpl", "foo", "foo", "mit", "Makefile"},
		"binary name":    {"cmd/{{BinaryName}}/main.go.tmpl", "my-cli", "mycli", "mit", "cmd/my-cli/main.go"},
		"lib package":    {"pkg/{{LibPackage}}/doc.go.tmpl", "my-lib", "mylib", "mit", "pkg/mylib/doc.go"},
		"github prefix":  {"github/workflows/ci.yml.tmpl", "foo", "foo", "mit", ".github/workflows/ci.yml"},
		"gitlab prefix":  {"gitlab/gitlab-ci.yml.tmpl", "foo", "foo", "mit", ".gitlab-ci.yml"},
		"gitignore":      {"gitignore.tmpl", "foo", "foo", "mit", ".gitignore"},
		"license mit":    {"LICENSE_mit.tmpl", "foo", "foo", "mit", "LICENSE"},
		"license apache": {"LICENSE_apache2.tmpl", "foo", "foo", "apache2", "LICENSE"},
		"changelog yaml": {"chlog.yaml.tmpl", "foo", "foo", "mit", "CHANGELOG.yaml"},
		"chlog config":   {"chlog-config.yaml.tmpl", "foo", "foo", "mit", ".chlog.yaml"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := resolveOutputPath(tt.relPath, tt.binaryName, tt.libPackage, tt.license)
			if got != tt.want {
				t.Errorf("resolveOutputPath(%q, %q, %q, %q) = %q, want %q",
					tt.relPath, tt.binaryName, tt.libPackage, tt.license, got, tt.want)
			}
		})
	}
}

func TestGenerateShellScriptsExecutable(t *testing.T) {
	cfg := bothConfig("test-project")
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
	cfg := bothConfig("my-tool")
	cfg.ModulePath = "github.com/alice/my-tool"
	cfg.BinaryName = "my-tool"
	cfg.Description = "A cool tool"
	cfg.Author = "Alice"
	cfg.GitUser = "alice"
	cfg.RepoURL = "https://github.com/alice/my-tool"
	cfg.EnvPrefix = "MY_TOOL"
	cfg.LibPackage = "mytool"

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Check go.mod has correct module path
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if !strings.Contains(string(gomod), "module github.com/alice/my-tool") {
		t.Errorf("go.mod should contain module path, got:\n%s", gomod)
	}

	// Check root.go imports correct module
	rootGo, err := os.ReadFile(filepath.Join(destDir, "cmd/my-tool/root.go"))
	if err != nil {
		t.Fatalf("reading root.go: %v", err)
	}
	if !strings.Contains(string(rootGo), "github.com/alice/my-tool") {
		t.Errorf("root.go should reference module path")
	}

	// Check Makefile has correct module path
	makefile, err := os.ReadFile(filepath.Join(destDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if !strings.Contains(string(makefile), "MODULE_PATH=github.com/alice/my-tool") {
		t.Errorf("Makefile should contain MODULE_PATH")
	}

	// Check pkg doc.go exists with correct package name
	docGo, err := os.ReadFile(filepath.Join(destDir, "pkg/mytool/doc.go"))
	if err != nil {
		t.Fatalf("reading pkg/mytool/doc.go: %v", err)
	}
	if !strings.Contains(string(docGo), "package mytool") {
		t.Errorf("doc.go should declare package mytool, got:\n%s", docGo)
	}
}
