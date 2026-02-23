package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:    "context",
	Short:  "Print the current repo context (for shell prompt use)",
	Long:   `Print the name of the active repo — current git root if inside one, otherwise the saved context set by 'wtt repo'. Exits silently with no output when no context is available. Intended for prompt integrations.`,
	Args:   cobra.NoArgs,
	Hidden: true, // power-user / prompt integration tool
	RunE:   runContext,
}

func runContext(_ *cobra.Command, _ []string) error {
	// 1. Prefer the actual git root of the current directory
	root, err := git.RepoRoot()
	if err == nil {
		fmt.Println(filepath.Base(root))
		return nil
	}

	// 2. Fall back to the saved context
	saved, err := globalconfig.GetCurrentRepo()
	if err != nil || saved == "" {
		// No context at all — output nothing, exit 0 (prompt stays clean)
		os.Exit(0)
	}
	fmt.Println(filepath.Base(saved))
	return nil
}
