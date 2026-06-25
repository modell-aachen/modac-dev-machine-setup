package provision

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

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
	Filter string
}

// ModuleEntry represents a provisioning module with its name and runner function
type ModuleEntry struct {
	Name   string
	Runner func(*output.Context, platform.Platform) error
}

// devboxUpdateModuleName is the module after which the executor may re-exec into
// a newly installed machine binary (see reexecAfterUpdate).
const devboxUpdateModuleName = "devbox-update"

// allModules defines the ordered list of all provisioning modules.
//
// Every module listed before devboxUpdateModuleName must be idempotent: when an
// update installs a newer binary, the executor re-execs and the provision runs
// from the top again on the new binary, so those modules run twice.
var allModules = []ModuleEntry{
	{"nix-conf", nixconf.Run},
	{devboxUpdateModuleName, devboxupdate.Run},
	{"onepassword", onepassword.Run},
	{"restore-backup", restorebackup.Run},
	{"packages", packages.Run},
	{"setup-envs", setupenvs.Run},
	{"asdf-packages", asdfpackages.Run},
	{"asdf", asdf.Run},
	{"kubectl-krew", kubectlkrew.Run},
	{"setup-k8s-cluster", setupk8scluster.Run},
	{"node", node.Run},
	{"nssdb", nssdb.Run},
	{"certificates", certificates.Run},
	{"setup-dev", setupdev.Run},
	{"completions", completions.Run},
	{"claude", claude.Run},
	{"github-auth-login", githubauthlogin.Run},
	{"gcloud-workforce-login", gcloudworkforcelogin.Run},
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
	// Create output context for nice formatting and logging
	out, err := output.New()
	if err != nil {
		return fmt.Errorf("failed to initialize output: %w", err)
	}
	defer out.Close()

	fmt.Printf("Machine Provisioning\n")
	fmt.Printf("Log file: %s\n\n", out.LogPath())

	plat, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("platform detection failed: %w", err)
	}

	modules := FilterModules(opts.Filter)

	// Resolved path of the binary we're running. After devbox-update we compare
	// it against the machine now on PATH and re-exec if it changed, so the rest
	// of the modules run on the updated binary.
	self := currentBinary()

	// Run all modules
	for _, module := range modules {
		if err := runModule(out, module, plat); err != nil {
			out.PrintError(fmt.Errorf("module %s failed: %w", module.Name, err))
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

// currentBinary returns the resolved path of the running machine binary, or ""
// if it can't be determined. An empty result makes reexecAfterUpdate treat the
// post-update binary as changed — better to reload than to silently keep running
// a possibly stale binary.
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

// binaryChanged reports whether resolved is a different machine binary than the
// one currently running (self). An unresolvable resolved ("") counts as no
// change so the caller keeps running rather than re-exec into nothing.
func binaryChanged(self, resolved string) bool {
	return resolved != "" && resolved != self
}

// reexecAfterUpdate restarts the provision under the machine binary now on PATH
// when devbox-update installed a different one, so the remaining modules run on
// it. devbox.GlobalUpdate already refreshed PATH into this process's env, which
// os.Environ() carries to the new image. The RestartGuardEnv guard makes the
// post-re-exec pass a no-op here, so this cannot loop. On success it never
// returns — syscall.Exec replaces the process image.
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
	return nil // unreachable on success: syscall.Exec replaces the process image
}

func ListModules() error {
	fmt.Println("Available modules:")
	for _, module := range GetAllModuleNames() {
		fmt.Printf("  - %s\n", module)
	}
	return nil
}

func runModule(out *output.Context, module ModuleEntry, plat platform.Platform) error {
	out.StartModule(module.Name)
	if err := module.Runner(out, plat); err != nil {
		out.Failure(fmt.Sprintf("Module failed: %v", err))
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
