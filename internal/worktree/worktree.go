package worktree

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/songtov/wtt/internal/git"
)

// Create creates a new git worktree for the given branch and returns its path.
// worktreeBaseDir is the absolute path of the base directory for worktrees.
// base, if non-empty, is passed as the start-point to git worktree add.
func Create(repoRoot, worktreeBaseDir, branch, base string) (string, error) {
	safeName := git.BranchToPath(branch)
	worktreePath := filepath.Join(worktreeBaseDir, safeName)

	args := []string{"worktree", "add", "-b", branch, worktreePath}
	if base != "" {
		args = append(args, base)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git worktree add: %w\n%s", err, out)
	}

	return worktreePath, nil
}

// Remove removes the git worktree at the given path.
// force skips the git-level check for modified/untracked files.
func Remove(repoRoot, worktreePath string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, worktreePath)

	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree remove: %w\n%s", err, out)
	}
	return nil
}
