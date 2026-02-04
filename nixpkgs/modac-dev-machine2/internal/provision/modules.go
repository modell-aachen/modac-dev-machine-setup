package provision

import "strings"

var allModules = []string{
	"packages",
	"setup-envs",
	"asdf-packages",
	"asdf",
	"kubectl-krew",
	"setup-k8s-cluster",
	"node",
	"certificates",
	"setup-dev",
	"completions",
	"claude",
	"github-auth-login",
	"install-modac-shell-helper",
	"orbstack",
	"docker-packages",
	"docker",
}

func GetAllModules() []string {
	return allModules
}

func FilterModules(filter string) []string {
	if filter == "" {
		return allModules
	}

	filterMap := make(map[string]bool)
	for _, name := range splitCSV(filter) {
		filterMap[name] = true
	}

	filtered := []string{}
	for _, module := range allModules {
		if filterMap[module] {
			filtered = append(filtered, module)
		}
	}

	return filtered
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
