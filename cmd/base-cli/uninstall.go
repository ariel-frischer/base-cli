package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallYes bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove base-cli from your system",
	RunE: func(cmd *cobra.Command, args []string) error {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locating binary: %w", err)
		}
		exe, err = filepath.EvalSymlinks(exe)
		if err != nil {
			return fmt.Errorf("resolving symlinks: %w", err)
		}

		fmt.Printf("Found binary: %s\n", fileRef(exe))

		if !uninstallYes {
			fmt.Print("Remove this binary? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := os.Remove(exe); err != nil {
			return fmt.Errorf("removing binary: %w", err)
		}

		// Clean up backups in the same directory
		dir := filepath.Dir(exe)
		base := filepath.Base(exe)
		pattern := filepath.Join(dir, base+".backup.*")
		backups, _ := filepath.Glob(pattern)
		for _, b := range backups {
			_ = os.Remove(b)
		}

		success("base-cli has been removed from %s", dir)
		if len(backups) > 0 {
			success("Cleaned up %d backup(s)", len(backups))
		}
		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallYes, "yes", false, "Skip confirmation prompt")
}
