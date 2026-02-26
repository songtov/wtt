package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/songtov/wtt/internal/fzf"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/spf13/cobra"
)

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Pick and switch the active repository context",
	Args:  cobra.NoArgs,
	RunE:  runRepoList,
}

func init() {
	repoCmd.AddCommand(repoListCmd)
}

func runRepoList(_ *cobra.Command, _ []string) error {
	repos, err := canonicalRepos()
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Fprintln(os.Stderr, "No repos registered yet. Run any wtt command from inside a git repo first.")
		return nil
	}

	selected, noneSelected, err := fzf.SelectRepoWithNone(repos)
	if err != nil {
		return err
	}

	if noneSelected {
		if err := globalconfig.ClearCurrentRepo(); err != nil {
			return fmt.Errorf("clearing repo context: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Cleared repo context.")
		return nil
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
