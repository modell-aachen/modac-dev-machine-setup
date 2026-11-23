package devbox

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modell-aachen/machine/internal/output"
)

// fakeDevboxScript fails or succeeds `devbox global update` and
// `devbox global shellenv` depending on FAKE_UPDATE_EXIT and FAKE_SHELLENV_EXIT.
// It records every shellenv call, which is what reinstalls missing packages of
// the global profile.
const fakeDevboxScript = `#!/bin/sh
case "$2" in
  update) exit "${FAKE_UPDATE_EXIT:-0}" ;;
  shellenv)
    : > "$FAKE_SHELLENV_MARKER"
    [ "${FAKE_SHELLENV_EXIT:-0}" = 0 ] || exit "$FAKE_SHELLENV_EXIT"
    echo 'export FAKE_REFRESHED=1' ;;
esac
`

func setupFakeDevbox(t *testing.T, updateExit, shellenvExit string) (out *output.Context, shellenvMarker string) {
	t.Helper()

	home := t.TempDir()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "devbox"), []byte(fakeDevboxScript), 0755); err != nil {
		t.Fatalf("writing fake devbox: %v", err)
	}

	shellenvMarker = filepath.Join(home, "shellenv-called")
	t.Setenv("HOME", home)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_UPDATE_EXIT", updateExit)
	t.Setenv("FAKE_SHELLENV_EXIT", shellenvExit)
	t.Setenv("FAKE_SHELLENV_MARKER", shellenvMarker)
	t.Setenv("FAKE_REFRESHED", "")

	out, err := output.New()
	if err != nil {
		t.Fatalf("output.New(): %v", err)
	}
	t.Cleanup(func() { _ = out.Close() })

	return out, shellenvMarker
}

func TestGlobalUpdate(t *testing.T) {
	tests := []struct {
		name          string
		updateExit    string
		shellenvExit  string
		wantErrs      []string
		wantRefreshed bool
	}{
		{"successful update refreshes the global environment", "0", "0", nil, true},
		{"failed update still restores the global profile", "1", "0", []string{"failed to update devbox"}, true},
		{"failed update and failed restore report both errors", "1", "1", []string{"failed to update devbox", "failed to load devbox shellenv"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, shellenvMarker := setupFakeDevbox(t, tt.updateExit, tt.shellenvExit)

			err := GlobalUpdate(out)

			if (err != nil) != (len(tt.wantErrs) > 0) {
				t.Fatalf("GlobalUpdate() error = %v, want errors %q", err, tt.wantErrs)
			}
			for _, want := range tt.wantErrs {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("GlobalUpdate() error = %q, missing %q", err, want)
				}
			}
			if _, statErr := os.Stat(shellenvMarker); statErr != nil {
				t.Error("devbox global shellenv was not called, the global profile stays unrepaired")
			}
			if got := os.Getenv("FAKE_REFRESHED") == "1"; got != tt.wantRefreshed {
				t.Errorf("environment refreshed = %v, want %v", got, tt.wantRefreshed)
			}
		})
	}
}
