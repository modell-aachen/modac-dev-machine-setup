package provision

import (
	"slices"
	"testing"

	"github.com/modell-aachen/machine/internal/config"
)

func moduleNames(modules []ModuleEntry) []string {
	names := make([]string, len(modules))
	for i, module := range modules {
		names[i] = module.Name
	}
	return names
}

func TestModulesFor(t *testing.T) {
	tests := []struct {
		name    string
		profile config.Profile
		filter  string
		want    []string
	}{
		{
			name:    "dev without filter returns all modules",
			profile: config.ProfileDev,
			filter:  "",
			want:    GetAllModuleNames(),
		},
		{
			name:    "dev with filter keeps definition order",
			profile: config.ProfileDev,
			filter:  "packages,nix-conf",
			want:    []string{"nix-conf", "packages"},
		},
		{
			name:    "dev with unknown names drops them silently",
			profile: config.ProfileDev,
			filter:  "nix-conf,does-not-exist",
			want:    []string{"nix-conf"},
		},
		{
			name:    "service without filter returns service modules in order",
			profile: config.ProfileService,
			filter:  "",
			want: []string{
				"nix-conf",
				"devbox-update",
				"onepassword",
				"kubectl-krew",
				"setup-k8s-cluster",
				"github-auth-login",
				"gcloud-workforce-login",
				"install-modac-shell-helper",
			},
		},
		{
			name:    "service with filter intersects with service set",
			profile: config.ProfileService,
			filter:  "kubectl-krew,orbstack,claude",
			want:    []string{"kubectl-krew"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := moduleNames(ModulesFor(tt.profile, tt.filter))
			if !slices.Equal(got, tt.want) {
				t.Errorf("ModulesFor(%q, %q) = %v, want %v", tt.profile, tt.filter, got, tt.want)
			}
		})
	}
}

// modulePrereqs documents the implicit dependencies encoded in allModules
// ordering. Every profile must contain each prerequisite of its modules,
// and the prerequisite must run first.
var modulePrereqs = []struct{ module, requires string }{
	{"devbox-update", "nix-conf"},
	{"restore-backup", "onepassword"},
	{"setup-envs", "onepassword"},
	{"kubectl-krew", "devbox-update"},
	{"setup-k8s-cluster", "onepassword"},
	{"setup-k8s-cluster", "devbox-update"},
	{"certificates", "nssdb"},
	{"install-modac-shell-helper", "github-auth-login"},
	{"install-modac-shell-helper", "devbox-update"},
}

func TestProfilesAreDependencyClosed(t *testing.T) {
	for _, profile := range []config.Profile{config.ProfileDev, config.ProfileService} {
		names := moduleNames(ModulesFor(profile, ""))
		for _, prereq := range modulePrereqs {
			moduleIdx := slices.Index(names, prereq.module)
			if moduleIdx == -1 {
				continue
			}
			requiresIdx := slices.Index(names, prereq.requires)
			if requiresIdx == -1 {
				t.Errorf("profile %q contains %q but not its prerequisite %q", profile, prereq.module, prereq.requires)
				continue
			}
			if requiresIdx > moduleIdx {
				t.Errorf("profile %q runs %q before its prerequisite %q", profile, prereq.module, prereq.requires)
			}
		}
	}
}

func TestPrereqTableReferencesExistingModules(t *testing.T) {
	names := GetAllModuleNames()
	for _, prereq := range modulePrereqs {
		if !slices.Contains(names, prereq.module) {
			t.Errorf("modulePrereqs references unknown module %q", prereq.module)
		}
		if !slices.Contains(names, prereq.requires) {
			t.Errorf("modulePrereqs references unknown prerequisite %q", prereq.requires)
		}
	}
}

func TestBinaryChanged(t *testing.T) {
	tests := []struct {
		name     string
		self     string
		resolved string
		want     bool
	}{
		{"same path means unchanged", "/nix/store/aaa/bin/machine", "/nix/store/aaa/bin/machine", false},
		{"different path means changed", "/nix/store/aaa/bin/machine", "/nix/store/bbb/bin/machine", true},
		{"unresolvable post-update binary counts as unchanged", "/nix/store/aaa/bin/machine", "", false},
		{"unknown current binary forces a reload", "", "/nix/store/bbb/bin/machine", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binaryChanged(tt.self, tt.resolved); got != tt.want {
				t.Errorf("binaryChanged(%q, %q) = %v, want %v", tt.self, tt.resolved, got, tt.want)
			}
		})
	}
}
