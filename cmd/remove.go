package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/songtov/wtt/internal/fzf"
	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/worktree"
	"github.com/spf13/cobra"
)

var forceRemove bool

var removeCmd = &cobra.Command{
	Use:   "remove [branch]",
	Short: "Remove a worktree",
	Long:  `Remove the git worktree associated with the given branch name. If no branch is given, opens an interactive picker.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Skip confirmation prompt")
}

func runRemove(_ *cobra.Command, args []string) error {
	repoRoot, err := repoRootWithFallback()
	if err != nil {
		return err
	}
	autoRegisterRepo(repoRoot)

	worktrees, err := git.ListWorktreesIn(repoRoot)
	if err != nil {
		return fmt.Errorf("listing worktrees: %w", err)
	}

	// Skip the main worktree (first entry) from removal candidates
	removable := worktrees[1:]

	var targetPath, branch string

	if len(args) == 0 {
		// No branch given: open interactive picker
		if len(removable) == 0 {
			return fmt.Errorf("no worktrees to remove")
		}
		selected, err := fzf.SelectWorktree(removable, false)
		if err != nil {
			return err
		}
		if selected == nil {
			return nil // user cancelled
		}
		targetPath = selected.Path
		branch = strings.TrimPrefix(selected.Branch, "refs/heads/")
	} else {
		branch = args[0]
		for _, wt := range removable {
			if wt.Branch == branch || wt.Branch == "refs/heads/"+branch {
				targetPath = wt.Path
				break
			}
		}
		if targetPath == "" {
			return fmt.Errorf("no worktree found for branch %q", branch)
		}
	}

	if !forceRemove {
		fmt.Fprintf(os.Stderr, "Remove worktree at %s? [y/N] ", targetPath)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(os.Stderr, "Aborted.")
			return nil
		}
	}

	if err := worktree.Remove(repoRoot, targetPath, forceRemove); err != nil {
		return err
	}

	// Clean up empty parent directory
	parent := filepath.Dir(targetPath)
	entries, err := os.ReadDir(parent)
	if err == nil && len(entries) == 0 {
		_ = os.Remove(parent)
	}

	fmt.Fprintf(os.Stderr, "Removed worktree for branch %q\n", branch)
	return nil
}
