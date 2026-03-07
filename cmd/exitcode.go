package cmd

import "strings"

const (
	exitCodeGeneric    = 1
	exitCodeValidation = 3
	exitCodeEnv        = 4
	exitCodeGit        = 5
)

func exitCodeForError(err error) int {
	if err == nil {
		return 0
	}
	msg := strings.ToLower(err.Error())

	// Validation / argument errors.
	if strings.Contains(msg, "requires") ||
		strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "cannot remove the main worktree") ||
		strings.Contains(msg, "no worktree found for branch") {
		return exitCodeValidation
	}

	// Environment / setup errors.
	if strings.Contains(msg, "not inside a git repository") ||
		strings.Contains(msg, "run 'wtt repo' to set a repo context") {
		return exitCodeEnv
	}

	// Git command failures.
	if strings.Contains(msg, "git worktree") ||
		strings.Contains(msg, "git check-ref-format") ||
		strings.HasPrefix(msg, "listing worktrees:") {
		return exitCodeGit
	}

	return exitCodeGeneric
}
