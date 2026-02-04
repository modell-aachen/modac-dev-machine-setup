package node

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run sets up Node.js tooling (yarn via corepack)
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Get devbox global path
	devboxPath, err := getDevboxGlobalPath()
	if err != nil {
		return fmt.Errorf("failed to get devbox global path: %w", err)
	}

	// Check if yarn is already installed
	yarnPath := filepath.Join(devboxPath, ".devbox", "virtenv", "nodejs", "corepack-bin", "yarn")
	if util.FileExists(yarnPath) {
		out.Skipped("Yarn is already installed")
		return nil
	}

	// Install yarn globally via corepack
	out.Step("Installing yarn via corepack")
	if err := out.RunCommand("corepack", "install", "-g", "yarn@latest"); err != nil {
		return fmt.Errorf("failed to install yarn: %w", err)
	}

	return nil
}

func getDevboxGlobalPath() (string, error) {
	cmd := exec.Command("devbox", "global", "path")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
