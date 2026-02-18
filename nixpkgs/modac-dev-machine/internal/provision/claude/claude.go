package claude

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run sets up Claude configuration
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
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
		out.Skipped("CLAUDE.md already exists")
		return nil
	}

	// Get templates directory
	templatesDir, err := util.GetTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	// Copy template
	out.Step("Creating CLAUDE.md from template")
	templatePath := filepath.Join(templatesDir, "team-claude.md")
	if err := copyFile(templatePath, claudeMdPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Set permissions
	if err := os.Chmod(claudeMdPath, 0664); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
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
