package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const wttTomlTemplate = `# wtt configuration
# See: https://github.com/songtov/wtt

# Directory for worktrees (relative to repo root)
# worktree_dir = "../<reponame>-worktrees"

# Files to copy into new worktrees
copy_files = [".env", ".gitignore"]

# Directories to copy into new worktrees
# copy_dirs = []

# Files to symlink (shared with main repo) into new worktrees
symlink_files = [".claude/settings.local.json"]

# Commands to run after creating a worktree
# post_create = []
`

var initForce bool

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing .wtt.toml")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .wtt.toml configuration file",
	Long:  `Scaffold a .wtt.toml configuration file in the repo root.`,
	Args:  cobra.NoArgs,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	repoRoot, err := repoRootWithFallback()
	if err != nil {
		return err
	}

	dest := filepath.Join(repoRoot, ".wtt.toml")

	if _, err := os.Stat(dest); err == nil && !initForce {
		fmt.Fprintf(os.Stderr, "`.wtt.toml` already exists. Use `wtt init -f` to overwrite.\n")
		return nil
	}

	if err := os.WriteFile(dest, []byte(wttTomlTemplate), 0o644); err != nil {
		return fmt.Errorf("writing .wtt.toml: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Created %s\n", dest)
	return nil
}
