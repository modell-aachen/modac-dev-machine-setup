package pnpm

import (
	"fmt"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run sets up Node.js tooling (pnpm via corepack)
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Get devbox global path
	devboxPath, err := util.GetDevboxGlobalPath()
	if err != nil {
		return fmt.Errorf("failed to get devbox global path: %w", err)
	}

	// Check if pnpm is already installed
	pnpmPath := filepath.Join(devboxPath, ".devbox", "virtenv", "nodejs", "corepack-bin", "pnpm")
	if util.FileExists(pnpmPath) {
		out.Skipped("Pnpm is already installed")
		return nil
	}

	// Install pnpm globally via corepack
	out.Step("Installing pnpm via corepack")
	if err := out.RunCommand("corepack", "install", "-g", "pnpm@latest"); err != nil {
		return fmt.Errorf("failed to install pnpm: %w", err)
	}

	return nil
}
