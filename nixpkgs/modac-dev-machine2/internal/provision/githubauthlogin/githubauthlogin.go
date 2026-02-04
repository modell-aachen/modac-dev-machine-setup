package githubauthlogin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine2/internal/output"
	"github.com/modell-aachen/machine2/internal/platform"
)

// Run authenticates GitHub CLI
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Check if already authenticated
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err == nil {
		out.Skipped("GitHub CLI is already authenticated")
		return nil
	}

	// Not authenticated, log in (this requires user interaction)
	out.Step("Logging into GitHub CLI (interactive)")
	loginCmd := exec.Command("gh", "auth", "login", "--web")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	loginCmd.Stdin = os.Stdin
	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate GitHub CLI: %w", err)
	}

	return nil
}
