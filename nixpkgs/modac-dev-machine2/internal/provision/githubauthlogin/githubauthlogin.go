package githubauthlogin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine2/internal/platform"
)

// Run authenticates GitHub CLI
func Run(plat platform.Platform) error {
	_ = plat
	// Check if already authenticated
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err == nil {
		fmt.Println("GitHub CLI is already authenticated.")
		return nil
	}

	// Not authenticated, log in
	fmt.Println("Logging into GitHub CLI...")
	loginCmd := exec.Command("gh", "auth", "login", "--web")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	loginCmd.Stdin = os.Stdin
	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate GitHub CLI: %w", err)
	}

	return nil
}
