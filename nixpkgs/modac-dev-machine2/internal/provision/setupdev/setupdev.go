package setupdev

import (
	"fmt"
	"os"

	"github.com/modell-aachen/machine2/internal/output"
	"github.com/modell-aachen/machine2/internal/platform"
)

// Run sets up development directory
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	reposDir := os.Getenv("REPOS_DIRECTORY")
	if reposDir == "" {
		return fmt.Errorf("REPOS_DIRECTORY environment variable not set")
	}

	if !dirExists(reposDir) {
		out.Step(fmt.Sprintf("Creating REPOS_DIRECTORY: %s", reposDir))
		if err := os.MkdirAll(reposDir, 0755); err != nil {
			return fmt.Errorf("failed to create REPOS_DIRECTORY: %w", err)
		}
	} else {
		out.Skipped("REPOS_DIRECTORY already exists")
	}

	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
