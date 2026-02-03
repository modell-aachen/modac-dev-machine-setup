package installmodacshellhelper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Run installs or upgrades modac-shell-helper
func Run() error {
	// Check if modac-shell-helper is installed
	installed, err := isInstalled()
	if err != nil {
		return fmt.Errorf("failed to check if modac-shell-helper is installed: %w", err)
	}

	if !installed {
		// Check if gh auth token exists
		hasToken, err := hasGhAuthToken()
		if err != nil {
			return fmt.Errorf("failed to check gh auth token: %w", err)
		}

		if !hasToken {
			return fmt.Errorf("you need to setup gh auth login")
		}

		// Setup git authentication
		fmt.Println("Setting up git authentication...")
		setupCmd := exec.Command("gh", "auth", "setup-git")
		setupCmd.Stdout = os.Stdout
		setupCmd.Stderr = os.Stderr
		if err := setupCmd.Run(); err != nil {
			return fmt.Errorf("failed to setup git authentication: %w", err)
		}

		// Install modac-shell-helper
		fmt.Println("Installing modac-shell-helper...")
		installCmd := exec.Command("uv", "tool", "install",
			"git+https://github.com/modell-aachen/modac-shell-helper.git@main",
			"--force")
		installCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install modac-shell-helper: %w", err)
		}
	} else {
		// Upgrade modac-shell-helper
		fmt.Println("Upgrading modac-shell-helper...")
		upgradeCmd := exec.Command("uv", "tool", "upgrade", "modac-shell-helper")
		upgradeCmd.Stdout = os.Stdout
		upgradeCmd.Stderr = os.Stderr
		if err := upgradeCmd.Run(); err != nil {
			return fmt.Errorf("failed to upgrade modac-shell-helper: %w", err)
		}
	}

	return nil
}

func isInstalled() (bool, error) {
	cmd := exec.Command("uv", "tool", "list")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return bytes.Contains(output, []byte("modac-shell-helper")), nil
}

func hasGhAuthToken() (bool, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return false, nil // No token if command fails
	}

	return len(bytes.TrimSpace(output)) > 0, nil
}
