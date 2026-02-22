package worktree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFiles copies the listed files from srcDir to dstDir.
// Files that don't exist in srcDir are silently skipped.
func CopyFiles(srcDir, dstDir string, files []string) error {
	for _, f := range files {
		src := filepath.Join(srcDir, f)
		dst := filepath.Join(dstDir, f)
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

// CopyDirs copies listed directories from srcDir to dstDir recursively.
func CopyDirs(srcDir, dstDir string, dirs []string) error {
	for _, d := range dirs {
		src := filepath.Join(srcDir, d)
		dst := filepath.Join(dstDir, d)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		if err := copyDir(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // silently skip
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s → %s: %w", src, dst, err)
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}
