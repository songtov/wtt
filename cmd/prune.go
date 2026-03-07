package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/worktree"
	"github.com/spf13/cobra"
)

var pruneForce bool
var pruneJSON bool

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove worktrees whose branches have been merged",
	Long:  `Prune automatically removes linked worktrees whose branches have been merged into the main worktree's HEAD.`,
	Args:  cobra.NoArgs,
	RunE:  runPrune,
}

func init() {
	pruneCmd.Flags().BoolVarP(&pruneForce, "force", "f", false, "Skip confirmation prompt")
	pruneCmd.Flags().BoolVar(&pruneJSON, "json", false, "Print result as JSON (non-interactive, requires --force)")
}

type pruneEntry struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
}

func runPrune(_ *cobra.Command, _ []string) error {
	if pruneJSON && !pruneForce {
		return fmt.Errorf("--json requires --force to avoid interactive confirmation")
	}

	repoRoot, err := repoRootWithFallback()
	if err != nil {
		return err
	}
	autoRegisterRepo(repoRoot)

	worktrees, err := git.ListWorktreesIn(repoRoot)
	if err != nil {
		return fmt.Errorf("listing worktrees: %w", err)
	}

	if len(worktrees) <= 1 {
		fmt.Fprintln(os.Stderr, "No linked worktrees found.")
		return nil
	}

	mergedBranches, err := git.MergedBranches(repoRoot)
	if err != nil {
		return err
	}

	var toRemove []git.Worktree
	for _, wt := range worktrees[1:] {
		if wt.Branch != "" && mergedBranches[wt.Branch] {
			toRemove = append(toRemove, wt)
		}
	}

	if len(toRemove) == 0 {
		fmt.Fprintln(os.Stderr, "No merged worktrees to prune.")
		return nil
	}

	removed := []pruneEntry{}
	scanner := bufio.NewScanner(os.Stdin)
	for _, wt := range toRemove {
		branch := strings.TrimPrefix(wt.Branch, "refs/heads/")

		if !pruneForce {
			fmt.Fprintf(os.Stderr, "Remove worktree at %s (branch %q, merged)? [y/N] ", wt.Path, branch)
			scanner.Scan()
			answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(os.Stderr, "Skipped.")
				continue
			}
		}

		if err := worktree.Remove(repoRoot, wt.Path, true); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove %s: %v\n", wt.Path, err)
			continue
		}

		// Clean up empty parent directory
		parent := filepath.Dir(wt.Path)
		entries, err := os.ReadDir(parent)
		if err == nil && len(entries) == 0 {
			_ = os.Remove(parent)
		}

		removed = append(removed, pruneEntry{Path: wt.Path, Branch: branch})
		if !pruneJSON {
			fmt.Fprintf(os.Stderr, "Removed worktree for branch %q\n", branch)
		}
	}

	if pruneJSON {
		return printJSON(struct {
			Removed []pruneEntry `json:"removed"`
		}{Removed: removed})
	}

	return nil
}
