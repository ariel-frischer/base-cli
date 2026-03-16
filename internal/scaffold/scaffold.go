// Package scaffold generates Go CLI project scaffolds from embedded templates.
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
	ProjectName string
	ModulePath  string
	BinaryName  string
	Description string
	Author      string
	Year        string
	GitUser     string
	RepoURL     string
	CIGitHub    bool
	CIGitLab    bool
	EnvPrefix   string
	License     string // "mit", "apache2", "none"
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

		// Conditional directory skipping
		if d.IsDir() {
			if (relPath == "github" || strings.HasPrefix(relPath, "github/")) && !cfg.CIGitHub {
				return fs.SkipDir
			}
			if (relPath == "gitlab" || strings.HasPrefix(relPath, "gitlab/")) && !cfg.CIGitLab {
				return fs.SkipDir
			}
			return nil
		}

		// Skip license files that don't match selection
		if cfg.License != "mit" && relPath == "LICENSE_mit.tmpl" {
			return nil
		}
		if cfg.License != "apache2" && relPath == "LICENSE_apache2.tmpl" {
			return nil
		}
		if cfg.License == "none" && (relPath == "LICENSE_mit.tmpl" || relPath == "LICENSE_apache2.tmpl") {
			return nil
		}

		// Resolve output path: strip .tmpl, replace {{BinaryName}}, fix special paths
		outPath := resolveOutputPath(relPath, cfg.BinaryName, cfg.License)
		fullPath := filepath.Join(destDir, outPath)

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", outPath, err)
		}

		// Read template content
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
		defer f.Close()

		if err := tmpl.Execute(f, cfg); err != nil {
			return fmt.Errorf("executing template %s: %w", relPath, err)
		}

		// Make shell scripts executable
		if strings.HasSuffix(outPath, ".sh") {
			if err := os.Chmod(fullPath, 0o755); err != nil {
				return fmt.Errorf("chmod %s: %w", outPath, err)
			}
		}

		return nil
	})
}

// resolveOutputPath converts a template path to the final output path.
func resolveOutputPath(relPath, binaryName, license string) string {
	// Strip .tmpl extension
	out := strings.TrimSuffix(relPath, ".tmpl")

	// Replace {{BinaryName}} placeholder in directory/file names
	out = strings.ReplaceAll(out, "{{BinaryName}}", binaryName)

	// Special path mappings
	switch {
	case strings.HasPrefix(out, "github/"):
		out = "." + out // github/ -> .github/
	case strings.HasPrefix(out, "gitlab/"):
		// gitlab/gitlab-ci.yml -> .gitlab-ci.yml
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
