package devbox

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/modell-aachen/machine/internal/output"
)

// GlobalUpdate runs `devbox global update` and refreshes the current process's
// devbox env so subsequent commands see any packages added by the update.
func GlobalUpdate(out *output.Context) error {
	out.Step("Updating devbox global environment")
	if err := out.RunCommand("devbox", "global", "update"); err != nil {
		return fmt.Errorf("failed to update devbox: %w", err)
	}
	return refreshGlobalEnv(out)
}

// refreshGlobalEnv re-sources `devbox global shellenv` into the current process.
// After `devbox global update` adds or removes packages, the inherited PATH and
// other devbox env vars are stale — subsequent commands launched by this process
// would fail to find newly installed binaries without this refresh.
func refreshGlobalEnv(out *output.Context) error {
	out.Step("Refreshing devbox global environment")

	cmd := exec.Command("bash", "-c", `eval "$(devbox global shellenv --preserve-path-stack -r)" && env`)
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
