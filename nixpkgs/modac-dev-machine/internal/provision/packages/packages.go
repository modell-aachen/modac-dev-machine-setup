package packages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run installs system packages based on the platform
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
	brewPackages := []string{
		"bash",
		"gettext",
		"gnu-getopt",
		"gpg",
		"openssh",
		"libfido2",
		"openvpn",
		"nmap",
		"fswatch",
		"gnu-sed",
	}

	brewCasks := []string{
		"visual-studio-code",
		"yubico-authenticator",
		"openvpn-connect",
		"orbstack",
	}

	// Install brew packages
	out.Step("Installing brew packages")
	args := append([]string{"install"}, brewPackages...)
	if err := out.RunCommand("brew", args...); err != nil {
		return fmt.Errorf("failed to install brew packages: %w", err)
	}

	// Install brew casks
	out.Step("Installing brew casks")
	args = append([]string{"install", "--cask"}, brewCasks...)
	if err := out.RunCommand("brew", args...); err != nil {
		return fmt.Errorf("failed to install brew casks: %w", err)
	}

	// Process Brewfile if it exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	brewfilePath := filepath.Join(homeDir, "Brewfile")
	if util.FileExists(brewfilePath) {
		out.Step("Processing Brewfile")

		// Check if bundle needs to be run
		checkCmd := exec.Command("brew", "bundle", "check", "--file="+brewfilePath)
		if err := checkCmd.Run(); err != nil {
			// Bundle check failed, run bundle
			if err := out.RunCommand("brew", "bundle", "--file="+brewfilePath); err != nil {
				return fmt.Errorf("failed to run brew bundle: %w", err)
			}
		} else {
			out.Skipped("Brewfile already satisfied")
		}
	}

	return nil
}

func runUbuntu(out *output.Context) error {
	// Update apt
	out.Step("Updating apt package lists")
	if err := out.RunCommand("sudo", "apt", "update"); err != nil {
		return fmt.Errorf("failed to update apt: %w", err)
	}

	aptPackages := []string{
		"python3-pip",
		"easy-rsa",
		"inotify-tools",
		"net-tools",
		"network-manager-openvpn",
		"network-manager-openvpn-gnome",
		"openvpn",
		"python-is-python3",
		"apt-transport-https",
		"ca-certificates",
		"gnupg",
		"libnss3-tools",
	}

	pythonBuildDeps := []string{
		"make",
		"build-essential",
		"libssl-dev",
		"zlib1g-dev",
		"libbz2-dev",
		"libreadline-dev",
		"libsqlite3-dev",
		"curl",
		"git",
		"libncursesw5-dev",
		"xz-utils",
		"tk-dev",
		"libxml2-dev",
		"libxmlsec1-dev",
		"libffi-dev",
		"liblzma-dev",
		"libzstd-dev",
	}

	// Install base packages
	out.Step("Installing apt packages")
	args := append([]string{"apt", "install", "-y"}, aptPackages...)
	if err := out.RunCommand("sudo", args...); err != nil {
		return fmt.Errorf("failed to install apt packages: %w", err)
	}

	// Install python build dependencies
	out.Step("Installing python build dependencies")
	args = append([]string{"apt", "install", "-y"}, pythonBuildDeps...)
	if err := out.RunCommand("sudo", args...); err != nil {
		return fmt.Errorf("failed to install python build dependencies: %w", err)
	}

	return nil
}
