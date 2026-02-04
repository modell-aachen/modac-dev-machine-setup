package provision

import (
	"fmt"
	"strings"

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

// ModuleEntry represents a provisioning module with its name and runner function
type ModuleEntry struct {
	Name   string
	Runner func(platform.Platform) error
}

// allModules defines the ordered list of all provisioning modules
var allModules = []ModuleEntry{
	{"packages", packages.Run},
	{"setup-envs", setupenvs.Run},
	{"asdf-packages", asdfpackages.Run},
	{"asdf", asdf.Run},
	{"kubectl-krew", kubectlkrew.Run},
	{"setup-k8s-cluster", setupk8scluster.Run},
	{"node", node.Run},
	{"certificates", certificates.Run},
	{"setup-dev", setupdev.Run},
	{"completions", completions.Run},
	{"claude", claude.Run},
	{"github-auth-login", githubauthlogin.Run},
	{"install-modac-shell-helper", installmodacshellhelper.Run},
	{"orbstack", orbstack.Run},
	{"docker-packages", dockerpackages.Run},
	{"docker", docker.Run},
}

// GetAllModuleNames returns the names of all available modules
func GetAllModuleNames() []string {
	names := make([]string, len(allModules))
	for i, module := range allModules {
		names[i] = module.Name
	}
	return names
}

// FilterModules returns a filtered list of modules based on the filter string
func FilterModules(filter string) []ModuleEntry {
	if filter == "" {
		return allModules
	}

	filterMap := make(map[string]bool)
	for _, name := range splitCSV(filter) {
		filterMap[name] = true
	}

	filtered := []ModuleEntry{}
	for _, module := range allModules {
		if filterMap[module.Name] {
			filtered = append(filtered, module)
		}
	}

	return filtered
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
			return fmt.Errorf("module %s failed: %w", module.Name, err)
		}
	}

	return nil
}

func ListModules() error {
	fmt.Println("Available modules:")
	for _, module := range GetAllModuleNames() {
		fmt.Printf("  - %s\n", module)
	}
	return nil
}

func runModule(module ModuleEntry, plat platform.Platform) error {
	fmt.Printf("Running %s\n", module.Name)
	return module.Runner(plat)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	result := []string{}
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
