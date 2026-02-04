package kubectlkrew

import (
	"fmt"

	"github.com/modell-aachen/machine2/internal/output"
	"github.com/modell-aachen/machine2/internal/platform"
)

// Run installs kubectl krew plugins
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	plugins := []string{"ctx", "ns", "konfig", "oidc-login"}

	for _, plugin := range plugins {
		out.Step(fmt.Sprintf("Installing krew plugin: %s", plugin))
		if err := out.RunCommand("krew", "install", plugin); err != nil {
			return fmt.Errorf("failed to install krew plugin %s: %w", plugin, err)
		}
	}

	out.Step("Upgrading krew plugins")
	if err := out.RunCommand("krew", "upgrade"); err != nil {
		return fmt.Errorf("failed to upgrade krew plugins: %w", err)
	}

	return nil
}
