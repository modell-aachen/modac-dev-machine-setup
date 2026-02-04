package completions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/platform"
	"github.com/modell-aachen/machine2/internal/util"
)

type completion struct {
	cmd     string
	version string
}

// Run installs shell completions for various tools
func Run(plat platform.Platform) error {
	_ = plat
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	shells := []string{"bash", "zsh"}
	completions := []completion{
		{"flux", "2_5_1"},
		{"op", "2_30_3"},
		{"helm", "3_17_3"},
		{"kubectl", "1_32_3"},
	}

	for _, shell := range shells {
		for _, comp := range completions {
			if err := installCompletion(homeDir, shell, comp.cmd, comp.version); err != nil {
				// Log error but continue with other completions
				fmt.Printf("Warning: failed to install %s completion for %s: %v\n", comp.cmd, shell, err)
			}
		}
	}

	return nil
}

func installCompletion(homeDir, shell, cmd, version string) error {
	shellPath := filepath.Join(homeDir, fmt.Sprintf(".%src", shell))
	completionsPath := filepath.Join(homeDir, fmt.Sprintf(".%s_completions", shell))
	cmdCompletionPath := filepath.Join(completionsPath, fmt.Sprintf("%s_%s.sh", cmd, version))

	// Check if shell rc file exists and completion doesn't exist yet
	if !util.FileExists(shellPath) {
		return nil // Skip if shell rc doesn't exist
	}

	if util.FileExists(cmdCompletionPath) {
		return nil // Already installed
	}

	// Create completions directory
	if err := os.MkdirAll(completionsPath, 0755); err != nil {
		return fmt.Errorf("failed to create completions directory: %w", err)
	}

	// Remove old completion files for this command
	pattern := filepath.Join(completionsPath, fmt.Sprintf("%s*.sh", cmd))
	matches, err := filepath.Glob(pattern)
	if err == nil {
		for _, match := range matches {
			os.Remove(match)
		}
	}

	// Generate completion
	fmt.Printf("Installing %s completion for %s under %s\n", cmd, shell, cmdCompletionPath)
	completionCmd := exec.Command(cmd, "completion", shell)
	output, err := completionCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate completion: %w", err)
	}

	// Write completion to file
	if err := os.WriteFile(cmdCompletionPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write completion file: %w", err)
	}

	return nil
}
