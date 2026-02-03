package provision

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/platform"
	"github.com/modell-aachen/machine2/internal/provision/asdf"
	"github.com/modell-aachen/machine2/internal/provision/asdfpackages"
	"github.com/modell-aachen/machine2/internal/provision/certificates"
	"github.com/modell-aachen/machine2/internal/provision/claude"
	"github.com/modell-aachen/machine2/internal/provision/completions"
	"github.com/modell-aachen/machine2/internal/provision/dockerpackages"
	"github.com/modell-aachen/machine2/internal/provision/githubauthlogin"
	"github.com/modell-aachen/machine2/internal/provision/installmodacshellhelper"
	"github.com/modell-aachen/machine2/internal/provision/kubectlkrew"
	"github.com/modell-aachen/machine2/internal/provision/node"
	"github.com/modell-aachen/machine2/internal/provision/orbstack"
	"github.com/modell-aachen/machine2/internal/provision/packages"
	"github.com/modell-aachen/machine2/internal/provision/setupdev"
	"github.com/modell-aachen/machine2/internal/provision/setupenvs"
	"github.com/modell-aachen/machine2/internal/provision/setupk8scluster"
)

type Options struct {
	Filter      string
	SkipInstall bool
}

func Execute(opts *Options) error {
	plat, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("platform detection failed: %w", err)
	}

	modules := FilterModules(opts.Filter)

	if !opts.SkipInstall {
		fmt.Println("Skipping install step (not implemented in Go CLI)")
		// Note: install command is excluded from this port
	}

	scriptsDir, err := getScriptsDir()
	if err != nil {
		return fmt.Errorf("failed to find scripts directory: %w", err)
	}

	for _, module := range modules {
		if err := runModule(module, plat, scriptsDir); err != nil {
			return fmt.Errorf("module %s failed: %w", module, err)
		}
	}

	return nil
}

func ListModules() error {
	fmt.Println("Available modules:")
	for _, module := range GetAllModules() {
		fmt.Printf("  - %s\n", module)
	}
	return nil
}

func runModule(module string, plat platform.Platform, scriptsDir string) error {
	fmt.Printf("Running %s\n", module)

	// Use Go implementations for ported modules
	switch module {
	case "packages":
		return packages.Run(plat)
	case "setup-envs":
		return setupenvs.Run()
	case "asdf-packages":
		return asdfpackages.Run(plat)
	case "asdf":
		return asdf.Run()
	case "kubectl-krew":
		return kubectlkrew.Run()
	case "setup-k8s-cluster":
		return setupk8scluster.Run()
	case "node":
		return node.Run()
	case "certificates":
		return certificates.Run()
	case "setup-dev":
		return setupdev.Run()
	case "completions":
		return completions.Run()
	case "claude":
		return claude.Run()
	case "github-auth-login":
		return githubauthlogin.Run()
	case "install-modac-shell-helper":
		return installmodacshellhelper.Run()
	case "orbstack":
		return orbstack.Run(plat)
	case "docker-packages":
		return dockerpackages.Run(plat)
	}

	// Fall back to bash script execution for other modules
	scriptPath := findScript(module, plat, scriptsDir)
	if scriptPath == "" {
		return fmt.Errorf("script not found for module: %s", module)
	}

	cmd := exec.Command("bash", scriptPath, scriptsDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", plat),
		"POETRY_VERSION=2.0.1",
		fmt.Sprintf("IS_%s=true", plat),
	)

	return cmd.Run()
}

func findScript(module string, plat platform.Platform, scriptsDir string) string {
	// Try platform-specific first
	platformPath := filepath.Join(scriptsDir, plat.String(), module+".bash")
	if fileExists(platformPath) {
		return platformPath
	}

	// Try shared script
	sharedPath := filepath.Join(scriptsDir, module+".bash")
	if fileExists(sharedPath) {
		return sharedPath
	}

	return ""
}

func getScriptsDir() (string, error) {
	// Get the directory of the currently running binary
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Binary is in bin/, scripts should be in ../share/machine2/provision-scripts/
	binDir := filepath.Dir(exePath)
	shareDir := filepath.Join(binDir, "..", "share", "machine2", "provision-scripts")

	// Check if the share directory exists
	if fileExists(shareDir) {
		return shareDir, nil
	}

	// Fallback: check if we're running from the repo (development mode)
	repoScriptsDir := filepath.Join(binDir, "..", "scripts", "provision")
	if fileExists(repoScriptsDir) {
		return repoScriptsDir, nil
	}

	return "", fmt.Errorf("scripts directory not found")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
