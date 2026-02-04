package docker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run sets up Docker buildx
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	markerFile := filepath.Join(homeDir, ".docker_buildx_builder_created")

	// Use CheckAndRun to handle marker file logic
	return out.CheckAndRun(markerFile, "Docker buildx builder already created", func() error {
		out.Step("Creating docker buildx builder")
		if err := out.RunCommand("docker", "buildx", "create", "--use"); err != nil {
			return fmt.Errorf("failed to create buildx builder: %w", err)
		}
		return nil
	})
}
