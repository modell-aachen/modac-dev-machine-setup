package asdfpackages

import (
	"fmt"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run installs packages required for asdf version manager
func Run(out *output.Context, plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		out.Skipped("No packages needed on Darwin")
		return nil
	case platform.Ubuntu:
		return runUbuntu(out)
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runUbuntu(out *output.Context) error {
	asdfDeps := []string{
		"automake",
		"autoconf",
		"libncurses5-dev",
		"libssl-dev",
	}

	out.Step("Installing asdf dependencies")
	args := append([]string{"apt", "install", "-y"}, asdfDeps...)
	if err := out.RunCommand("sudo", args...); err != nil {
		return fmt.Errorf("failed to install asdf dependencies: %w", err)
	}

	return nil
}
