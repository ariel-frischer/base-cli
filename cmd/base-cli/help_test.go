package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSplitFlagLine(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"flag with description": {
			input: "--name string   Project name",
			want:  []string{"--name string", "Project name"},
		},
		"short and long flag": {
			input: "-n, --name string   Project name",
			want:  []string{"-n, --name string", "Project name"},
		},
		"flag only no description": {
			input: "--verbose",
			want:  []string{"--verbose"},
		},
		"empty string": {
			input: "",
			want:  []string{""},
		},
		"flag with wide padding": {
			input: "--ci string        CI provider (default \"github\")",
			want:  []string{"--ci string", "CI provider (default \"github\")"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := splitFlagLine(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("splitFlagLine(%q) returned %d parts, want %d: %v", tt.input, len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitFlagLine(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestVisibleSubcommands(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	visible := &cobra.Command{Use: "init", Short: "Initialize"}
	hidden := &cobra.Command{Use: "secret", Hidden: true}
	helpCmd := &cobra.Command{Use: "help", Short: "Help"}

	root.AddCommand(visible, hidden, helpCmd)

	got := visibleSubcommands(root)
	if len(got) != 1 {
		t.Fatalf("visibleSubcommands() returned %d commands, want 1", len(got))
	}
	if got[0].Name() != "init" {
		t.Errorf("visibleSubcommands()[0] = %q, want %q", got[0].Name(), "init")
	}
}

func TestMaxCommandNameLen(t *testing.T) {
	cmds := []*cobra.Command{
		{Use: "a"},
		{Use: "longer"},
		{Use: "mid"},
	}

	got := maxCommandNameLen(cmds)
	if got != 6 {
		t.Errorf("maxCommandNameLen() = %d, want 6", got)
	}
}

func TestMaxCommandNameLenEmpty(t *testing.T) {
	got := maxCommandNameLen(nil)
	if got != 0 {
		t.Errorf("maxCommandNameLen(nil) = %d, want 0", got)
	}
}
