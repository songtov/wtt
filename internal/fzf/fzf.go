package fzf

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// SelectRepo presents a list of repo paths for selection via fzf (or numbered
// fallback). Returns the selected absolute repo path, or "" if cancelled.
func SelectRepo(repos []string) (string, error) {
	if len(repos) == 0 {
		return "", fmt.Errorf("no repos registered yet; run wtt commands inside a git repo first")
	}
	if hasFzf() {
		return selectRepoWithFzf(repos)
	}
	return selectRepoNumbered(repos)
}

func selectRepoWithFzf(repos []string) (string, error) {
	var input strings.Builder
	for i, p := range repos {
		fmt.Fprintf(&input, "%d\t%s\t%s\n", i, filepath.Base(p), p)
	}

	// Show "name   path" columns; hide the index column
	cmd := exec.Command("fzf", "--with-nth=2,3", "--delimiter=\t", "--ansi")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.ExitCode() == 130 {
			return "", nil // user cancelled
		}
		return "", fmt.Errorf("fzf: %w", err)
	}

	selected := strings.TrimSpace(out.String())
	parts := strings.SplitN(selected, "\t", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected fzf output: %q", selected)
	}
	idx, err := strconv.Atoi(parts[0])
	if err != nil || idx < 0 || idx >= len(repos) {
		return "", fmt.Errorf("unexpected fzf index: %q", parts[0])
	}
	return repos[idx], nil
}

// SelectRepoWithNone is like SelectRepo but prepends a "None" option that
// represents clearing the current repo context.
// Returns (path, noneSelected, err). noneSelected is true when the user
// explicitly chose "(none)"; path is "" and noneSelected is false when cancelled.
func SelectRepoWithNone(repos []string) (string, bool, error) {
	if hasFzf() {
		return selectRepoWithNoneFzf(repos)
	}
	return selectRepoWithNoneNumbered(repos)
}

func selectRepoWithNoneFzf(repos []string) (string, bool, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "-1\t(none)\t\n")
	for i, p := range repos {
		fmt.Fprintf(&input, "%d\t%s\t%s\n", i, filepath.Base(p), p)
	}

	cmd := exec.Command("fzf", "--with-nth=2,3", "--delimiter=\t", "--ansi")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.ExitCode() == 130 {
			return "", false, nil
		}
		return "", false, fmt.Errorf("fzf: %w", err)
	}

	selected := strings.TrimSpace(out.String())
	parts := strings.SplitN(selected, "\t", 3)
	if len(parts) < 2 {
		return "", false, fmt.Errorf("unexpected fzf output: %q", selected)
	}
	if parts[0] == "-1" {
		return "", true, nil
	}
	idx, err := strconv.Atoi(parts[0])
	if err != nil || idx < 0 || idx >= len(repos) {
		return "", false, fmt.Errorf("unexpected fzf index: %q", parts[0])
	}
	return repos[idx], false, nil
}

func selectRepoWithNoneNumbered(repos []string) (string, bool, error) {
	fmt.Fprintln(os.Stderr, "Select a repo (0 to clear):")
	fmt.Fprintf(os.Stderr, "  [0] (none) – clear repo context\n")
	for i, p := range repos {
		fmt.Fprintf(os.Stderr, "  [%d] %s  %s\n", i+1, filepath.Base(p), p)
	}
	fmt.Fprint(os.Stderr, "Enter number: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return "", false, nil
	}
	n, err := strconv.Atoi(text)
	if err != nil || n < 0 || n > len(repos) {
		return "", false, fmt.Errorf("invalid selection %q", text)
	}
	if n == 0 {
		return "", true, nil
	}
	return repos[n-1], false, nil
}

func selectRepoNumbered(repos []string) (string, error) {
	fmt.Fprintln(os.Stderr, "Select a repo:")
	for i, p := range repos {
		fmt.Fprintf(os.Stderr, "  [%d] %s  %s\n", i+1, filepath.Base(p), p)
	}
	fmt.Fprint(os.Stderr, "Enter number: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return "", nil
	}
	n, err := strconv.Atoi(text)
	if err != nil || n < 1 || n > len(repos) {
		return "", fmt.Errorf("invalid selection %q", text)
	}
	return repos[n-1], nil
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
