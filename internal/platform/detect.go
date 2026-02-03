package platform

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Platform string

const (
	Darwin Platform = "Darwin"
	Ubuntu Platform = "Ubuntu"
)

func (p Platform) String() string {
	return string(p)
}

func Detect() (Platform, error) {
	switch runtime.GOOS {
	case "darwin":
		return Darwin, nil
	case "linux":
		return detectLinuxDistro()
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func detectLinuxDistro() (Platform, error) {
	// Check /etc/os-release for distribution info
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("failed to read /etc/os-release: %w", err)
	}

	content := string(data)
	if strings.Contains(strings.ToLower(content), "ubuntu") {
		return Ubuntu, nil
	}

	return "", fmt.Errorf("unsupported Linux distribution (only Ubuntu is supported)")
}
