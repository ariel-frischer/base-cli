// Package scaffold generates Go project scaffolds from embedded templates.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// Config holds the template variables for scaffold generation.
type Config struct {
	ProjectName   string
	ModulePath    string
	BinaryName    string
	Description   string
	Author        string
	Year          string
	GitUser       string
	RepoURL       string
	CIGitHub      bool
	CIGitLab      bool
	EnvPrefix     string
	License       string // "mit", "apache2", "none"
	Layout        string // "both", "cli", "lib"
	HasCLI        bool   // true for "both" and "cli"
	HasLib        bool   // true for "both" and "lib"
	LibPackage    string // Go-safe package name (hyphens stripped)
	Goreleaser    bool   // Include goreleaser config and release workflow
	Community     bool   // Include community files (issue templates, PR template, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY)
	Changelog     bool   // Include changelog files (CHANGELOG.yaml, CHANGELOG.md, .chlog.yaml)
	AgentMDClaude bool   // Include CLAUDE.md and .skills/
	AgentMDAgents bool   // Include AGENTS.md
	Config        bool   // Include internal/config package and config subcommands
	Todo          bool   // Include TODO.md
}

// Generate walks the embedded template tree and writes rendered files to destDir.
func Generate(cfg Config, destDir string) error {
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	return fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Strip "templates/" prefix
		relPath := strings.TrimPrefix(path, "templates/")
		if relPath == "" {
			return nil
		}

		if d.IsDir() {
			return skipDir(relPath, cfg)
		}

		if skipFile(relPath, cfg) {
			return nil
		}

		// Resolve output path: strip .tmpl, replace placeholders, fix special paths
		outPath := resolveOutputPath(relPath, cfg.BinaryName, cfg.LibPackage, cfg.License)
		fullPath := filepath.Join(destDir, outPath)

		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", outPath, err)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", path, err)
		}

		// Parse and execute template with [% %] delimiters to avoid
		// conflicts with goreleaser {{ }}, Make $, and bash [[ ]]
		tmpl, err := template.New(relPath).Delims("[%", "%]").Funcs(funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", relPath, err)
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", fullPath, err)
		}

		if err := tmpl.Execute(f, cfg); err != nil {
			_ = f.Close()
			return fmt.Errorf("executing template %s: %w", relPath, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("closing file %s: %w", fullPath, err)
		}

		if strings.HasSuffix(outPath, ".sh") {
			if err := os.Chmod(fullPath, 0o755); err != nil {
				return fmt.Errorf("chmod %s: %w", outPath, err)
			}
		}

		return nil
	})
}

// skipDir returns fs.SkipDir for directories that should be excluded based on config.
func skipDir(relPath string, cfg Config) error {
	// CI provider filtering
	if matchesPrefix(relPath, "github") && !cfg.CIGitHub {
		return fs.SkipDir
	}
	if matchesPrefix(relPath, "gitlab") && !cfg.CIGitLab {
		return fs.SkipDir
	}

	// Community files: skip issue templates dir
	if !cfg.Community && matchesPrefix(relPath, "github/ISSUE_TEMPLATE") {
		return fs.SkipDir
	}

	// Agent MD: skip .skills/ when Claude not selected
	if !cfg.AgentMDClaude && matchesPrefix(relPath, "skills") {
		return fs.SkipDir
	}

	// Config: skip internal/config when disabled
	if !cfg.Config && matchesPrefix(relPath, "internal/config") {
		return fs.SkipDir
	}

	// Layout filtering
	if !cfg.HasCLI {
		if matchesPrefix(relPath, "cmd") || matchesPrefix(relPath, "internal") {
			return fs.SkipDir
		}
		if matchesPrefix(relPath, "scripts") {
			return fs.SkipDir
		}
	}
	if !cfg.HasLib && matchesPrefix(relPath, "pkg") {
		return fs.SkipDir
	}

	return nil
}

// skipFile returns true for files that should be excluded based on config.
func skipFile(relPath string, cfg Config) bool {
	// License filtering
	if cfg.License != "mit" && relPath == "LICENSE_mit.tmpl" {
		return true
	}
	if cfg.License != "apache2" && relPath == "LICENSE_apache2.tmpl" {
		return true
	}

	// CLI-only files skipped for lib layout
	if !cfg.HasCLI {
		switch relPath {
		case "install.sh.tmpl", "uninstall.sh.tmpl", "goreleaser.yaml.tmpl":
			return true
		}
	}

	// Community files
	if !cfg.Community {
		switch relPath {
		case "CONTRIBUTING.md.tmpl", "CODE_OF_CONDUCT.md.tmpl", "SECURITY.md.tmpl",
			"github/pull_request_template.md.tmpl":
			return true
		}
	}

	// Changelog-related files
	if !cfg.Changelog {
		switch relPath {
		case "chlog.yaml.tmpl", "chlog-config.yaml.tmpl", "CHANGELOG.md.tmpl":
			return true
		}
	}

	// Goreleaser-related files
	if !cfg.Goreleaser {
		switch relPath {
		case "goreleaser.yaml.tmpl", "scripts/release.sh.tmpl",
			"github/workflows/release.yml.tmpl":
			return true
		}
	}

	// Agent MD files
	if !cfg.AgentMDClaude && relPath == "CLAUDE.md.tmpl" {
		return true
	}
	if !cfg.AgentMDAgents && relPath == "AGENTS.md.tmpl" {
		return true
	}

	// Config command file
	if !cfg.Config && relPath == "cmd/{{BinaryName}}/config.go.tmpl" {
		return true
	}

	// TODO.md
	if !cfg.Todo && relPath == "TODO.md.tmpl" {
		return true
	}

	return false
}

// matchesPrefix returns true if relPath equals prefix or starts with prefix/.
func matchesPrefix(relPath, prefix string) bool {
	return relPath == prefix || strings.HasPrefix(relPath, prefix+"/")
}

// resolveOutputPath converts a template path to the final output path.
func resolveOutputPath(relPath, binaryName, libPackage, license string) string {
	out := strings.TrimSuffix(relPath, ".tmpl")

	out = strings.ReplaceAll(out, "{{BinaryName}}", binaryName)
	out = strings.ReplaceAll(out, "{{LibPackage}}", libPackage)

	// Convert gitkeep → .gitkeep (used for keeping empty directories in git)
	if filepath.Base(out) == "gitkeep" {
		out = filepath.Join(filepath.Dir(out), ".gitkeep")
	}

	switch {
	case strings.HasPrefix(out, "skills/"):
		out = "." + out
	case strings.HasPrefix(out, "github/"):
		out = "." + out
	case strings.HasPrefix(out, "gitlab/"):
		out = "." + strings.TrimPrefix(out, "gitlab/")
	case out == "goreleaser.yaml":
		out = ".goreleaser.yaml"
	case out == "gitignore":
		out = ".gitignore"
	case out == "LICENSE_mit" || out == "LICENSE_apache2":
		out = "LICENSE"
	case out == "chlog.yaml":
		out = "CHANGELOG.yaml"
	case out == "chlog-config.yaml":
		out = ".chlog.yaml"
	}

	return out
}
