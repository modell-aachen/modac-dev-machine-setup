package installmodacshellhelper

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run installs or upgrades modac-shell-helper
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
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
		out.Step("Setting up git authentication")
		if err := out.RunCommand("gh", "auth", "setup-git"); err != nil {
			return fmt.Errorf("failed to setup git authentication: %w", err)
		}

		// Install modac-shell-helper
		out.Step("Installing modac-shell-helper")
		installCmd := exec.Command("uv", "tool", "install",
			"git+https://github.com/modell-aachen/modac-shell-helper.git@main",
			"--force")
		installCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if err := out.RunCommand(installCmd.Path, installCmd.Args[1:]...); err != nil {
			return fmt.Errorf("failed to install modac-shell-helper: %w", err)
		}
	} else {
		// Upgrade modac-shell-helper
		out.Step("Upgrading modac-shell-helper")
		if err := out.RunCommand("uv", "tool", "upgrade", "modac-shell-helper"); err != nil {
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
