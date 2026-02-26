package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/songtov/wtt/internal/fzf"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/spf13/cobra"
)

var repoRemoveForce bool

var repoRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a repository from the known repos list",
	Args:  cobra.NoArgs,
	RunE:  runRepoRemove,
}

func init() {
	repoRemoveCmd.Flags().BoolVarP(&repoRemoveForce, "force", "f", false, "Skip confirmation prompt")
	repoCmd.AddCommand(repoRemoveCmd)
}

func runRepoRemove(_ *cobra.Command, _ []string) error {
	repos, err := canonicalRepos()
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Fprintln(os.Stderr, "No repos registered yet.")
		return nil
	}

	selected, err := fzf.SelectRepo(repos)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil // user cancelled
	}

	if !repoRemoveForce {
		fmt.Fprintf(os.Stderr, "Remove %s from known repos? [y/N] ", filepath.Base(selected))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			fmt.Fprintln(os.Stderr, "Aborted.")
			return nil
		}
	}

	if err := globalconfig.RemoveRepo(selected); err != nil {
		return fmt.Errorf("removing repo: %w", err)
	}

	current, _ := globalconfig.GetCurrentRepo()
	if current == selected {
		_ = globalconfig.ClearCurrentRepo()
		fmt.Fprintf(os.Stderr, "Removed %s and cleared repo context.\n", filepath.Base(selected))
	} else {
		fmt.Fprintf(os.Stderr, "Removed %s from known repos.\n", filepath.Base(selected))
	}
	return nil
}
