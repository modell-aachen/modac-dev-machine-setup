package orbstack

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/modell-aachen/machine2/internal/platform"
)

// Run manages OrbStack service
func Run(plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin()
	case platform.Ubuntu:
		fmt.Println("Skipping for OrbStack on Ubuntu")
		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin() error {
	// Check OrbStack status
	statusCmd := exec.Command("orbctl", "status")
	output, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check OrbStack status: %w", err)
	}

	status := string(bytes.TrimSpace(output))
	if status == "Stopped" {
		// Start OrbStack
		fmt.Println("Starting OrbStack...")
		startCmd := exec.Command("orbctl", "start")
		startCmd.Stdout = os.Stdout
		startCmd.Stderr = os.Stderr
		if err := startCmd.Run(); err != nil {
			return fmt.Errorf("failed to start OrbStack: %w", err)
		}

		// Login to OrbStack
		fmt.Println("Login to OrbStack ..")
		loginCmd := exec.Command("orbctl", "login")
		loginCmd.Stdout = os.Stdout
		loginCmd.Stderr = os.Stderr
		if err := loginCmd.Run(); err != nil {
			return fmt.Errorf("failed to login to OrbStack: %w", err)
		}
	}

	return nil
}
