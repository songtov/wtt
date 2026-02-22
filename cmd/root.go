package cmd

import (
	"fmt"
	"os"

	"github.com/songtov/wtt/internal/git"
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

	worktrees, err := git.ListWorktrees()
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
