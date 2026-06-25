package devboxupdate

import (
	"os"

	"github.com/modell-aachen/machine/internal/devbox"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// RestartGuardEnv guards against an infinite re-exec loop: the executor sets it
// on the process it re-execs into after an update, so the second pass skips the
// update instead of triggering another restart.
const RestartGuardEnv = "MACHINE_SELF_UPDATED"

// AlreadyUpdated reports whether this process is the post-re-exec pass.
func AlreadyUpdated() bool {
	return os.Getenv(RestartGuardEnv) == "1"
}

// Run updates the devbox global environment. Restarting into a newer machine
// binary (so modules added or changed by the update run in the same provision)
// is the executor's job — see provision.Execute. On the post-re-exec pass this
// skips, so a single `machine provision` never updates twice.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	if AlreadyUpdated() {
		out.Skipped("devbox already updated earlier this run")
		return nil
	}

	return devbox.GlobalUpdate(out)
}
