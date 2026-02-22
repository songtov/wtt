package cmd

import (
	"fmt"
	"os"

	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/globalconfig"
	"github.com/songtov/wtt/internal/shell"
	"github.com/spf13/cobra"
)

var initShell string

var rootCmd = &cobra.Command{
	Use:   "wtt [branch]",
	Short: "Git worktree manager",
	Long:  `wtt wraps git worktree to create, remove, and navigate worktrees easily.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&initShell, "init", "", "Print shell init function (zsh, bash, fish)")
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(repoCmd)
	rootCmd.AddCommand(contextCmd)
}

// repoRootWithFallback returns the git repo root for the current directory.
// When not inside a git repo it falls back to the context saved by "wtt repo".
func repoRootWithFallback() (string, error) {
	root, err := git.RepoRoot()
	if err == nil {
		return root, nil
	}
	current, cfgErr := globalconfig.GetCurrentRepo()
	if cfgErr != nil || current == "" {
		return "", fmt.Errorf("not inside a git repository (run 'wtt repo' to set a repo context)")
	}
	return current, nil
}

// autoRegisterRepo silently adds repoRoot to the known-repos list so it shows
// up in "wtt repo". Errors are intentionally ignored — registration is best-effort.
func autoRegisterRepo(repoRoot string) {
	_ = globalconfig.RegisterRepo(repoRoot)
}

func runRoot(cmd *cobra.Command, args []string) error {
	if initShell != "" {
		fn, err := shell.InitFunction(initShell)
		if err != nil {
			return err
		}
		fmt.Print(fn)
		return nil
	}

	if len(args) == 0 {
		return cmd.Help()
	}

	branch := args[0]

	repoRoot, err := repoRootWithFallback()
	if err != nil {
		return err
	}
	autoRegisterRepo(repoRoot)

	worktrees, err := git.ListWorktreesIn(repoRoot)
	if err != nil {
		return fmt.Errorf("listing worktrees: %w", err)
	}

	for _, wt := range worktrees {
		if wt.Branch == branch || wt.Branch == "refs/heads/"+branch {
			fmt.Println(wt.Path)
			return nil
		}
	}

	return fmt.Errorf("no worktree found for branch %q", branch)
}
