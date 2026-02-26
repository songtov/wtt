package globalconfig

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "wtt")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// GetCurrentRepo returns the path of the currently selected repo context.
// Returns an empty string (no error) if no context has been set yet.
func GetCurrentRepo() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(dir, "current_repo"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SetCurrentRepo saves the given repo path as the active context.
func SetCurrentRepo(repoPath string) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "current_repo"), []byte(repoPath+"\n"), 0644)
}

// GetKnownRepos returns all registered repo paths.
func GetKnownRepos() ([]string, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "repos"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var repos []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			repos = append(repos, line)
		}
	}
	return repos, nil
}

// SetKnownRepos overwrites the repos list with the given slice.
func SetKnownRepos(repos []string) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	var sb strings.Builder
	for _, r := range repos {
		sb.WriteString(r + "\n")
	}
	return os.WriteFile(filepath.Join(dir, "repos"), []byte(sb.String()), 0644)
}

// RemoveRepo removes a repo path from the known repos list.
func RemoveRepo(repoPath string) error {
	repos, err := GetKnownRepos()
	if err != nil {
		return err
	}
	var filtered []string
	for _, r := range repos {
		if r != repoPath {
			filtered = append(filtered, r)
		}
	}
	return SetKnownRepos(filtered)
}

// ClearCurrentRepo removes the saved current-repo context file.
func ClearCurrentRepo() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(dir, "current_repo"))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// RegisterRepo adds a repo path to the known repos list. It is idempotent.
func RegisterRepo(repoPath string) error {
	repos, err := GetKnownRepos()
	if err != nil {
		return err
	}
	for _, r := range repos {
		if r == repoPath {
			return nil // already registered
		}
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "repos"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(repoPath + "\n")
	return err
}
