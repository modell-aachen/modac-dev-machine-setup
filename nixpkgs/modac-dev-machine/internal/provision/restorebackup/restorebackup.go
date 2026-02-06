package restorebackup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/backup"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run restores backup from 1Password if this is the first provisioning
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	markerFile := filepath.Join(homeDir, ".machine", "restore-backup-done")

	// Only run this once during first provisioning
	return out.CheckAndRun(markerFile, "Backup already restored", func() error {
		out.Step("Checking for existing backup in 1Password (first pass)")

		// First restore pass
		if err := backup.Restore(""); err != nil {
			// If error is about 1Password not being installed or not signed in,
			// skip gracefully since onepassword module should have handled this
			out.Skipped(fmt.Sprintf("Could not restore backup (first pass): %v", err))
		} else {
			out.Success("Backup restored successfully (first pass)")
		}

		out.Step("Restoring backup (second pass)")

		// Second restore pass
		if err := backup.Restore(""); err != nil {
			out.Skipped(fmt.Sprintf("Could not restore backup (second pass): %v", err))
		} else {
			out.Success("Backup restored successfully (second pass)")
		}

		// Run devbox global update
		out.Step("Updating devbox global environment")
		if err := out.RunCommand("devbox", "global", "update"); err != nil {
			return fmt.Errorf("failed to update devbox: %w", err)
		}

		return nil
	})
}
