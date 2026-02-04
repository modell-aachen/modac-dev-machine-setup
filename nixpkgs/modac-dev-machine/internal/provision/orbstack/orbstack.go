package orbstack

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run manages OrbStack service
func Run(out *output.Context, plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin(out)
	case platform.Ubuntu:
		out.Skipped("OrbStack not needed on Ubuntu")
		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin(out *output.Context) error {
	// Check OrbStack status
	statusCmd := exec.Command("orbctl", "status")
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check OrbStack status: %w", err)
	}

	status := string(bytes.TrimSpace(statusOutput))
	if status == "Stopped" {
		// Start OrbStack
		out.Step("Starting OrbStack")
		if err := out.RunCommand("orbctl", "start"); err != nil {
			return fmt.Errorf("failed to start OrbStack: %w", err)
		}

		// Login to OrbStack
		out.Step("Logging into OrbStack")
		if err := out.RunCommand("orbctl", "login"); err != nil {
			return fmt.Errorf("failed to login to OrbStack: %w", err)
		}
	} else {
		out.Skipped("OrbStack is already running")
	}

	return nil
}
