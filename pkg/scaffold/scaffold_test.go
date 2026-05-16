package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bothConfig(name string) Config {
	return Config{
		ProjectName:   name,
		ModulePath:    "github.com/test/" + name,
		BinaryName:    name,
		Description:   "A test project",
		Author:        "Test Author",
		Year:          "2026",
		GitUser:       "test",
		RepoURL:       "https://github.com/test/" + name,
		CIGitHub:      true,
		CIGitLab:      false,
		EnvPrefix:     "TEST_PROJECT",
		License:       "mit",
		Layout:        "both",
		HasCLI:        true,
		HasLib:        true,
		LibPackage:    strings.ReplaceAll(name, "-", ""),
		Goreleaser:    true,
		Community:     true,
		Changelog:     true,
		Config:        true,
		AgentMDClaude: true,
		AgentMDAgents: true,
		Todo:          true,
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
		"cmd/test-project/command_test.go",
		"cmd/test-project/ui.go",
		"cmd/test-project/help.go",
		"cmd/test-project/config.go",
		"internal/version/version.go",
		"internal/version/version_test.go",
		"internal/config/config.go",
		"internal/config/config_test.go",
		"pkg/testproject/doc.go",
		"Makefile",
		".goreleaser.yaml",
		"go.mod",
		"README.md",
		"CLAUDE.md",
		"AGENTS.md",
		"LICENSE",
		".gitignore",
		"install.sh",
		"uninstall.sh",
		"scripts/release.sh",
		".github/workflows/ci.yml",
		".github/workflows/release.yml",
		".github/ISSUE_TEMPLATE/bug_report.md",
		".github/ISSUE_TEMPLATE/feature_request.md",
		".github/pull_request_template.md",
		"CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md",
		"SECURITY.md",
		"CHANGELOG.yaml",
		"CHANGELOG.md",
		".chlog.yaml",
		"TODO.md",
		"assets/.gitkeep",
		"pkg/testproject/testdata/sample.yaml",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(destDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", f)
		}
	}

	rootContent, err := os.ReadFile(filepath.Join(destDir, "cmd/test-project/root.go"))
	if err != nil {
		t.Fatalf("reading generated root.go: %v", err)
	}
	for _, want := range []string{`StringVar(&configPathOverride, "config"`, `TEST_PROJECT_CONFIG`, `selectedConfigPath`} {
		if !strings.Contains(string(rootContent), want) {
			t.Errorf("generated root.go missing %q", want)
		}
	}

	configContent, err := os.ReadFile(filepath.Join(destDir, "cmd/test-project/config.go"))
	if err != nil {
		t.Fatalf("reading generated config.go: %v", err)
	}
	for _, want := range []string{`configSetCmd`, `configGetCmd`, `configToggleCmd`, `configKeysCmd`, `setYAMLPath`} {
		if !strings.Contains(string(configContent), want) {
			t.Errorf("generated config.go missing %q", want)
		}
	}

	makefile, err := os.ReadFile(filepath.Join(destDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading generated Makefile: %v", err)
	}
	makefileContent := string(makefile)
	for _, want := range []string{
		"bin: build ## Alias for build",
		"install-global: go-install ## Alias for go-install",
		"go-install: ## Install test-project to GOPATH/bin",
	} {
		if !strings.Contains(makefileContent, want) {
			t.Errorf("generated Makefile missing %q", want)
		}
	}

	securityContent, err := os.ReadFile(filepath.Join(destDir, "SECURITY.md"))
	if err != nil {
		t.Fatalf("reading generated SECURITY.md: %v", err)
	}
	for _, want := range []string{
		"## Supported Versions",
		"## Reporting a Vulnerability",
		"Please report suspected vulnerabilities privately",
		"## Response Expectations",
		"## Public Disclosure",
		"https://github.com/test/test-project/issues",
	} {
		if !strings.Contains(string(securityContent), want) {
			t.Errorf("generated SECURITY.md missing %q", want)
		}
	}

	versionContent, err := os.ReadFile(filepath.Join(destDir, "cmd/test-project/version.go"))
	if err != nil {
		t.Fatalf("reading generated version.go: %v", err)
	}
	if !strings.Contains(string(versionContent), `Aliases: []string{"v"}`) {
		t.Error("generated version.go should include v alias")
	}

	commandTestContent, err := os.ReadFile(filepath.Join(destDir, "cmd/test-project/command_test.go"))
	if err != nil {
		t.Fatalf("reading generated command_test.go: %v", err)
	}
	for _, want := range []string{
		`TestVersionCommandSmoke`,
		`TestVersionAliasSmoke`,
		`TestHelpCommandSmoke`,
		`TestCompletionCommandSmoke`,
		`TestConfigPathCommandSmoke`,
		`TestConfigPathFlagOverride`,
	} {
		if !strings.Contains(string(commandTestContent), want) {
			t.Errorf("generated command_test.go missing %q", want)
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
		"cmd/my-cli/help.go",
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
		"pkg/mylib/testdata/sample.yaml",
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
	for _, notWant := range []string{"bin:", "go-install:", "install-global:"} {
		if strings.Contains(string(makefile), notWant) {
			t.Errorf("lib-only Makefile should not contain %q", notWant)
		}
	}
}

func TestGenerateNoGoreleaser(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Goreleaser = false

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldNotExist := []string{
		".goreleaser.yaml",
		".github/workflows/release.yml",
		"scripts/release.sh",
	}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist when Goreleaser=false", f)
		}
	}

	// CI workflow should still exist
	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); os.IsNotExist(err) {
		t.Error("ci.yml should still exist when Goreleaser=false")
	}
}

func TestGenerateNoCommunity(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Community = false

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldNotExist := []string{
		".github/ISSUE_TEMPLATE/bug_report.md",
		".github/ISSUE_TEMPLATE/feature_request.md",
		".github/pull_request_template.md",
		"CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md",
		"SECURITY.md",
	}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist when Community=false", f)
		}
	}

	// CI workflow should still exist
	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); os.IsNotExist(err) {
		t.Error("ci.yml should still exist when Community=false")
	}
}

func TestGenerateNoChangelog(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Changelog = false

	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldNotExist := []string{
		"CHANGELOG.yaml",
		"CHANGELOG.md",
		".chlog.yaml",
	}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist when Changelog=false", f)
		}
	}

	// CI workflow should still exist
	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); os.IsNotExist(err) {
		t.Error("ci.yml should still exist when Changelog=false")
	}
}

func TestGeneratedReleaseScriptPreflight(t *testing.T) {
	cfg := bothConfig("my-cli")
	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "scripts/release.sh"))
	if err != nil {
		t.Fatalf("reading release script: %v", err)
	}
	releaseScript := string(content)

	for _, want := range []string{
		`SEMVER_REGEX=`,
		`[[ ! "$VERSION" =~ $SEMVER_REGEX ]]`,
		`git status --porcelain`,
		`git rev-parse -q --verify "refs/tags/${TAG}"`,
		`git ls-remote --exit-code --tags origin`,
		`make test`,
		`make lint`,
		`make build`,
		`goreleaser check`,
		`goreleaser release --snapshot --clean --skip=publish`,
		`require_cmd chlog`,
		`chlog validate`,
		`chlog check`,
		`chlog extract "${BARE_VERSION}" > .release/notes.md`,
	} {
		if !strings.Contains(releaseScript, want) {
			t.Errorf("release script missing %q", want)
		}
	}
}

func TestGeneratedReleaseScriptWithoutChangelogDoesNotRequireChlog(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Changelog = false
	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "scripts/release.sh"))
	if err != nil {
		t.Fatalf("reading release script: %v", err)
	}
	if strings.Contains(string(content), "chlog") {
		t.Error("release script should not reference chlog when changelog is disabled")
	}
}

func TestGeneratedMakefileReleaseTargetsUsePreflightScript(t *testing.T) {
	cfg := bothConfig("my-cli")
	destDir := t.TempDir()

	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	makefile := string(content)

	for _, want := range []string{
		`BUILD_VERSION?=$(shell git tag --sort=-v:refname 2>/dev/null | head -1)`,
		`-X ${MODULE_PATH}/internal/version.Version=${BUILD_VERSION}`,
		`ifndef VERSION`,
		`release: prep-release`,
		`@./scripts/release.sh $(VERSION)`,
	} {
		if !strings.Contains(makefile, want) {
			t.Errorf("Makefile missing %q", want)
		}
	}
	if strings.Contains(makefile, "git tag -a $(VERSION)") {
		t.Error("Makefile release target should delegate tagging to scripts/release.sh")
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
		"release yml":    {"github/workflows/release.yml.tmpl", "foo", "foo", "mit", ".github/workflows/release.yml"},
		"assets gitkeep": {"assets/gitkeep.tmpl", "foo", "foo", "mit", "assets/.gitkeep"},
		"testdata":       {"pkg/{{LibPackage}}/testdata/sample.yaml.tmpl", "foo", "foo", "mit", "pkg/foo/testdata/sample.yaml"},
		"changelog md":   {"CHANGELOG.md.tmpl", "foo", "foo", "mit", "CHANGELOG.md"},
		"todo md":        {"TODO.md.tmpl", "foo", "foo", "mit", "TODO.md"},
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

func TestGenerateBothCI(t *testing.T) {
	cfg := bothConfig("dual-ci")
	cfg.CIGitHub = true
	cfg.CIGitLab = true

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, ".github/workflows/ci.yml")); os.IsNotExist(err) {
		t.Error("GitHub CI should exist when CIGitHub=true")
	}
	if _, err := os.Stat(filepath.Join(destDir, ".gitlab-ci.yml")); os.IsNotExist(err) {
		t.Error("GitLab CI should exist when CIGitLab=true")
	}
}

func TestGenerateNoCI(t *testing.T) {
	cfg := bothConfig("no-ci")
	cfg.CIGitHub = false
	cfg.CIGitLab = false

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, ".github")); err == nil {
		t.Error(".github/ should not exist when CIGitHub=false")
	}
	if _, err := os.Stat(filepath.Join(destDir, ".gitlab-ci.yml")); err == nil {
		t.Error(".gitlab-ci.yml should not exist when CIGitLab=false")
	}
}

func TestGenerateReadOnlyDestDir(t *testing.T) {
	cfg := bothConfig("fail")
	destDir := filepath.Join(t.TempDir(), "readonly")
	if err := os.MkdirAll(destDir, 0o555); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(destDir, 0o755) }() // cleanup

	err := Generate(cfg, destDir)
	if err == nil {
		t.Error("expected error when writing to read-only directory")
	}
}

func TestSkipDir(t *testing.T) {
	tests := map[string]struct {
		relPath string
		cfg     Config
		want    error
	}{
		"github skipped": {
			relPath: "github/workflows",
			cfg:     Config{CIGitHub: false},
			want:    fs.SkipDir,
		},
		"github kept": {
			relPath: "github/workflows",
			cfg:     Config{CIGitHub: true},
			want:    nil,
		},
		"gitlab skipped": {
			relPath: "gitlab",
			cfg:     Config{CIGitLab: false},
			want:    fs.SkipDir,
		},
		"gitlab kept": {
			relPath: "gitlab",
			cfg:     Config{CIGitLab: true},
			want:    nil,
		},
		"cmd skipped for lib": {
			relPath: "cmd/myapp",
			cfg:     Config{HasCLI: false, HasLib: true},
			want:    fs.SkipDir,
		},
		"internal skipped for lib": {
			relPath: "internal/version",
			cfg:     Config{HasCLI: false, HasLib: true},
			want:    fs.SkipDir,
		},
		"scripts skipped for lib": {
			relPath: "scripts",
			cfg:     Config{HasCLI: false, HasLib: true},
			want:    fs.SkipDir,
		},
		"pkg skipped for cli": {
			relPath: "pkg/mylib",
			cfg:     Config{HasCLI: true, HasLib: false},
			want:    fs.SkipDir,
		},
		"pkg kept for lib": {
			relPath: "pkg/mylib",
			cfg:     Config{HasCLI: false, HasLib: true},
			want:    nil,
		},
		"skills skipped when no claude": {
			relPath: "skills/default",
			cfg:     Config{AgentMDClaude: false},
			want:    fs.SkipDir,
		},
		"skills kept when claude": {
			relPath: "skills/default",
			cfg:     Config{AgentMDClaude: true},
			want:    nil,
		},
		"internal/config skipped when Config=false": {
			relPath: "internal/config",
			cfg:     Config{HasCLI: true, Config: false},
			want:    fs.SkipDir,
		},
		"internal/config kept when Config=true": {
			relPath: "internal/config",
			cfg:     Config{HasCLI: true, Config: true},
			want:    nil,
		},
		"unrelated dir kept": {
			relPath: "docs",
			cfg:     Config{},
			want:    nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := skipDir(tt.relPath, tt.cfg)
			if got != tt.want {
				t.Errorf("skipDir(%q) = %v, want %v", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestSkipFile(t *testing.T) {
	tests := map[string]struct {
		relPath string
		cfg     Config
		want    bool
	}{
		"mit license kept": {
			relPath: "LICENSE_mit.tmpl",
			cfg:     Config{License: "mit", HasCLI: true},
			want:    false,
		},
		"mit license skipped for apache": {
			relPath: "LICENSE_mit.tmpl",
			cfg:     Config{License: "apache2", HasCLI: true},
			want:    true,
		},
		"apache license skipped for mit": {
			relPath: "LICENSE_apache2.tmpl",
			cfg:     Config{License: "mit", HasCLI: true},
			want:    true,
		},
		"apache license kept": {
			relPath: "LICENSE_apache2.tmpl",
			cfg:     Config{License: "apache2", HasCLI: true},
			want:    false,
		},
		"install.sh skipped for lib": {
			relPath: "install.sh.tmpl",
			cfg:     Config{HasCLI: false},
			want:    true,
		},
		"goreleaser skipped for lib": {
			relPath: "goreleaser.yaml.tmpl",
			cfg:     Config{HasCLI: false},
			want:    true,
		},
		"uninstall.sh skipped for lib": {
			relPath: "uninstall.sh.tmpl",
			cfg:     Config{HasCLI: false},
			want:    true,
		},
		"changelog yaml skipped": {
			relPath: "chlog.yaml.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Changelog: false},
			want:    true,
		},
		"changelog config skipped": {
			relPath: "chlog-config.yaml.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Changelog: false},
			want:    true,
		},
		"changelog md skipped": {
			relPath: "CHANGELOG.md.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Changelog: false},
			want:    true,
		},
		"changelog yaml kept": {
			relPath: "chlog.yaml.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Changelog: true},
			want:    false,
		},
		"CLAUDE.md skipped when no claude": {
			relPath: "CLAUDE.md.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, AgentMDClaude: false},
			want:    true,
		},
		"CLAUDE.md kept when claude": {
			relPath: "CLAUDE.md.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, AgentMDClaude: true},
			want:    false,
		},
		"AGENTS.md skipped when no agents": {
			relPath: "AGENTS.md.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, AgentMDAgents: false},
			want:    true,
		},
		"AGENTS.md kept when agents": {
			relPath: "AGENTS.md.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, AgentMDAgents: true},
			want:    false,
		},
		"config.go.tmpl skipped when Config=false": {
			relPath: "cmd/{{BinaryName}}/config.go.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Config: false},
			want:    true,
		},
		"config.go.tmpl kept when Config=true": {
			relPath: "cmd/{{BinaryName}}/config.go.tmpl",
			cfg:     Config{License: "mit", HasCLI: true, Config: true},
			want:    false,
		},
		"TODO.md skipped when Todo=false": {
			relPath: "TODO.md.tmpl",
			cfg:     Config{License: "mit", Todo: false},
			want:    true,
		},
		"TODO.md kept when Todo=true": {
			relPath: "TODO.md.tmpl",
			cfg:     Config{License: "mit", Todo: true},
			want:    false,
		},
		"SECURITY.md skipped when no community": {
			relPath: "SECURITY.md.tmpl",
			cfg:     Config{License: "mit", Community: false},
			want:    true,
		},
		"SECURITY.md kept when community": {
			relPath: "SECURITY.md.tmpl",
			cfg:     Config{License: "mit", Community: true},
			want:    false,
		},
		"regular file kept": {
			relPath: "Makefile.tmpl",
			cfg:     Config{License: "mit", HasCLI: true},
			want:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := skipFile(tt.relPath, tt.cfg)
			if got != tt.want {
				t.Errorf("skipFile(%q) = %v, want %v", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestMatchesPrefix(t *testing.T) {
	tests := map[string]struct {
		relPath string
		prefix  string
		want    bool
	}{
		"exact match":       {"github", "github", true},
		"prefix with slash": {"github/workflows", "github", true},
		"no match":          {"gitlab", "github", false},
		"partial no match":  {"githubx", "github", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := matchesPrefix(tt.relPath, tt.prefix)
			if got != tt.want {
				t.Errorf("matchesPrefix(%q, %q) = %v, want %v", tt.relPath, tt.prefix, got, tt.want)
			}
		})
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

func TestGenerateAgentMDClaudeOnly(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.AgentMDClaude = true
	cfg.AgentMDAgents = false

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "CLAUDE.md")); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist when AgentMDClaude=true")
	}
	if _, err := os.Stat(filepath.Join(destDir, ".skills/default/SKILL.md")); os.IsNotExist(err) {
		t.Error(".skills/ should exist when AgentMDClaude=true")
	}
	if _, err := os.Stat(filepath.Join(destDir, "AGENTS.md")); err == nil {
		t.Error("AGENTS.md should not exist when AgentMDAgents=false")
	}
}

func TestGenerateAgentMDAgentsOnly(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.AgentMDClaude = false
	cfg.AgentMDAgents = true

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "AGENTS.md")); os.IsNotExist(err) {
		t.Error("AGENTS.md should exist when AgentMDAgents=true")
	}
	if _, err := os.Stat(filepath.Join(destDir, "CLAUDE.md")); err == nil {
		t.Error("CLAUDE.md should not exist when AgentMDClaude=false")
	}
	if _, err := os.Stat(filepath.Join(destDir, ".skills")); err == nil {
		t.Error(".skills/ should not exist when AgentMDClaude=false")
	}
}

func TestGenerateAgentsMDContentForCLIAndLibrary(t *testing.T) {
	cfg := bothConfig("my-cli")

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("reading AGENTS.md: %v", err)
	}
	agents := string(content)

	for _, want := range []string{
		"- **Layout:** CLI + library",
		"make build          # Build ./bin/my-cli with version ldflags",
		"make install-global # Alias for go-install",
		"go run ./cmd/my-cli config keys",
		"Path priority: root `--config`, then `$TEST_PROJECT_CONFIG`, then the default path.",
		"CHANGELOG.yaml        # changelog source",
		"Release flow is `make prep-release VERSION=vX.Y.Z`, which runs `scripts/release.sh`.",
		"Stage explicit files only, not `git add .` or `git add -A`.",
		"Keep fixtures in `pkg/mycli/testdata/`",
	} {
		if !strings.Contains(agents, want) {
			t.Errorf("AGENTS.md missing %q", want)
		}
	}
}

func TestGenerateAgentsMDContentForDisabledOptions(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Config = false
	cfg.Changelog = false
	cfg.Goreleaser = false

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("reading AGENTS.md: %v", err)
	}
	agents := string(content)

	for _, want := range []string{
		"This scaffold was generated without config support.",
		"Release automation was not generated.",
		"Changelog support was disabled",
	} {
		if !strings.Contains(agents, want) {
			t.Errorf("AGENTS.md missing %q", want)
		}
	}
	for _, notWant := range []string{
		"go run ./cmd/my-cli config keys",
		"CHANGELOG.yaml        # changelog source",
		"Run `chlog check` after changelog edits.",
		"make prep-release VERSION=v0.1.0",
	} {
		if strings.Contains(agents, notWant) {
			t.Errorf("AGENTS.md should not contain %q", notWant)
		}
	}
}

func TestGenerateAgentsMDContentForLibraryOnly(t *testing.T) {
	cfg := bothConfig("my-lib")
	cfg.Layout = "lib"
	cfg.HasCLI = false
	cfg.HasLib = true
	cfg.Config = false
	cfg.Goreleaser = false
	cfg.BinaryName = "my-lib"
	cfg.LibPackage = "mylib"

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("reading AGENTS.md: %v", err)
	}
	agents := string(content)

	for _, want := range []string{
		"- **Layout:** Library only",
		"Library layout has no generated CLI config package or config commands.",
		"Library layout has no generated binary release flow.",
		"pkg/mylib/      # public library package",
	} {
		if !strings.Contains(agents, want) {
			t.Errorf("AGENTS.md missing %q", want)
		}
	}
	for _, notWant := range []string{
		"make build",
		"go run ./cmd/my-lib",
		"cmd/my-lib/",
		"GitHub Actions CI + release",
	} {
		if strings.Contains(agents, notWant) {
			t.Errorf("AGENTS.md should not contain %q", notWant)
		}
	}
}

func TestGenerateNoConfig(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.Config = false

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldNotExist := []string{
		"internal/config/config.go",
		"internal/config/config_test.go",
		"cmd/my-cli/config.go",
	}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist when Config=false", f)
		}
	}

	// Other internal files should still exist
	if _, err := os.Stat(filepath.Join(destDir, "internal/version/version.go")); os.IsNotExist(err) {
		t.Error("internal/version/version.go should still exist when Config=false")
	}

	// go.mod should not contain yaml.v3
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if strings.Contains(string(gomod), "yaml.v3") {
		t.Error("go.mod should not contain yaml.v3 when Config=false")
	}
}

func TestGenerateConfigContent(t *testing.T) {
	cfg := bothConfig("my-tool")
	cfg.ModulePath = "github.com/alice/my-tool"
	cfg.Config = true

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// config.go should reference the binary name and module path
	configGo, err := os.ReadFile(filepath.Join(destDir, "internal/config/config.go"))
	if err != nil {
		t.Fatalf("reading internal/config/config.go: %v", err)
	}
	if !strings.Contains(string(configGo), "my-tool") {
		t.Error("internal/config/config.go should reference binary name")
	}
	if !strings.Contains(string(configGo), "gopkg.in/yaml.v3") {
		t.Error("internal/config/config.go should import yaml.v3")
	}

	// cmd config.go should reference the module path
	cmdConfigGo, err := os.ReadFile(filepath.Join(destDir, "cmd/my-tool/config.go"))
	if err != nil {
		t.Fatalf("reading cmd/my-tool/config.go: %v", err)
	}
	if !strings.Contains(string(cmdConfigGo), "github.com/alice/my-tool/internal/config") {
		t.Error("cmd/my-tool/config.go should import internal/config")
	}

	// go.mod should contain yaml.v3
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if !strings.Contains(string(gomod), "gopkg.in/yaml.v3") {
		t.Error("go.mod should contain yaml.v3 when Config=true")
	}

	// root.go should register configCmd
	rootGo, err := os.ReadFile(filepath.Join(destDir, "cmd/my-tool/root.go"))
	if err != nil {
		t.Fatalf("reading cmd/my-tool/root.go: %v", err)
	}
	if !strings.Contains(string(rootGo), "configCmd") {
		t.Error("root.go should register configCmd when Config=true")
	}
}

func TestGenerateConfigLibLayout(t *testing.T) {
	cfg := bothConfig("my-lib")
	cfg.Layout = "lib"
	cfg.HasCLI = false
	cfg.HasLib = true
	cfg.Config = true // should be ignored for lib layout

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Config files should not exist: lib layout skips all of internal/ and cmd/
	shouldNotExist := []string{"internal", "cmd"}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist for lib layout", f)
		}
	}
}

func TestGenerateAgentMDNone(t *testing.T) {
	cfg := bothConfig("my-cli")
	cfg.AgentMDClaude = false
	cfg.AgentMDAgents = false

	destDir := t.TempDir()
	if err := Generate(cfg, destDir); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldNotExist := []string{"CLAUDE.md", "AGENTS.md", ".skills"}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(destDir, f)); err == nil {
			t.Errorf("%s should not exist when both AgentMD flags are false", f)
		}
	}
}
