package devbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/modell-aachen/machine/internal/output"
)

// GlobalUpdate runs `devbox global update` and refreshes the current process's
// devbox env so subsequent commands see any packages added by the update.
//
// The refresh also runs when the update fails: devbox removes outdated packages
// from the global profile before installing their replacements, so a failed
// install leaves binaries like machine and claude missing. The shellenv call
// in refreshGlobalEnv reinstalls them.
func GlobalUpdate(out *output.Context) error {
	out.Step("Updating devbox global environment")
	if err := out.RunCommand("devbox", "global", "update"); err != nil {
		updateErr := fmt.Errorf("failed to update devbox: %w", err)
		return errors.Join(updateErr, refreshGlobalEnv(out))
	}
	return refreshGlobalEnv(out)
}

// refreshGlobalEnv re-sources `devbox global shellenv` into the current process.
// After `devbox global update` adds or removes packages, the inherited PATH and
// other devbox env vars are stale — subsequent commands launched by this process
// would fail to find newly installed binaries without this refresh.
func refreshGlobalEnv(out *output.Context) error {
	out.Step("Refreshing devbox global environment")

	// The assignment keeps devbox's exit code, which `eval "$(...)"` would swallow.
	cmd := exec.Command("bash", "-c", `shellenv="$(devbox global shellenv --preserve-path-stack -r)" && eval "$shellenv" && env`)
	cmd.Stderr = out.MultiWriter(nil)
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to load devbox shellenv: %w", err)
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key, value := line[:idx], line[idx+1:]
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set env %s: %w", key, err)
		}
	}

	return nil
}
