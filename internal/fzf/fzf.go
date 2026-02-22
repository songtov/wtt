package fzf

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/songtov/wtt/internal/git"
)

// Select presents a list of worktrees for selection.
// If fzf is available it uses an interactive picker; otherwise a numbered list.
// Returns the selected worktree path, or "" if the user cancelled.
func Select(worktrees []git.Worktree) (string, error) {
	wt, err := SelectWorktree(worktrees)
	if err != nil || wt == nil {
		return "", err
	}
	return wt.Path, nil
}

// SelectWorktree presents a list of worktrees for selection and returns the
// selected Worktree, or nil if the user cancelled.
func SelectWorktree(worktrees []git.Worktree) (*git.Worktree, error) {
	if len(worktrees) == 0 {
		return nil, fmt.Errorf("no worktrees found")
	}

	if hasFzf() {
		return selectWorktreeWithFzf(worktrees)
	}
	return selectWorktreeNumbered(worktrees)
}

func hasFzf() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

func selectWorktreeWithFzf(worktrees []git.Worktree) (*git.Worktree, error) {
	var input strings.Builder
	for i, wt := range worktrees {
		branch := wt.Branch
		if branch == "" {
			branch = "(detached)"
		}
		branch = strings.TrimPrefix(branch, "refs/heads/")
		fmt.Fprintf(&input, "%d\t%s\n", i, branch)
	}

	cmd := exec.Command("fzf", "--with-nth=2", "--delimiter=\t", "--ansi")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.ExitCode() == 130 {
			return nil, nil
		}
		return nil, fmt.Errorf("fzf: %w", err)
	}

	selected := strings.TrimSpace(out.String())
	parts := strings.SplitN(selected, "\t", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected fzf output: %q", selected)
	}
	idx, err := strconv.Atoi(parts[0])
	if err != nil || idx < 0 || idx >= len(worktrees) {
		return nil, fmt.Errorf("unexpected fzf index: %q", parts[0])
	}
	return &worktrees[idx], nil
}

func selectWorktreeNumbered(worktrees []git.Worktree) (*git.Worktree, error) {
	fmt.Fprintln(os.Stderr, "Select a worktree:")
	for i, wt := range worktrees {
		branch := wt.Branch
		if branch == "" {
			branch = "(detached)"
		}
		branch = strings.TrimPrefix(branch, "refs/heads/")
		fmt.Fprintf(os.Stderr, "  [%d] %s  %s\n", i+1, branch, wt.Path)
	}
	fmt.Fprint(os.Stderr, "Enter number: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return nil, nil
	}

	n, err := strconv.Atoi(text)
	if err != nil || n < 1 || n > len(worktrees) {
		return nil, fmt.Errorf("invalid selection %q", text)
	}
	return &worktrees[n-1], nil
}
