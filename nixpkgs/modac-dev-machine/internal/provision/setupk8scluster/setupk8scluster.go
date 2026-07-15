package setupk8scluster

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
)

// Run sets up Kubernetes cluster configuration
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Check if signed into 1Password CLI
	if err := checkOpSignIn(); err != nil {
		return err
	}

	// Download and execute kubeconfig setup script
	scriptURL := "https://modell-aachen.github.io/k8s-kubeconfig-setup/kubeconfig-setup.sh"
	out.Step("Downloading kubeconfig setup script")

	resp, err := http.Get(scriptURL)
	if err != nil {
		return fmt.Errorf("failed to download kubeconfig setup script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download script: HTTP %d", resp.StatusCode)
	}

	scriptContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read script content: %w", err)
	}

	// Execute script with --merge flag
	out.Step("Executing kubeconfig setup script")
	cmd := exec.Command("bash", "-s", "--", "--merge")
	cmd.Stdin = bytes.NewReader(scriptContent)
	cmd.Stdout = out.MultiWriter(nil)
	cmd.Stderr = out.MultiWriter(nil)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute kubeconfig setup script: %w", err)
	}

	return nil
}

func checkOpSignIn() error {
	// whoami works with both the desktop app integration and an exported
	// OP_SESSION token; op signin would prompt for a password without stdin
	if err := exec.Command("op", "whoami").Run(); err != nil {
		return fmt.Errorf("you are not logged in to 1Password CLI. Please log in (eval $(op signin)) and try again")
	}
	return nil
}
