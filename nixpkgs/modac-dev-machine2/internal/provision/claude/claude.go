package claude

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/util"
)

// Run sets up Claude configuration
func Run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")

	// Create .claude directory
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	claudeMdPath := filepath.Join(claudeDir, "CLAUDE.md")

	// Check if CLAUDE.md already exists
	if util.FileExists(claudeMdPath) {
		return nil
	}

	// Get templates directory
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	// Copy template
	templatePath := filepath.Join(templatesDir, "team-claude.md")
	if err := copyFile(templatePath, claudeMdPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Set permissions
	if err := os.Chmod(claudeMdPath, 0664); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	fmt.Printf("Created %s\n", claudeMdPath)
	return nil
}

func getTemplatesDir() (string, error) {
	// Get the directory of the currently running binary
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Binary is in bin/, templates should be in ../share/machine2/templates/
	binDir := filepath.Dir(exePath)
	shareDir := filepath.Join(binDir, "..", "share", "machine2", "templates")

	// Check if the share directory exists
	if util.FileExists(shareDir) {
		return shareDir, nil
	}

	// Fallback: check if we're running from the repo (development mode)
	repoTemplatesDir := filepath.Join(binDir, "..", "scripts", "templates")
	if util.FileExists(repoTemplatesDir) {
		return repoTemplatesDir, nil
	}

	return "", fmt.Errorf("templates directory not found")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
