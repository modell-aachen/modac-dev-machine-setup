package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile string

const (
	ProfileDev     Profile = "dev"
	ProfileService Profile = "service"
)

func ParseProfile(s string) (Profile, error) {
	switch Profile(s) {
	case ProfileDev, ProfileService:
		return Profile(s), nil
	default:
		return "", fmt.Errorf("unknown profile %q (valid: %s, %s)", s, ProfileDev, ProfileService)
	}
}

func ProfilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".machine", "profile"), nil
}

// LoadProfile reads the persisted machine profile. A missing file means the
// machine was provisioned before profiles existed, so it defaults to dev.
func LoadProfile() (Profile, error) {
	path, err := ProfilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ProfileDev, nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read profile file: %w", err)
	}

	profile, err := ParseProfile(strings.TrimSpace(string(data)))
	if err != nil {
		return "", fmt.Errorf("invalid profile in %s: %w", path, err)
	}
	return profile, nil
}

func SaveProfile(profile Profile) error {
	path, err := ProfilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(profile+"\n"), 0o644); err != nil {
		return fmt.Errorf("failed to write profile file: %w", err)
	}
	return nil
}
