package worktree

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/songtov/wtt/internal/git"
)

// Create creates a new git worktree for the given branch and returns its path.
// worktreeBaseDir is the absolute path of the base directory for worktrees.
func Create(repoRoot, worktreeBaseDir, branch string) (string, error) {
	safeName := git.BranchToPath(branch)
	worktreePath := filepath.Join(worktreeBaseDir, safeName)

	cmd := exec.Command("git", "worktree", "add", "-b", branch, worktreePath)
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
