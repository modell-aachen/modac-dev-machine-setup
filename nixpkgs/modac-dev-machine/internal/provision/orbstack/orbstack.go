package orbstack

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

const (
	statusRunning = "Running"
	// startAttempts and pollInterval bound how long we wait for OrbStack to
	// finish booting after a start. orbctl start can return before the engine
	// is ready, so we poll the status instead of trusting the start exit code.
	startAttempts = 30
	pollInterval  = 2 * time.Second
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
	if orbStatus() == statusRunning {
		out.Skipped("OrbStack is already running")
		return nil
	}

	// The status is "Stopped" or could not be determined (orbctl returns an
	// empty status in some cold-start states). Both cases lead to the same
	// action: attempt to start OrbStack. orbctl start may exit non-zero while
	// the engine is still booting, so its error is logged but not treated as
	// fatal; the post-start status is authoritative.
	out.Step("Starting OrbStack")
	if err := out.RunCommand("orbctl", "start"); err != nil {
		out.Info(fmt.Sprintf("orbctl start exited with %v, verifying status", err))
	}

	sleep := func() { time.Sleep(pollInterval) }
	if err := waitForRunning(orbStatus, sleep, startAttempts); err != nil {
		return err
	}

	out.Step("Logging into OrbStack")
	if err := out.RunCommand("orbctl", "login"); err != nil {
		return fmt.Errorf("failed to login to OrbStack: %w", err)
	}

	return nil
}

// orbStatus returns the trimmed output of `orbctl status`. orbctl exits
// non-zero when OrbStack is not running, so the status string is authoritative
// rather than the exit code; an empty result means the status is unknown.
func orbStatus() string {
	out, _ := exec.Command("orbctl", "status").Output()
	return string(bytes.TrimSpace(out))
}

// waitForRunning polls status until OrbStack reports running, waiting between
// checks. It checks attempts+1 times in total and returns an error only if the
// engine never reaches the running state.
func waitForRunning(status func() string, wait func(), attempts int) error {
	for i := 0; ; i++ {
		if status() == statusRunning {
			return nil
		}
		if i >= attempts {
			return fmt.Errorf("OrbStack did not reach %q state after start", statusRunning)
		}
		wait()
	}
}
