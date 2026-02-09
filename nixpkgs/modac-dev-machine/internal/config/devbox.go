package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BackupConfig struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Vault string `json:"vault"`
	Type  string `json:"type"` // "file" or "directory"
}

type DevboxConfig struct {
	Backups []BackupConfig `json:"backups"`
}

func DevboxPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "share", "devbox", "global", "default", "devbox.json"), nil
}

func LoadDevbox() (*DevboxConfig, error) {
	path, err := DevboxPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read devbox.json: %w", err)
	}

	var config DevboxConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse devbox.json: %w", err)
	}

	return &config, nil
}

func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "$HOME/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(homeDir, path[6:]), nil
	}
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}
