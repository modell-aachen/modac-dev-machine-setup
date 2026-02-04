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

type moduleRunner func(platform.Platform) error

var moduleRegistry = map[string]moduleRunner{
	"packages":                   packages.Run,
	"setup-envs":                 setupenvs.Run,
	"asdf-packages":              asdfpackages.Run,
	"asdf":                       asdf.Run,
	"kubectl-krew":               kubectlkrew.Run,
	"setup-k8s-cluster":          setupk8scluster.Run,
	"node":                       node.Run,
	"certificates":               certificates.Run,
	"setup-dev":                  setupdev.Run,
	"completions":                completions.Run,
	"claude":                     claude.Run,
	"github-auth-login":          githubauthlogin.Run,
	"install-modac-shell-helper": installmodacshellhelper.Run,
	"orbstack":                   orbstack.Run,
	"docker-packages":            dockerpackages.Run,
	"docker":                     docker.Run,
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

	runner, ok := moduleRegistry[module]
	if !ok {
		return fmt.Errorf("unknown module: %s", module)
	}

	return runner(plat)
}
