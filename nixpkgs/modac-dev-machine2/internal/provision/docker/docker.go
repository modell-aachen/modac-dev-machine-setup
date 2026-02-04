package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/util"
)

// Run sets up Docker buildx
func Run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	markerFile := filepath.Join(homeDir, ".docker_buildx_builder_created")

	// Check if buildx builder already created
	if util.FileExists(markerFile) {
		return nil
	}

	// Create buildx builder
	fmt.Println("Creating docker buildx builder...")
	cmd := exec.Command("docker", "buildx", "create", "--use")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create buildx builder: %w", err)
	}

	// Create marker file
	if err := os.WriteFile(markerFile, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create marker file: %w", err)
	}

	return nil
}
