package gcloudworkforcelogin

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// ariadne-gcp-workforce-login is shipped with the binary and installed into the
// user's PATH. It is a client-go exec credential plugin: it mints a fresh Azure
// AD id_token from a stored refresh token, exchanges it at Google STS, and
// prints an ExecCredential for kubectl. Because it talks to STS directly, gcloud
// is not on the runtime path and no GOOGLE_EXTERNAL_ACCOUNT_ALLOW_EXECUTABLES
// gate is needed anywhere. GKE re-evaluates the user's group memberships
// whenever the cached token is re-minted.
//
// The kubeconfig that wires GKE users to this helper is managed by the team in
// 1Password, so this module only installs the helper and performs the one-time
// interactive device-code login. Workforce Identity Federation replaces Identity
// Service for GKE (discontinued 2026-07-01).
//
//go:embed ariadne-gcp-workforce-login
var execScript []byte

const execName = "ariadne-gcp-workforce-login"

// Run installs the credential helper and ensures a valid Azure refresh token,
// signing the user in interactively (device code) if necessary.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}
	execPath := filepath.Join(home, ".local", "bin", execName)

	if err := installExecutable(out, execPath); err != nil {
		return err
	}

	if hasSession(execPath) {
		// A refresh token exists; validate the full path (Azure + STS) so a
		// broken exchange surfaces here instead of on the first kubectl, and is
		// not mistaken for "needs login".
		if err := checkSession(execPath); err != nil {
			return fmt.Errorf("workforce session present but token exchange failed; "+
				"check QUOTA_PROJECT/serviceusage and the provider audience: %w", err)
		}
		out.Skipped("Workforce session already valid")
		return nil
	}

	out.Step("Logging into Azure AD via device code (interactive, one-time)")
	if err := deviceLogin(execPath); err != nil {
		return fmt.Errorf("failed to obtain Azure refresh token: %w", err)
	}
	if err := checkSession(execPath); err != nil {
		return fmt.Errorf("token exchange failed after login; "+
			"check QUOTA_PROJECT/serviceusage and the provider audience: %w", err)
	}

	return nil
}

func installExecutable(out *output.Context, execPath string) error {
	out.Step("Installing " + execName)
	if err := os.MkdirAll(filepath.Dir(execPath), 0o755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}
	if err := os.WriteFile(execPath, execScript, 0o755); err != nil {
		return fmt.Errorf("failed to write %s: %w", execName, err)
	}
	return nil
}

// hasSession reports whether a refresh token is present (no network call).
func hasSession(execPath string) bool {
	return exec.Command(execPath, "--has-session").Run() == nil
}

// checkSession validates the full credential path (Azure refresh + STS
// exchange), surfacing the helper's error message on failure.
func checkSession(execPath string) error {
	out, err := exec.Command(execPath, "--check").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
	}
	return nil
}

func deviceLogin(execPath string) error {
	cmd := exec.Command(execPath, "--login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
