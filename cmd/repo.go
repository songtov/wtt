package cmd

import (
	"fmt"

	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repository contexts",
	Long: `Switch between repositories and manage the active repo context.
All subsequent wtt commands (create, list, remove) will operate on the
active repository even when run from outside it — just like kubens for kubectl.`,
}

// canonicalRepos loads the known repos, normalizes each to its main repo root,
// deduplicates, and saves the cleaned list back if it changed.
func canonicalRepos() ([]string, error) {
	repos, err := globalconfig.GetKnownRepos()
	if err != nil {
		return nil, fmt.Errorf("loading repos: %w", err)
	}

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
	return canonical, nil
}
