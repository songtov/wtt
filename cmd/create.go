package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/songtov/wtt/internal/config"
	"github.com/songtov/wtt/internal/git"
	"github.com/songtov/wtt/internal/namegen"
	"github.com/songtov/wtt/internal/worktree"
	"github.com/spf13/cobra"
)

var createBase string
var createClaude bool

func init() {
	createCmd.Flags().StringVarP(&createBase, "base", "b", "", "Base commit/branch/ref to create the worktree from (default: HEAD)")
	createCmd.Flags().BoolVarP(&createClaude, "claude", "c", false, "start Claude Code in the new worktree")
}

var createCmd = &cobra.Command{
	Use:   "create [branch]",
	Short: "Create a new worktree",
	Long: `Create a new git worktree for the given branch.
If no branch name is given, a random name is generated.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCreate,
}

func runCreate(cmd *cobra.Command, args []string) error {
	repoRoot, err := repoRootWithFallback()
	if err != nil {
		return err
	}
	autoRegisterRepo(repoRoot)

	repoName := filepath.Base(repoRoot)

	cfg, err := config.Load(repoRoot, repoName)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var branch string
	if len(args) == 0 {
		branch = namegen.Generate(repoName)
	} else {
		branch = args[0]
		if err := git.ValidateBranchName(branch); err != nil {
			return err
		}
	}

	// Resolve worktree base dir relative to repoRoot
	worktreeBaseDir := cfg.WorktreeDir
	if !filepath.IsAbs(worktreeBaseDir) {
		worktreeBaseDir = filepath.Clean(filepath.Join(repoRoot, worktreeBaseDir))
	}

	fmt.Fprintf(os.Stderr, "Creating worktree for branch %q...\n", branch)

	worktreePath, err := worktree.Create(repoRoot, worktreeBaseDir, branch, createBase)
	if err != nil {
		return err
	}

	// Copy files
	if err := worktree.CopyFiles(repoRoot, worktreePath, cfg.CopyFiles); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: copying files: %v\n", err)
	}
	if err := worktree.CopyDirs(repoRoot, worktreePath, cfg.CopyDirs); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: copying dirs: %v\n", err)
	}
	if err := worktree.SymlinkFiles(repoRoot, worktreePath, cfg.SymlinkFiles); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: symlinking files: %v\n", err)
	}

	// Run post_create commands
	for _, command := range cfg.PostCreate {
		fmt.Fprintf(os.Stderr, "Running: %s\n", command)
		c := exec.Command("sh", "-c", command)
		c.Dir = worktreePath
		c.Stdout = os.Stderr
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: post_create command failed: %v\n", err)
		}
	}

	// Launch Claude Code in the worktree if requested
	if createClaude {
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("cannot open terminal: %w", err)
		}
		defer tty.Close()

		claude := exec.Command("claude")
		claude.Dir = worktreePath
		claude.Stdin = tty
		claude.Stdout = tty
		claude.Stderr = tty
		if err := claude.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "claude exited: %v\n", err)
		}
	}

	// Print the path so the shell wrapper can cd to it
	fmt.Println(worktreePath)
	return nil
}
