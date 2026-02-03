package setupdev

import (
	"fmt"
	"os"
)

// Run sets up development directory
func Run() error {
	reposDir := os.Getenv("REPOS_DIRECTORY")
	if reposDir == "" {
		return fmt.Errorf("REPOS_DIRECTORY environment variable not set")
	}

	if !dirExists(reposDir) {
		fmt.Printf("Creating REPOS_DIRECTORY: %s\n", reposDir)
		if err := os.MkdirAll(reposDir, 0755); err != nil {
			return fmt.Errorf("failed to create REPOS_DIRECTORY: %w", err)
		}
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
