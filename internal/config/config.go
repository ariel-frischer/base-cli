package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultDir returns ~/.config/base-cli.
func DefaultDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "base-cli")
}

// DefaultPath returns ~/.config/base-cli/config.yaml.
func DefaultPath() string {
	return filepath.Join(DefaultDir(), "config.yaml")
}

// Config holds user-level defaults for base-cli init.
// All fields are optional — zero values mean "use CLI default".
type Config struct {
	Author       string `yaml:"author,omitempty"`
	License      string `yaml:"license,omitempty"`
	CI           string `yaml:"ci,omitempty"`
	Layout       string `yaml:"layout,omitempty"`
	AgentMD      string `yaml:"agent_md,omitempty"`
	NoGitInit    *bool  `yaml:"no_git_init,omitempty"`
	NoGoreleaser *bool  `yaml:"no_goreleaser,omitempty"`
	NoCommunity  *bool  `yaml:"no_community,omitempty"`
	NoChangelog  *bool  `yaml:"no_changelog,omitempty"`
	NoConfig     *bool  `yaml:"no_config,omitempty"`
	Todo         *bool  `yaml:"todo,omitempty"`
}

// Load reads config from path. Returns empty config (no error) if file is missing.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

// Save writes config to path, creating parent directories as needed.
func Save(cfg *Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config %s: %w", path, err)
	}
	return nil
}

// BoolVal returns the value of a *bool, or the fallback if nil.
func BoolVal(p *bool, fallback bool) bool {
	if p != nil {
		return *p
	}
	return fallback
}

// BoolPtr returns a pointer to b.
func BoolPtr(b bool) *bool {
	return &b
}
