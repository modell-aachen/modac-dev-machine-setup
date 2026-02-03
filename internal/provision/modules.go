package provision

var allModules = []string{
	"packages",
	"setup-envs",
	"asdf-packages",
	"asdf",
	"kubectl-krew",
	"setup-k8s-cluster",
	"k3s-network",
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
	for _, part := range splitString(s, ',') {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s string, sep rune) []string {
	result := []string{}
	current := ""
	for _, ch := range s {
		if ch == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" || len(result) > 0 {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && isSpace(s[start]) {
		start++
	}

	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
