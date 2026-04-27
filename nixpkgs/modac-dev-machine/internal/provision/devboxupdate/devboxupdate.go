package devboxupdate

import (
	"github.com/modell-aachen/machine/internal/devbox"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run updates the devbox global environment
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	return devbox.GlobalUpdate(out)
}
