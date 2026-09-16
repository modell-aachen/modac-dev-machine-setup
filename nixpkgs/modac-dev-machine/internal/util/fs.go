package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileExists checks if a file or directory exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetTemplatesDir returns the path to the templates directory, checking
// the production path (share/machine/templates) first, then falling back
// to the development repo path (scripts/templates).
func GetTemplatesDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	binDir := filepath.Dir(exePath)

	shareDir := filepath.Join(binDir, "..", "share", "machine", "templates")
	if FileExists(shareDir) {
		return shareDir, nil
	}

	repoTemplatesDir := filepath.Join(binDir, "..", "scripts", "templates")
	if FileExists(repoTemplatesDir) {
		return repoTemplatesDir, nil
	}

	return "", fmt.Errorf("templates directory not found")
}

// GetDevboxGlobalPath returns the global path of devbox by executing
// the command `devbox global path`.
func GetDevboxGlobalPath() (string, error) {
	cmd := exec.Command("devbox", "global", "path")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
