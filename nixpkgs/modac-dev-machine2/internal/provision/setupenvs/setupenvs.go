package setupenvs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine2/internal/output"
	"github.com/modell-aachen/machine2/internal/platform"
	"github.com/modell-aachen/machine2/internal/util"
)

type DevboxConfig struct {
	EnvFrom       string            `json:"env_from,omitempty"`
	OpSecretsTpl  map[string]string `json:"op_secrets_tpl,omitempty"`
	OtherFields   map[string]any    `json:"-"`
}

// Run sets up environment variables and secrets integration
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Get devbox global path
	devboxPath, err := getDevboxGlobalPath()
	if err != nil {
		return fmt.Errorf("failed to get devbox global path: %w", err)
	}

	configPath := filepath.Join(devboxPath, "devbox.json")
	tmpPath := filepath.Join(devboxPath, "tmp.json")

	// Remove tmp file if it exists
	if util.FileExists(tmpPath) {
		if err := os.Remove(tmpPath); err != nil {
			return fmt.Errorf("failed to remove tmp file: %w", err)
		}
	}

	// Read current config
	config, err := readDevboxConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read devbox config: %w", err)
	}

	// Add env_from if not present
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	envFromPath := filepath.Join(homeDir, ".secrets", ".env")
	if config.EnvFrom != envFromPath {
		out.Step("Adding env_from to devbox config")
		config.EnvFrom = envFromPath
		if err := writeDevboxConfig(configPath, config); err != nil {
			return fmt.Errorf("failed to write env_from: %w", err)
		}
	} else {
		out.Skipped("env_from already configured")
	}

	// Get templates directory
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	// Merge op_secrets_tpl from template
	out.Step("Merging op_secrets_tpl with template")
	templatePath := filepath.Join(templatesDir, "devbox.json")
	template, err := readDevboxConfig(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Merge template secrets into config (only add new keys)
	if config.OpSecretsTpl == nil {
		config.OpSecretsTpl = make(map[string]string)
	}
	for key, value := range template.OpSecretsTpl {
		if _, exists := config.OpSecretsTpl[key]; !exists {
			config.OpSecretsTpl[key] = value
		}
	}

	if err := writeDevboxConfig(configPath, config); err != nil {
		return fmt.Errorf("failed to write merged config: %w", err)
	}

	// Remove "export " from $HOME/.env if present
	envFile := filepath.Join(homeDir, ".env")
	if util.FileExists(envFile) {
		content, err := os.ReadFile(envFile)
		if err != nil {
			return fmt.Errorf("failed to read .env file: %w", err)
		}

		if bytes.Contains(content, []byte("export ")) {
			out.Step("Removing export from .env file")
			modified := bytes.ReplaceAll(content, []byte("\nexport "), []byte("\n"))
			modified = bytes.TrimPrefix(modified, []byte("export "))
			if err := os.WriteFile(envFile, modified, 0644); err != nil {
				return fmt.Errorf("failed to write .env file: %w", err)
			}
		}
	}

	// Create secrets directory
	secretsDir := filepath.Join(homeDir, ".secrets")
	if err := os.MkdirAll(secretsDir, 0755); err != nil {
		return fmt.Errorf("failed to create secrets directory: %w", err)
	}

	// Generate env.tpl from op_secrets_tpl
	out.Step("Generating env template")
	envTplPath := filepath.Join(secretsDir, "env.tpl")
	if err := generateEnvTemplate(config.OpSecretsTpl, envTplPath); err != nil {
		return fmt.Errorf("failed to generate env template: %w", err)
	}

	// Inject secrets using 1Password CLI
	out.Step("Injecting secrets from 1Password")
	envPath := filepath.Join(secretsDir, ".env")
	if err := out.RunCommand("op", "inject", "--in-file", envTplPath, "--out-file", envPath, "--force"); err != nil {
		return fmt.Errorf("failed to inject secrets: %w", err)
	}

	return nil
}

func getDevboxGlobalPath() (string, error) {
	cmd := exec.Command("devbox", "global", "path")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getTemplatesDir() (string, error) {
	// Get the directory of the currently running binary
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Binary is in bin/, templates should be in ../share/machine2/templates/
	binDir := filepath.Dir(exePath)
	shareDir := filepath.Join(binDir, "..", "share", "machine2", "templates")

	// Check if the share directory exists
	if util.FileExists(shareDir) {
		return shareDir, nil
	}

	// Fallback: check if we're running from the repo (development mode)
	repoTemplatesDir := filepath.Join(binDir, "..", "scripts", "templates")
	if util.FileExists(repoTemplatesDir) {
		return repoTemplatesDir, nil
	}

	return "", fmt.Errorf("templates directory not found")
}

func readDevboxConfig(path string) (*DevboxConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// First unmarshal into a map to preserve unknown fields
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	config := &DevboxConfig{}

	// Extract known fields
	if envFrom, ok := raw["env_from"].(string); ok {
		config.EnvFrom = envFrom
	}

	if opSecrets, ok := raw["op_secrets_tpl"].(map[string]any); ok {
		config.OpSecretsTpl = make(map[string]string)
		for k, v := range opSecrets {
			if str, ok := v.(string); ok {
				config.OpSecretsTpl[k] = str
			}
		}
	}

	// Store other fields
	config.OtherFields = raw

	return config, nil
}

func writeDevboxConfig(path string, config *DevboxConfig) error {
	// Start with other fields
	output := config.OtherFields
	if output == nil {
		output = make(map[string]any)
	}

	// Update known fields
	if config.EnvFrom != "" {
		output["env_from"] = config.EnvFrom
	}

	if config.OpSecretsTpl != nil {
		output["op_secrets_tpl"] = config.OpSecretsTpl
	}

	// Marshal with indentation
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func generateEnvTemplate(secrets map[string]string, outputPath string) error {
	var lines []string
	for key, value := range secrets {
		// Shell-quote the value (simple implementation)
		quotedValue := shellQuote(value)
		lines = append(lines, fmt.Sprintf("%s=%s", key, quotedValue))
	}

	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}

func shellQuote(s string) string {
	// Simple shell quoting: wrap in single quotes and escape single quotes
	s = strings.ReplaceAll(s, "'", "'\"'\"'")
	return "'" + s + "'"
}
