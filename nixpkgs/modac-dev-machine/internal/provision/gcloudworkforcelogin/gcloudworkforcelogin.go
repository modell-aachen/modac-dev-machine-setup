package gcloudworkforcelogin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Workforce Identity Federation pool and provider that back GKE cluster
// authentication. Identity Service for GKE is discontinued on 2026-07-01;
// afterwards kubectl authenticates via gke-gcloud-auth-plugin, which obtains a
// token from the active gcloud credential. This module makes that credential a
// Workforce Identity (federated via Azure AD) instead of a Google account.
const (
	workforcePool     = "ariadne-workforce"
	workforceProvider = "ariadne-workforce"
)

// Run configures gcloud for Workforce Identity Federation and signs the user
// in, so that gke-gcloud-auth-plugin can obtain GKE tokens with no further
// manual setup.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}
	loginConfig := filepath.Join(home, ".config", "gcloud", "ariadne-login-config.json")

	// 1. Generate the (non-secret) login config describing the workforce pool
	//    provider, unless it already exists. This must succeed without an active
	//    credential, since it is a prerequisite for signing in.
	if _, statErr := os.Stat(loginConfig); statErr != nil {
		out.Step("Generating Workforce Identity login config")
		provider := fmt.Sprintf("locations/global/workforcePools/%s/providers/%s", workforcePool, workforceProvider)
		genCmd := exec.Command("gcloud", "iam", "workforce-pools", "create-login-config",
			provider, "--output-file="+loginConfig)
		genCmd.Stdout = out.MultiWriter(nil)
		genCmd.Stderr = out.MultiWriter(nil)
		if err := genCmd.Run(); err != nil {
			return fmt.Errorf("failed to create workforce login config: %w", err)
		}
	} else {
		out.Skipped("Workforce login config already present")
	}

	// 2. Make gcloud use the workforce login config for browser sign-in, so a
	//    plain `gcloud auth login` runs the Workforce (Azure AD) flow.
	out.Step("Configuring gcloud to use the Workforce login config")
	setCmd := exec.Command("gcloud", "config", "set", "auth/login_config_file", loginConfig)
	setCmd.Stdout = out.MultiWriter(nil)
	setCmd.Stderr = out.MultiWriter(nil)
	if err := setCmd.Run(); err != nil {
		return fmt.Errorf("failed to set gcloud login config: %w", err)
	}

	// 3. Sign in if there is no active credential yet (interactive browser flow).
	if err := exec.Command("gcloud", "auth", "print-access-token").Run(); err == nil {
		out.Skipped("gcloud already has an active credential")
		return nil
	}

	out.Step("Logging into Google Cloud via Workforce Identity (interactive)")
	loginCmd := exec.Command("gcloud", "auth", "login")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	loginCmd.Stdin = os.Stdin
	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate to Google Cloud: %w", err)
	}

	return nil
}
