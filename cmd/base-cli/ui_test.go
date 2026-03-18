package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func init() {
	// Disable colors in tests for predictable output.
	color.NoColor = true
}

func TestHighlight(t *testing.T) {
	got := highlight("hello")
	if got != "hello" {
		t.Errorf("highlight(%q) = %q, want %q", "hello", got, "hello")
	}
}

func TestFileRef(t *testing.T) {
	got := fileRef("path/to/file")
	if got != "path/to/file" {
		t.Errorf("fileRef(%q) = %q, want %q", "path/to/file", got, "path/to/file")
	}
}

func TestSuccess(t *testing.T) {
	out := captureStdoutUI(t, func() {
		success("done %s", "ok")
	})
	if !strings.Contains(out, "done ok") {
		t.Errorf("success() output = %q, want to contain %q", out, "done ok")
	}
}

func TestWarn(t *testing.T) {
	out := captureStdoutUI(t, func() {
		warn("caution %d", 42)
	})
	if !strings.Contains(out, "caution 42") {
		t.Errorf("warn() output = %q, want to contain %q", out, "caution 42")
	}
}

func captureStdoutUI(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
