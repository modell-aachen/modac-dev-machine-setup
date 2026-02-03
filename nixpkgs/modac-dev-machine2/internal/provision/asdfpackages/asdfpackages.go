package asdfpackages

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine2/internal/platform"
)

// Run installs packages required for asdf version manager
func Run(plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		// No packages needed on Darwin
		return nil
	case platform.Ubuntu:
		return runUbuntu()
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runUbuntu() error {
	fmt.Println("Installing asdf dependencies...")

	asdfDeps := []string{
		"automake",
		"autoconf",
		"libncurses5-dev",
		"libssl-dev",
	}

	args := append([]string{"apt", "install", "-y"}, asdfDeps...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install asdf dependencies: %w", err)
	}

	return nil
}
