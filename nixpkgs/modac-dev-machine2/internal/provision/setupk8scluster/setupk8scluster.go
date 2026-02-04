package setupk8scluster

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/modell-aachen/machine2/internal/output"
	"github.com/modell-aachen/machine2/internal/platform"
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
	cmd := exec.Command("op", "signin")
	output, err := cmd.CombinedOutput()

	// If the output contains "[ERROR]", user is not signed in
	if strings.Contains(string(output), "[ERROR]") {
		return fmt.Errorf("you are not logged in to 1Password CLI. Please log into 1Password CLI and try again")
	}

	// If there's an error and it's not about being signed in, return it
	if err != nil && !strings.Contains(string(output), "signed in") {
		return fmt.Errorf("failed to check 1Password CLI sign-in status: %w", err)
	}

	return nil
}
