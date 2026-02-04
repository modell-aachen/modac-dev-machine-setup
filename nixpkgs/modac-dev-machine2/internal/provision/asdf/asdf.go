package asdf

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine2/internal/platform"
)

type DevboxConfig struct {
	Env         map[string]string `json:"env,omitempty"`
	OtherFields map[string]any    `json:"-"`
}

// Run configures asdf version manager
func Run(plat platform.Platform) error {
	_ = plat
	// Add asdf plugins
	if err := addPlugin("erlang"); err != nil {
		return fmt.Errorf("failed to add erlang plugin: %w", err)
	}

	if err := addPlugin("elixir"); err != nil {
		return fmt.Errorf("failed to add elixir plugin: %w", err)
	}

	// Get ASDF_DIR
	asdfDir, err := getAsdfDir()
	if err != nil {
		return fmt.Errorf("failed to get ASDF_DIR: %w", err)
	}

	// Get devbox config path
	devboxPath, err := getDevboxGlobalPath()
	if err != nil {
		return fmt.Errorf("failed to get devbox global path: %w", err)
	}

	configPath := filepath.Join(devboxPath, "devbox.json")

	// Read config
	config, err := readDevboxConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read devbox config: %w", err)
	}

	// Check if ASDF_DIR is already set
	if config.Env != nil {
		if existingDir, exists := config.Env["ASDF_DIR"]; exists && existingDir == asdfDir {
			// Already configured
			return nil
		}
	}

	// Add ASDF_DIR to env
	fmt.Println("Adding ASDF_DIR to devbox config")
	if config.Env == nil {
		config.Env = make(map[string]string)
	}
	config.Env["ASDF_DIR"] = asdfDir

	// Write config
	if err := writeDevboxConfig(configPath, config); err != nil {
		return fmt.Errorf("failed to write devbox config: %w", err)
	}

	return nil
}

func addPlugin(name string) error {
	fmt.Printf("Adding asdf plugin: %s\n", name)
	cmd := exec.Command("asdf", "plugin", "add", name)
	// Ignore errors - plugin might already be added
	_ = cmd.Run()
	return nil
}

func getAsdfDir() (string, error) {
	cmd := exec.Command("asdf", "info")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output to find ASDF_DATA_DIR=...
	// The bash script greps for "ASDF_DIR" which matches ASDF_DATA_DIR
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ASDF_DATA_DIR=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("ASDF_DATA_DIR not found in asdf info output")
}

func getDevboxGlobalPath() (string, error) {
	cmd := exec.Command("devbox", "global", "path")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
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

	// Extract env field
	if envRaw, ok := raw["env"].(map[string]any); ok {
		config.Env = make(map[string]string)
		for k, v := range envRaw {
			if str, ok := v.(string); ok {
				config.Env[k] = str
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

	// Update env field
	if config.Env != nil {
		output["env"] = config.Env
	}

	// Marshal with indentation
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
