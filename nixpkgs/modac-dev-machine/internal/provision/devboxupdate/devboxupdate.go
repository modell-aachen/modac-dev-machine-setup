package devboxupdate

import (
	"fmt"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run updates the devbox global environment
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	out.Step("Updating devbox global environment")
	if err := out.RunCommand("devbox", "global", "update"); err != nil {
		return fmt.Errorf("failed to update devbox: %w", err)
	}

	return nil
}
