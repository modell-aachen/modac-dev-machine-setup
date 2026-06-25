package devboxupdate

import (
	"os"

	"github.com/modell-aachen/machine/internal/devbox"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

const RestartGuardEnv = "MACHINE_SELF_UPDATED"

func AlreadyUpdated() bool {
	return os.Getenv(RestartGuardEnv) == "1"
}

func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	if AlreadyUpdated() {
		out.Skipped("devbox already updated earlier this run")
		return nil
	}

	return devbox.GlobalUpdate(out)
}
