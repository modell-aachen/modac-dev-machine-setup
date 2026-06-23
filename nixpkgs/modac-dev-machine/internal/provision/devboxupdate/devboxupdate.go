package devboxupdate

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/modell-aachen/machine/internal/devbox"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// selfUpdatedEnv guards against an infinite re-exec loop: it is set on the
// process we exec into after updating, so the second pass skips the update.
const selfUpdatedEnv = "MACHINE_SELF_UPDATED"

// alreadyUpdated reports whether this process is the post-re-exec pass.
func alreadyUpdated() bool {
	return os.Getenv(selfUpdatedEnv) == "1"
}

// Run updates the devbox global environment and then, on the first pass, re-execs
// the (possibly newer) machine binary so that modules added or changed by the
// update run in the same `machine provision` invocation. The whole provision
// restarts from the top under the new binary; modules are idempotent, so the few
// that run before this one simply run again.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	if alreadyUpdated() {
		out.Skipped("devbox already updated earlier this run")
		return nil
	}

	if err := devbox.GlobalUpdate(out); err != nil {
		return err
	}

	// Restart into whatever machine is now on PATH so the rest of provisioning
	// uses the updated binary. devbox.GlobalUpdate already refreshed PATH into
	// this process's env, which os.Environ() carries to the new image.
	bin, err := exec.LookPath("machine")
	if err != nil {
		// machine not resolvable on PATH; continue with the current binary.
		return nil
	}

	out.Step("Restarting machine with the updated environment")
	if err := syscall.Exec(bin, os.Args, append(os.Environ(), selfUpdatedEnv+"=1")); err != nil {
		return fmt.Errorf("failed to restart machine after update: %w", err)
	}
	return nil // unreachable on success: syscall.Exec replaces the process image
}
