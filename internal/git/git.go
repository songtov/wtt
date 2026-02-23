package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree entry.
type Worktree struct {
	Path   string
	Head   string
	Branch string
}

// RepoRoot returns the absolute path of the repository root.
func RepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git repository")
	}
	return strings.TrimSpace(string(out)), nil
}

// MainRepoRoot returns the absolute path of the main (primary) repository
// root, even when called from inside a linked worktree. It uses
// ListWorktrees(), whose first entry is always the main worktree regardless
// of the current working directory.
func MainRepoRoot() (string, error) {
	worktrees, err := ListWorktrees()
	if err != nil {
		return "", fmt.Errorf("not inside a git repository: %w", err)
	}
	if len(worktrees) == 0 {
		return "", fmt.Errorf("could not determine main repository root")
	}
	return worktrees[0].Path, nil
}

// RepoName returns the base name of the repository root directory.
func RepoName() (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

// ListWorktrees returns all worktrees for the current repo.
func ListWorktrees() ([]Worktree, error) {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return parseWorktrees(string(out)), nil
}

func parseWorktrees(raw string) []Worktree {
	var worktrees []Worktree
	var current Worktree

	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
			}
			current = Worktree{}
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Head = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
		}
	}
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}
	return worktrees
}

// ValidateBranchName checks that a branch name is safe for git.
func ValidateBranchName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	// Use git to validate
	if err := exec.Command("git", "check-ref-format", "--branch", name).Run(); err != nil {
		return fmt.Errorf("invalid branch name %q", name)
	}
	return nil
}

// BranchToPath converts a branch name to a safe directory name.
// feature/login → feature-login
func BranchToPath(branch string) string {
	return strings.ReplaceAll(branch, "/", "-")
}

// MainRepoRootOf resolves any path inside a git repo (including linked
// worktrees) to the primary worktree root. It uses `git -C <dir>` so it
// works regardless of the process's current working directory.
func MainRepoRootOf(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %s", dir)
	}
	worktrees := parseWorktrees(string(out))
	if len(worktrees) == 0 {
		return "", fmt.Errorf("could not determine main repo root for %s", dir)
	}
	return worktrees[0].Path, nil
}

// ListWorktreesIn returns all worktrees for the repo at the given root path.
// It uses `git -C <repoRoot>` so it works from any working directory.
func ListWorktreesIn(repoRoot string) ([]Worktree, error) {
	out, err := exec.Command("git", "-C", repoRoot, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return parseWorktrees(string(out)), nil
}
