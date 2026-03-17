package main

import "testing"

func TestTruncateCommit(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"long hash":    {"abcdef1234567890", "abcdef12"},
		"exactly 8":    {"abcdef12", "abcdef12"},
		"short hash":   {"abc", "abc"},
		"empty":        {"", ""},
		"dev":          {"dev", "dev"},
		"8+ boundary":  {"123456789", "12345678"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := truncateCommit(tt.input)
			if got != tt.want {
				t.Errorf("truncateCommit(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
