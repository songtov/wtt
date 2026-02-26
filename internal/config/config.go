package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const configFile = ".wtt.toml"

// Config holds the wtt configuration.
type Config struct {
	WorktreeDir  string   `toml:"worktree_dir"`
	CopyFiles    []string `toml:"copy_files"`
	CopyDirs     []string `toml:"copy_dirs"`
	SymlinkFiles []string `toml:"symlink_files"`
	PostCreate   []string `toml:"post_create"`
}

// Load reads .wtt.toml from repoRoot and merges with defaults.
func Load(repoRoot, repoName string) (*Config, error) {
	cfg := defaults(repoName)

	path := filepath.Join(repoRoot, configFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	var fileCfg Config
	if _, err := toml.DecodeFile(path, &fileCfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	// Merge: non-zero file values override defaults
	if fileCfg.WorktreeDir != "" {
		cfg.WorktreeDir = fileCfg.WorktreeDir
	}
	if len(fileCfg.CopyFiles) > 0 {
		cfg.CopyFiles = fileCfg.CopyFiles
	}
	if len(fileCfg.CopyDirs) > 0 {
		cfg.CopyDirs = fileCfg.CopyDirs
	}
	if len(fileCfg.SymlinkFiles) > 0 {
		cfg.SymlinkFiles = fileCfg.SymlinkFiles
	}
	if len(fileCfg.PostCreate) > 0 {
		cfg.PostCreate = fileCfg.PostCreate
	}

	return cfg, nil
}

func defaults(repoName string) *Config {
	return &Config{
		WorktreeDir:  fmt.Sprintf("../%s-worktrees", repoName),
		CopyFiles:    []string{".gitignore"},
		CopyDirs:     []string{},
		SymlinkFiles: []string{},
		PostCreate:   []string{},
	}
}
