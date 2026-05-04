package kubectlkrew

import (
	"fmt"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run installs kubectl krew plugins
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	plugins := []string{"konfig", "oidc-login"}

	// Remove plugins that have been replaced by other tools (e.g. ctx/ns -> kubie).
	// Errors are tolerated because the plugin may not be installed on fresh machines.
	obsolete := []string{"ctx", "ns"}
	for _, plugin := range obsolete {
		out.Step(fmt.Sprintf("Removing obsolete krew plugin: %s", plugin))
		_ = out.RunCommand("krew", "uninstall", plugin)
	}

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
