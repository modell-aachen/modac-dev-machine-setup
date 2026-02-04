package node

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine2/internal/util"
)

// Run sets up Node.js tooling (yarn via corepack)
func Run() error {
	// Get devbox global path
	devboxPath, err := getDevboxGlobalPath()
	if err != nil {
		return fmt.Errorf("failed to get devbox global path: %w", err)
	}

	// Check if yarn is already installed
	yarnPath := filepath.Join(devboxPath, ".devbox", "virtenv", "nodejs", "corepack-bin", "yarn")
	if util.FileExists(yarnPath) {
		fmt.Println("Yarn is already installed")
		return nil
	}

	// Install yarn globally via corepack
	fmt.Println("Installing yarn via corepack...")
	cmd := exec.Command("corepack", "install", "-g", "yarn@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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
