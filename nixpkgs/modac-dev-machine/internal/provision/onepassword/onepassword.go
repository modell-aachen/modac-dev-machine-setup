package onepassword

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/modell-aachen/machine/internal/config"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run installs 1Password and 1Password CLI based on the platform
func Run(out *output.Context, plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin(out)
	case platform.Ubuntu:
		return runUbuntu(out)
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin(out *output.Context) error {
	// Check if 1password-cli is already installed
	cmd := exec.Command("brew", "list", "1password-cli")
	if err := cmd.Run(); err == nil {
		out.Skipped("1Password already installed")
		return postInstallSetup(out, platform.Darwin)
	}

	out.Step("Installing 1Password and CLI via Homebrew")
	if err := out.RunCommand("brew", "install", "--cask", "1password", "1password-cli"); err != nil {
		return fmt.Errorf("failed to install 1Password: %w", err)
	}

	return postInstallSetup(out, platform.Darwin)
}

func runUbuntu(out *output.Context) error {
	cliOnly, err := useCLIOnly()
	if err != nil {
		return err
	}
	if cliOnly {
		return runUbuntuCLIOnly(out)
	}

	// Check if 1password is already installed
	if isDebPackageInstalled("1password") {
		out.Skipped("1Password already installed")
		if err := exportToDistroboxIfNeeded(out); err != nil {
			return err
		}
		return postInstallSetup(out, platform.Ubuntu)
	}

	if err := ensureAptRepo(out); err != nil {
		return err
	}

	out.Step("Installing 1Password and CLI")
	if err := out.RunCommand("sudo", "apt", "install", "-y", "1password", "1password-cli"); err != nil {
		return fmt.Errorf("failed to install 1Password: %w", err)
	}

	if isDistrobox() {
		out.Step("Installing audio libraries for Distrobox")
		if err := out.RunCommand("sudo", "apt", "install", "-y", "libasound2t64"); err != nil {
			return fmt.Errorf("failed to install audio libraries: %w", err)
		}
	}

	if err := exportToDistroboxIfNeeded(out); err != nil {
		return err
	}

	return postInstallSetup(out, platform.Ubuntu)
}

// useCLIOnly decides whether to skip the desktop app entirely: service
// machines never need it, and under WSL there is no desktop to run it on.
func useCLIOnly() (bool, error) {
	if platform.IsWSL() {
		return true, nil
	}
	profile, err := config.LoadProfile()
	if err != nil {
		return false, err
	}
	return profile == config.ProfileService, nil
}

func runUbuntuCLIOnly(out *output.Context) error {
	if isDebPackageInstalled("1password-cli") {
		out.Skipped("1Password CLI already installed")
		return postInstallSetupCLIOnly(out)
	}

	if err := ensureAptRepo(out); err != nil {
		return err
	}

	out.Step("Updating apt package list")
	if err := out.RunCommand("sudo", "apt", "update"); err != nil {
		return fmt.Errorf("failed to update apt: %w", err)
	}

	out.Step("Installing 1Password CLI")
	if err := out.RunCommand("sudo", "apt", "install", "-y", "1password-cli"); err != nil {
		return fmt.Errorf("failed to install 1Password CLI: %w", err)
	}

	return postInstallSetupCLIOnly(out)
}

func isDebPackageInstalled(name string) bool {
	cmd := exec.Command("dpkg", "-l", name)
	output, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() != 0 {
		return false
	}
	// Package is actually installed when the status line starts with 'ii'
	return len(output) > 2 && output[0] == 'i' && output[1] == 'i'
}

func ensureAptRepo(out *output.Context) error {
	// Minimal Ubuntu images (e.g. fresh WSL) ship without gnupg
	if _, err := exec.LookPath("gpg"); err != nil {
		out.Step("Installing gnupg")
		if err := out.RunCommand("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt: %w", err)
		}
		if err := out.RunCommand("sudo", "apt", "install", "-y", "gnupg"); err != nil {
			return fmt.Errorf("failed to install gnupg: %w", err)
		}
	}

	opKeyring := "/usr/share/keyrings/1password-archive-keyring.gpg"
	if _, err := os.Stat(opKeyring); os.IsNotExist(err) {
		out.Step("Adding 1Password GPG keyring")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/keys/1password.asc | sudo gpg --dearmor --output "+opKeyring); err != nil {
			return fmt.Errorf("failed to add 1Password keyring: %w", err)
		}
	}

	opSourceList := "/etc/apt/sources.list.d/1password.list"
	if _, err := os.Stat(opSourceList); os.IsNotExist(err) {
		out.Step("Adding 1Password apt repository")
		if err := out.RunCommand("bash", "-c", fmt.Sprintf("echo 'deb [arch=amd64 signed-by=%s] https://downloads.1password.com/linux/debian/amd64 stable main' | sudo tee %s", opKeyring, opSourceList)); err != nil {
			return fmt.Errorf("failed to add 1Password repository: %w", err)
		}
	}

	opGpg := "/usr/share/debsig/keyrings/AC2D62742012EA22/debsig.gpg"
	if _, err := os.Stat(opGpg); os.IsNotExist(err) {
		out.Step("Setting up 1Password package signing")

		if err := out.RunCommand("sudo", "mkdir", "-p", "/etc/debsig/policies/AC2D62742012EA22"); err != nil {
			return fmt.Errorf("failed to create debsig policy directory: %w", err)
		}

		out.Step("Adding debsig policy")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/debian/debsig/1password.pol | sudo tee /etc/debsig/policies/AC2D62742012EA22/1password.pol"); err != nil {
			return fmt.Errorf("failed to add debsig policy: %w", err)
		}

		if err := out.RunCommand("sudo", "mkdir", "-p", "/usr/share/debsig/keyrings/AC2D62742012EA22"); err != nil {
			return fmt.Errorf("failed to create debsig keyring directory: %w", err)
		}

		out.Step("Adding debsig keyring")
		if err := out.RunCommand("bash", "-c", "curl -sS https://downloads.1password.com/linux/keys/1password.asc | sudo gpg --dearmor --output "+opGpg); err != nil {
			return fmt.Errorf("failed to add debsig keyring: %w", err)
		}

		out.Step("Updating apt package list")
		if err := out.RunCommand("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt: %w", err)
		}
	}

	return nil
}

func isDistrobox() bool {
	return os.Getenv("CONTAINER_ID") != ""
}

func exportToDistroboxIfNeeded(out *output.Context) error {
	if !isDistrobox() {
		return nil
	}

	out.Step("Exporting 1Password app to host")
	if err := out.RunCommand("distrobox-export", "--app", "1password"); err != nil {
		return fmt.Errorf("failed to export 1Password app: %w", err)
	}

	return nil
}

func postInstallSetup(out *output.Context, plat platform.Platform) error {
	// Step 1: Open the 1Password app
	if err := openOnePasswordApp(out, plat); err != nil {
		return err
	}

	// Step 2: Ensure CLI integration is enabled
	if err := ensureCLIIntegration(out); err != nil {
		return err
	}

	// Step 3: Sign in the user
	if err := signInUser(out); err != nil {
		return err
	}

	return nil
}

func openOnePasswordApp(out *output.Context, plat platform.Platform) error {
	out.Step("Opening 1Password app")

	var cmd *exec.Cmd
	switch plat {
	case platform.Darwin:
		cmd = exec.Command("open", "-a", "1Password", "--args", "--silent")
	case platform.Ubuntu:
		cmd = exec.Command("1password", "--silent")
	default:
		return fmt.Errorf("unsupported platform for opening 1Password: %s", plat)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open 1Password app: %w", err)
	}

	out.Success("1Password app opened")
	out.Step("Waiting for 1Password app to be ready")

	// Poll until the app is ready (max 60 seconds)
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for 1Password app to be ready. Please ensure the app is unlocked and try again")
		case <-ticker.C:
			// Check if the CLI can connect to the app
			checkCmd := exec.Command("op", "account", "list")
			if err := checkCmd.Run(); err == nil {
				out.Success("1Password app is ready")
				return nil
			}
			// If error contains "cannot connect", app is still starting up - keep polling
		}
	}
}

func ensureCLIIntegration(out *output.Context) error {
	out.Step("Checking 1Password CLI integration")

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd := exec.Command("op", "--format", "json", "account", "list")
		output, err := cmd.CombinedOutput()

		if err == nil && len(output) > 0 {
			// Try to parse as JSON array to confirm it's valid
			var accounts []interface{}
			if err := json.Unmarshal(output, &accounts); err == nil && len(accounts) > 0 {
				out.Success("1Password CLI integration is enabled")
				return nil
			}
		}

		// CLI integration is not enabled
		fmt.Println("\n⚠️  1Password CLI integration is not enabled.")
		fmt.Println("To enable CLI integration:")
		fmt.Println("  1. Open the 1Password app")
		fmt.Println("  2. Go to Settings → Developer")
		fmt.Println("  3. Enable 'Connect with 1Password CLI'")
		fmt.Println("\nPress Enter once you've enabled CLI integration to continue...")

		_, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}
	}
}

// postInstallSetupCLIOnly signs in without the desktop app. The session token
// is exported into this process so later modules (setup-k8s-cluster) can use
// op within the same provision run. A 1Password service account via
// OP_SERVICE_ACCOUNT_TOKEN would be the non-interactive alternative.
func postInstallSetupCLIOnly(out *output.Context) error {
	if !hasConfiguredAccount() {
		out.Step("Adding 1Password account")
		cmd := exec.Command("op", "account", "add")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add 1Password account: %w", err)
		}
	}

	out.Step("Signing in to 1Password CLI")
	var stdout bytes.Buffer
	cmd := exec.Command("op", "signin")
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to sign in to 1Password CLI: %w", err)
	}

	for name, token := range parseSessionExports(stdout.String()) {
		if err := os.Setenv(name, token); err != nil {
			return fmt.Errorf("failed to export 1Password session: %w", err)
		}
	}

	out.Success("Signed in to 1Password CLI")
	return nil
}

func hasConfiguredAccount() bool {
	output, err := exec.Command("op", "--format", "json", "account", "list").Output()
	if err != nil {
		return false
	}
	var accounts []interface{}
	return json.Unmarshal(output, &accounts) == nil && len(accounts) > 0
}

var sessionExportPattern = regexp.MustCompile(`(?m)^export (OP_SESSION_[A-Za-z0-9_]+)="([^"]+)"$`)

func parseSessionExports(signinOutput string) map[string]string {
	sessions := map[string]string{}
	for _, match := range sessionExportPattern.FindAllStringSubmatch(signinOutput, -1) {
		sessions[match[1]] = match[2]
	}
	return sessions
}

func signInUser(out *output.Context) error {
	out.Step("Signing in to 1Password")

	cmd := exec.Command("op", "signin")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Check if already signed in
		checkCmd := exec.Command("op", "account", "get")
		if checkErr := checkCmd.Run(); checkErr == nil {
			out.Success("Already signed in to 1Password")
			return nil
		}
		return fmt.Errorf("failed to sign in to 1Password: %w", err)
	}

	out.Success("Successfully signed in to 1Password")
	return nil
}
