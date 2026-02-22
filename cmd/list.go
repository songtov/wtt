package cmd

import (
	"fmt"

	"github.com/songtov/wtt/internal/fzf"
	"github.com/songtov/wtt/internal/git"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and navigate worktrees",
	Long:  `List all worktrees and interactively select one to navigate to (requires fzf).`,
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return fmt.Errorf("listing worktrees: %w", err)
	}

	path, err := fzf.Select(worktrees)
	if err != nil {
		return err
	}
	if path == "" {
		return nil // user cancelled
	}

	// Print path so the shell wrapper can cd to it
	fmt.Println(path)
	return nil
}
