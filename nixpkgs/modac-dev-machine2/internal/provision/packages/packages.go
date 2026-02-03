package packages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/platform"
)

// Run installs system packages based on the platform
func Run(plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin()
	case platform.Ubuntu:
		return runUbuntu()
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin() error {
	fmt.Println("Installing brew packages...")

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
	args := append([]string{"install"}, brewPackages...)
	cmd := exec.Command("brew", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install brew packages: %w", err)
	}

	// Install brew casks
	fmt.Println("Installing brew casks...")
	args = append([]string{"install", "--cask"}, brewCasks...)
	cmd = exec.Command("brew", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install brew casks: %w", err)
	}

	// Process Brewfile if it exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	brewfilePath := filepath.Join(homeDir, "Brewfile")
	if fileExists(brewfilePath) {
		fmt.Println("Found Brewfile, processing...")

		// Check if bundle needs to be run
		checkCmd := exec.Command("brew", "bundle", "check", "--file="+brewfilePath)
		if err := checkCmd.Run(); err != nil {
			// Bundle check failed, run bundle
			bundleCmd := exec.Command("brew", "bundle", "--file="+brewfilePath)
			bundleCmd.Stdout = os.Stdout
			bundleCmd.Stderr = os.Stderr
			if err := bundleCmd.Run(); err != nil {
				return fmt.Errorf("failed to run brew bundle: %w", err)
			}
		}
	}

	return nil
}

func runUbuntu() error {
	fmt.Println("Updating apt package lists...")

	// Update apt
	cmd := exec.Command("sudo", "apt", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update apt: %w", err)
	}

	fmt.Println("Installing apt packages...")

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
	args := append([]string{"apt", "install", "-y"}, aptPackages...)
	cmd = exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install apt packages: %w", err)
	}

	// Install python build dependencies
	fmt.Println("Installing python build dependencies...")
	args = append([]string{"apt", "install", "-y"}, pythonBuildDeps...)
	cmd = exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install python build dependencies: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
