package onepassword

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run installs 1Password and 1Password CLI based on the platform
func Run(out *output.Context, plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin(out)
	case platform.Ubuntu:
		return runUbuntu(out)
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin(out *output.Context) error {
	// Check if 1password-cli is already installed
	cmd := exec.Command("brew", "list", "1password-cli")
	if err := cmd.Run(); err == nil {
		out.Skipped("1Password already installed")
		return nil
	}

	out.Step("Installing 1Password and CLI via Homebrew")
	if err := out.RunCommand("brew", "install", "--cask", "1password", "1password-cli"); err != nil {
		return fmt.Errorf("failed to install 1Password: %w", err)
	}

	return nil
}

func runUbuntu(out *output.Context) error {
	// Check if 1password is already installed
	cmd := exec.Command("dpkg", "-l", "1password")
	output, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() == 0 && len(output) > 0 {
		// Check if package is actually installed (starts with 'ii')
		if len(output) > 2 && output[0] == 'i' && output[1] == 'i' {
			out.Skipped("1Password already installed")
			return exportToDistroboxIfNeeded(out)
		}
	}

	opKeyring := "/usr/share/keyrings/1password-archive-keyring.gpg"
	if _, err := os.Stat(opKeyring); os.IsNotExist(err) {
		out.Step("Adding 1Password GPG keyring")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/keys/1password.asc | sudo gpg --dearmor --output "+opKeyring); err != nil {
			return fmt.Errorf("failed to add 1Password keyring: %w", err)
		}
	}

	opSourceList := "/etc/apt/sources.list.d/1password.list"
	if _, err := os.Stat(opSourceList); os.IsNotExist(err) {
		out.Step("Adding 1Password apt repository")
		if err := out.RunCommand("bash", "-c", fmt.Sprintf("echo 'deb [arch=amd64 signed-by=%s] https://downloads.1password.com/linux/debian/amd64 stable main' | sudo tee %s", opKeyring, opSourceList)); err != nil {
			return fmt.Errorf("failed to add 1Password repository: %w", err)
		}
	}

	opGpg := "/usr/share/debsig/keyrings/AC2D62742012EA22/debsig.gpg"
	if _, err := os.Stat(opGpg); os.IsNotExist(err) {
		out.Step("Setting up 1Password package signing")

		if err := out.RunCommand("sudo", "mkdir", "-p", "/etc/debsig/policies/AC2D62742012EA22"); err != nil {
			return fmt.Errorf("failed to create debsig policy directory: %w", err)
		}

		out.Step("Adding debsig policy")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/debian/debsig/1password.pol | sudo tee /etc/debsig/policies/AC2D62742012EA22/1password.pol"); err != nil {
			return fmt.Errorf("failed to add debsig policy: %w", err)
		}

		if err := out.RunCommand("sudo", "mkdir", "-p", "/usr/share/debsig/keyrings/AC2D62742012EA22"); err != nil {
			return fmt.Errorf("failed to create debsig keyring directory: %w", err)
		}

		out.Step("Adding debsig keyring")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/keys/1password.asc | sudo gpg --dearmor --output "+opGpg); err != nil {
			return fmt.Errorf("failed to add debsig keyring: %w", err)
		}

		out.Step("Updating apt package list")
		if err := out.RunCommand("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt: %w", err)
		}
	}

	out.Step("Installing 1Password and CLI")
	if err := out.RunCommand("sudo", "apt", "install", "-y", "1password", "1password-cli"); err != nil {
		return fmt.Errorf("failed to install 1Password: %w", err)
	}

	if isDistrobox() {
		out.Step("Installing audio libraries for Distrobox")
		if err := out.RunCommand("sudo", "apt", "install", "-y", "libasound2t64"); err != nil {
			return fmt.Errorf("failed to install audio libraries: %w", err)
		}
	}

	return exportToDistroboxIfNeeded(out)
}

func isDistrobox() bool {
	return os.Getenv("CONTAINER_ID") != ""
}

func exportToDistroboxIfNeeded(out *output.Context) error {
	if !isDistrobox() {
		return nil
	}

	out.Step("Exporting 1Password app to host")
	if err := out.RunCommand("distrobox-export", "--app", "1password"); err != nil {
		return fmt.Errorf("failed to export 1Password app: %w", err)
	}

	return nil
}
