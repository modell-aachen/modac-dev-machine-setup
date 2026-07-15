package provision

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/modell-aachen/machine/internal/config"
	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/provision/asdf"
	"github.com/modell-aachen/machine/internal/provision/asdfpackages"
	"github.com/modell-aachen/machine/internal/provision/certificates"
	"github.com/modell-aachen/machine/internal/provision/claude"
	"github.com/modell-aachen/machine/internal/provision/completions"
	"github.com/modell-aachen/machine/internal/provision/devboxupdate"
	"github.com/modell-aachen/machine/internal/provision/docker"
	"github.com/modell-aachen/machine/internal/provision/dockerpackages"
	"github.com/modell-aachen/machine/internal/provision/gcloudworkforcelogin"
	"github.com/modell-aachen/machine/internal/provision/githubauthlogin"
	"github.com/modell-aachen/machine/internal/provision/installmodacshellhelper"
	"github.com/modell-aachen/machine/internal/provision/kubectlkrew"
	"github.com/modell-aachen/machine/internal/provision/nixconf"
	"github.com/modell-aachen/machine/internal/provision/node"
	"github.com/modell-aachen/machine/internal/provision/nssdb"
	"github.com/modell-aachen/machine/internal/provision/onepassword"
	"github.com/modell-aachen/machine/internal/provision/orbstack"
	"github.com/modell-aachen/machine/internal/provision/packages"
	"github.com/modell-aachen/machine/internal/provision/restorebackup"
	"github.com/modell-aachen/machine/internal/provision/setupdev"
	"github.com/modell-aachen/machine/internal/provision/setupenvs"
	"github.com/modell-aachen/machine/internal/provision/setupk8scluster"
)

type Options struct {
	Filter  string
	Profile config.Profile
}

type ModuleEntry struct {
	Name    string
	Runner  func(*output.Context, platform.Platform) error
	Service bool // part of the minimal service machine profile
}

const devboxUpdateModuleName = "devbox-update"

var allModules = []ModuleEntry{
	{Name: "nix-conf", Runner: nixconf.Run, Service: true},
	{Name: devboxUpdateModuleName, Runner: devboxupdate.Run, Service: true},
	{Name: "onepassword", Runner: onepassword.Run, Service: true},
	{Name: "restore-backup", Runner: restorebackup.Run},
	{Name: "packages", Runner: packages.Run},
	{Name: "setup-envs", Runner: setupenvs.Run},
	{Name: "asdf-packages", Runner: asdfpackages.Run},
	{Name: "asdf", Runner: asdf.Run},
	{Name: "kubectl-krew", Runner: kubectlkrew.Run, Service: true},
	{Name: "setup-k8s-cluster", Runner: setupk8scluster.Run, Service: true},
	{Name: "node", Runner: node.Run},
	{Name: "nssdb", Runner: nssdb.Run},
	{Name: "certificates", Runner: certificates.Run},
	{Name: "setup-dev", Runner: setupdev.Run},
	{Name: "completions", Runner: completions.Run},
	{Name: "claude", Runner: claude.Run},
	{Name: "github-auth-login", Runner: githubauthlogin.Run, Service: true},
	{Name: "gcloud-workforce-login", Runner: gcloudworkforcelogin.Run, Service: true},
	{Name: "install-modac-shell-helper", Runner: installmodacshellhelper.Run, Service: true},
	{Name: "orbstack", Runner: orbstack.Run},
	{Name: "docker-packages", Runner: dockerpackages.Run},
	{Name: "docker", Runner: docker.Run},
}

func GetAllModuleNames() []string {
	names := make([]string, len(allModules))
	for i, module := range allModules {
		names[i] = module.Name
	}
	return names
}

// ModulesFor returns the modules for a profile, optionally narrowed by a
// comma-separated filter. Definition order is always preserved; filter names
// outside the profile are dropped silently.
func ModulesFor(profile config.Profile, filter string) []ModuleEntry {
	filterMap := make(map[string]bool)
	for _, name := range splitCSV(filter) {
		filterMap[name] = true
	}

	modules := []ModuleEntry{}
	for _, module := range allModules {
		if profile == config.ProfileService && !module.Service {
			continue
		}
		if filter != "" && !filterMap[module.Name] {
			continue
		}
		modules = append(modules, module)
	}

	return modules
}

func Execute(opts *Options) error {
	out, err := output.New()
	if err != nil {
		return fmt.Errorf("failed to initialize output: %w", err)
	}
	defer out.Close()

	fmt.Printf("Machine Provisioning (%s profile)\n", opts.Profile)
	fmt.Printf("Log file: %s\n\n", out.LogPath())

	plat, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("platform detection failed: %w", err)
	}

	modules := ModulesFor(opts.Profile, opts.Filter)

	self := currentBinary()

	for _, module := range modules {
		if err := runModule(out, module, plat); err != nil {
			return fmt.Errorf("module %s failed: %w", module.Name, err)
		}
		if module.Name == devboxUpdateModuleName {
			if err := reexecAfterUpdate(out, self); err != nil {
				out.PrintError(err)
				return err
			}
		}
	}

	out.PrintSummary()
	return nil
}

func currentBinary() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved
	}
	return exe
}

func binaryChanged(self, resolved string) bool {
	return resolved != "" && resolved != self
}

func reexecAfterUpdate(out *output.Context, self string) error {
	if devboxupdate.AlreadyUpdated() {
		return nil
	}

	bin, err := exec.LookPath("machine")
	if err != nil {
		out.Skipped("machine not found on PATH after update; continuing on the current binary")
		return nil
	}

	resolved, err := filepath.EvalSymlinks(bin)
	if err != nil {
		resolved = bin
	}

	if !binaryChanged(self, resolved) {
		out.Skipped("machine binary unchanged; no restart needed")
		return nil
	}

	out.Step("Restarting machine with the updated binary")
	env := append(os.Environ(), devboxupdate.RestartGuardEnv+"=1")
	if err := syscall.Exec(bin, os.Args, env); err != nil {
		return fmt.Errorf("failed to restart machine after update: %w", err)
	}
	return nil
}

func ListModules() error {
	fmt.Println("Available modules (* = included in service profile):")
	for _, module := range allModules {
		marker := ""
		if module.Service {
			marker = " *"
		}
		fmt.Printf("  - %s%s\n", module.Name, marker)
	}
	return nil
}

func runModule(out *output.Context, module ModuleEntry, plat platform.Platform) error {
	out.StartModule(module.Name)
	if err := module.Runner(out, plat); err != nil {
		out.Failure("Module failed")
		return err
	}
	out.Success("Module completed")
	return nil
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
