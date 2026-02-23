package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/songtov/wtt/internal/fzf"
	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Select the active repository context",
	Long: `Pick a repository with fzf and remember it as the active context.
All subsequent wtt commands (create, list, remove) will operate on that
repository even when run from outside it — just like kubens for kubectl.`,
	Args: cobra.NoArgs,
	RunE: runRepo,
}

func runRepo(_ *cobra.Command, _ []string) error {
	repos, err := globalconfig.GetKnownRepos()
	if err != nil {
		return fmt.Errorf("loading repos: %w", err)
	}
	if len(repos) == 0 {
		fmt.Fprintln(os.Stderr, "No repos registered yet. Run any wtt command from inside a git repo first.")
		return nil
	}

	// Normalize each entry to its main repo root and deduplicate.
	// This cleans up stale worktree paths registered before the fix.
	seen := map[string]bool{}
	var canonical []string
	for _, r := range repos {
		main, err := git.MainRepoRootOf(r)
		if err != nil {
			continue // path gone or not a git repo — drop it
		}
		if !seen[main] {
			seen[main] = true
			canonical = append(canonical, main)
		}
	}
	if len(canonical) != len(repos) {
		_ = globalconfig.SetKnownRepos(canonical)
	}
	repos = canonical

	if len(repos) == 0 {
		fmt.Fprintln(os.Stderr, "No repos registered yet. Run any wtt command from inside a git repo first.")
		return nil
	}

	selected, err := fzf.SelectRepo(repos)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil // user cancelled
	}

	if err := globalconfig.SetCurrentRepo(selected); err != nil {
		return fmt.Errorf("saving repo context: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Switched to repo: %s\n", filepath.Base(selected))
	return nil
}
