package kubectlkrew

import (
	"fmt"
	"os"
	"os/exec"
)

// Run installs kubectl krew plugins
func Run() error {
	plugins := []string{"ctx", "ns", "konfig", "oidc-login"}

	for _, plugin := range plugins {
		fmt.Printf("Installing krew plugin: %s\n", plugin)
		cmd := exec.Command("krew", "install", plugin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install krew plugin %s: %w", plugin, err)
		}
	}

	fmt.Println("Upgrading krew plugins...")
	cmd := exec.Command("krew", "upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to upgrade krew plugins: %w", err)
	}

	return nil
}
