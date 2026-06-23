package gcloudworkforcelogin

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// ariadne-gcp-workforce-login is shipped with the binary and installed into the
// user's PATH. It mints fresh Azure AD id_tokens from a stored refresh token so
// GKE sees current group memberships on every token refresh.
//
//go:embed ariadne-gcp-workforce-login
var execScript []byte

// Workforce Identity Federation, federated via Azure AD, backs GKE cluster
// authentication (Identity Service for GKE is discontinued on 2026-07-01).
// Instead of a browser sign-in that freezes the group set until the next login,
// we configure an executable-sourced external_account credential: gke-gcloud-
// auth-plugin obtains a token from gcloud, gcloud re-runs the executable on
// every STS refresh, and the executable re-mints an Azure id_token whose
// `groups` claim reflects current (incl. JIT-activated) memberships.
const (
	workforcePool     = "ariadne-workforce"
	workforceProvider = "ariadne-workforce"
	// quotaProject attributes API quota/billing for the org-wide workforce
	// identity's calls; the principal needs serviceusage.services.use on it.
	quotaProject = "modac-dev"
	execName     = "ariadne-gcp-workforce-login"
	// allowExecutables is gcloud's security gate for executable-sourced creds.
	allowExecutables = "GOOGLE_EXTERNAL_ACCOUNT_ALLOW_EXECUTABLES=1"
)

// Run installs the credential helper, writes the credential config, ensures a
// valid Azure refresh token (one-time interactive device-code login), and
// activates the credential so gke-gcloud-auth-plugin can obtain GKE tokens.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}
	execPath := filepath.Join(home, ".local", "bin", execName)
	credConfig := filepath.Join(home, ".config", "gcloud", "ariadne-cred-config.json")

	if err := installExecutable(out, execPath); err != nil {
		return err
	}
	if err := writeCredConfig(out, credConfig, execPath); err != nil {
		return err
	}

	if refreshTokenValid(execPath) {
		out.Skipped("Azure refresh token already valid")
	} else {
		out.Step("Logging into Azure AD via device code (interactive, one-time)")
		if err := deviceLogin(execPath); err != nil {
			return fmt.Errorf("failed to obtain Azure refresh token: %w", err)
		}
	}

	out.Step("Activating the Workforce Identity credential for gcloud")
	if err := activateCredential(out, credConfig); err != nil {
		return fmt.Errorf("failed to activate workforce credential: %w", err)
	}
	if err := verifyCredential(); err != nil {
		return fmt.Errorf("workforce credential is not usable: %w", err)
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

func writeCredConfig(out *output.Context, path, execPath string) error {
	out.Step("Writing Workforce Identity credential config")
	data, err := buildCredConfig(execPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create gcloud config directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write credential config: %w", err)
	}
	return nil
}

// buildCredConfig renders the external_account credential config that points
// gcloud at the credential helper.
func buildCredConfig(execPath string) ([]byte, error) {
	cfg := map[string]any{
		"type": "external_account",
		"audience": fmt.Sprintf(
			"//iam.googleapis.com/locations/global/workforcePools/%s/providers/%s",
			workforcePool, workforceProvider),
		"subject_token_type":          "urn:ietf:params:oauth:token-type:id_token",
		"token_url":                   "https://sts.googleapis.com/v1/token",
		"workforce_pool_user_project": quotaProject,
		"credential_source": map[string]any{
			"executable": map[string]any{
				"command":        execPath,
				"timeout_millis": 10000,
			},
		},
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credential config: %w", err)
	}
	return append(data, '\n'), nil
}

// refreshTokenValid reports whether the helper can already mint a token, i.e.
// a usable refresh token is present and no interactive login is required.
func refreshTokenValid(execPath string) bool {
	output, err := exec.Command(execPath).Output()
	if err != nil {
		return false
	}
	return bytes.Contains(output, []byte(`"success":true`))
}

func deviceLogin(execPath string) error {
	cmd := exec.Command(execPath, "--login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func activateCredential(out *output.Context, credConfig string) error {
	cmd := exec.Command("gcloud", "auth", "login", "--cred-file="+credConfig, "--quiet")
	cmd.Env = append(os.Environ(), allowExecutables)
	cmd.Stdout = out.MultiWriter(nil)
	cmd.Stderr = out.MultiWriter(nil)
	return cmd.Run()
}

func verifyCredential() error {
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	cmd.Env = append(os.Environ(), allowExecutables)
	return cmd.Run()
}
