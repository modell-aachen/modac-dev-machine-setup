package provision

import (
	"fmt"

	"github.com/modell-aachen/machine2/internal/platform"
	"github.com/modell-aachen/machine2/internal/provision/asdf"
	"github.com/modell-aachen/machine2/internal/provision/asdfpackages"
	"github.com/modell-aachen/machine2/internal/provision/certificates"
	"github.com/modell-aachen/machine2/internal/provision/claude"
	"github.com/modell-aachen/machine2/internal/provision/completions"
	"github.com/modell-aachen/machine2/internal/provision/docker"
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

	for _, module := range modules {
		if err := runModule(module, plat); err != nil {
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

func runModule(module string, plat platform.Platform) error {
	fmt.Printf("Running %s\n", module)

	// All modules are now implemented in Go
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
	case "docker":
		return docker.Run()
	default:
		return fmt.Errorf("unknown module: %s", module)
	}
}
