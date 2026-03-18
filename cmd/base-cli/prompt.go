package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// stdinReader is a shared buffered reader for interactive prompts.
// Using a single reader avoids data loss when bufio reads ahead from a pipe.
var stdinReader *bufio.Reader

func getStdinReader() *bufio.Reader {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	return stdinReader
}

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func prompt(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("  %s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("  %s: ", label)
	}
	answer, _ := getStdinReader().ReadString('\n')
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return defaultVal
	}
	return answer
}

// promptChoice prompts for a value that must be one of the allowed choices.
// Re-prompts on invalid input.
func promptChoice(label string, choices []string, defaultVal string) string {
	hint := strings.Join(choices, "/")
	for {
		val := prompt(fmt.Sprintf("%s (%s)", label, hint), defaultVal)
		if slices.Contains(choices, val) {
			return val
		}
		warn("  invalid choice %q — must be one of: %s", val, hint)
	}
}

func gitConfigValue(key string) string {
	out, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// detectGitUser tries gh, then glab, then git config to find a username.
func detectGitUser() string {
	if u := ghUser(); u != "" {
		return u
	}
	if u := glabUser(); u != "" {
		return u
	}
	return gitConfigValue("user.name")
}

// detectHost returns "gitlab.com" if glab is authenticated but gh is not,
// otherwise "github.com".
func detectHost() string {
	if ghUser() != "" {
		return "github.com"
	}
	if glabUser() != "" {
		return "gitlab.com"
	}
	return "github.com"
}

func ghUser() string {
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func glabUser() string {
	out, err := exec.Command("glab", "api", "user", "--jq", ".username").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}
