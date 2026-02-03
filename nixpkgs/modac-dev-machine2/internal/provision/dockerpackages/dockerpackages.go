package dockerpackages

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine2/internal/platform"
)

// Run installs Docker packages
func Run(plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin()
	case platform.Ubuntu:
		return runUbuntu()
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func runDarwin() error {
	// Check if docker is installed via brew
	brewListCmd := exec.Command("brew", "list")
	output, err := brewListCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list brew packages: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	dockerInstalled := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "docker" {
			dockerInstalled = true
			break
		}
	}

	if !dockerInstalled {
		// Install docker packages
		fmt.Println("Installing Docker packages...")
		installCmd := exec.Command("brew", "install",
			"docker",
			"docker-buildx",
			"docker-compose",
			"docker-completion")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install docker packages: %w", err)
		}

		// Set docker context to orbstack
		fmt.Println("Setting docker context to orbstack...")
		contextCmd := exec.Command("docker", "context", "use", "orbstack")
		contextCmd.Stdout = os.Stdout
		contextCmd.Stderr = os.Stderr
		if err := contextCmd.Run(); err != nil {
			return fmt.Errorf("failed to set docker context: %w", err)
		}
	}

	// Setup docker config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dockerConfigPath := filepath.Join(homeDir, ".docker", "config.json")
	needsConfig := false

	if !fileExists(dockerConfigPath) {
		needsConfig = true
	} else {
		content, err := os.ReadFile(dockerConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read docker config: %w", err)
		}
		if !bytes.Contains(content, []byte("cliPluginsExtraDirs")) {
			needsConfig = true
		}
	}

	if needsConfig {
		homebrewPrefix := os.Getenv("HOMEBREW_PREFIX")
		if homebrewPrefix == "" {
			homebrewPrefix = "/opt/homebrew" // Default for Apple Silicon
		}

		dockerDir := filepath.Join(homeDir, ".docker")
		if err := os.MkdirAll(dockerDir, 0755); err != nil {
			return fmt.Errorf("failed to create .docker directory: %w", err)
		}

		configContent := fmt.Sprintf("{\n\t\"cliPluginsExtraDirs\": [ \"%s/lib/docker/cli-plugins\" ]\n}\n",
			homebrewPrefix)
		if err := os.WriteFile(dockerConfigPath, []byte(configContent), 0644); err != nil {
			return fmt.Errorf("failed to write docker config: %w", err)
		}
	}

	return nil
}

func runUbuntu() error {
	// Check if running in distrobox
	containerID := os.Getenv("CONTAINER_ID")
	if containerID != "" {
		// Check if docker command exists
		_, err := exec.LookPath("docker")
		if err != nil {
			// Docker not found, create symlinks
			fmt.Println("Running inside a distrobox, linking docker")

			// Create symlink for docker
			symlinkCmd := exec.Command("sudo", "ln", "-s", "/usr/bin/distrobox-host-exec", "/usr/local/bin/docker")
			symlinkCmd.Stdout = os.Stdout
			symlinkCmd.Stderr = os.Stderr
			if err := symlinkCmd.Run(); err != nil {
				return fmt.Errorf("failed to create docker symlink: %w", err)
			}

			// Link .docker directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}

			user := os.Getenv("USER")
			hostDockerDir := fmt.Sprintf("/run/host/home/%s/.docker", user)
			localDockerDir := filepath.Join(homeDir, ".docker")

			linkCmd := exec.Command("ln", "-s", hostDockerDir, localDockerDir)
			linkCmd.Stdout = os.Stdout
			linkCmd.Stderr = os.Stderr
			if err := linkCmd.Run(); err != nil {
				return fmt.Errorf("failed to link .docker directory: %w", err)
			}
		} else {
			fmt.Println("Running inside a distrobox, skipping docker install")
		}
		return nil
	}

	// Not in distrobox, install docker normally

	// Check and install GPG key
	if !fileExists("/etc/apt/keyrings/docker.asc") {
		fmt.Println("Installing docker GPG key...")

		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "ca-certificates", "curl")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install prerequisites: %w", err)
		}

		mkdirCmd := exec.Command("sudo", "install", "-m", "0755", "-d", "/etc/apt/keyrings")
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("failed to create keyrings directory: %w", err)
		}

		curlCmd := exec.Command("sudo", "curl", "-fsSL",
			"https://download.docker.com/linux/ubuntu/gpg",
			"-o", "/etc/apt/keyrings/docker.asc")
		curlCmd.Stdout = os.Stdout
		curlCmd.Stderr = os.Stderr
		if err := curlCmd.Run(); err != nil {
			return fmt.Errorf("failed to download docker GPG key: %w", err)
		}

		chmodCmd := exec.Command("sudo", "chmod", "a+r", "/etc/apt/keyrings/docker.asc")
		if err := chmodCmd.Run(); err != nil {
			return fmt.Errorf("failed to set key permissions: %w", err)
		}
	}

	// Check and add docker repository
	if !fileExists("/etc/apt/sources.list.d/docker.list") {
		fmt.Println("Adding docker repository...")

		// Get Ubuntu codename
		lsbCmd := exec.Command("lsb_release", "-cs")
		codename, err := lsbCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get Ubuntu codename: %w", err)
		}

		repoLine := fmt.Sprintf("deb [signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu %s stable",
			strings.TrimSpace(string(codename)))

		teeCmd := exec.Command("sudo", "tee", "/etc/apt/sources.list.d/docker.list")
		teeCmd.Stdin = strings.NewReader(repoLine + "\n")
		teeCmd.Stdout = os.Stdout
		if err := teeCmd.Run(); err != nil {
			return fmt.Errorf("failed to add docker repository: %w", err)
		}

		updateCmd := exec.Command("sudo", "apt", "update")
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update apt: %w", err)
		}
	}

	// Install docker packages
	fmt.Println("Installing Docker packages...")
	installCmd := exec.Command("sudo", "apt-get", "install", "-y",
		"docker-ce",
		"docker-ce-cli",
		"containerd.io",
		"docker-buildx-plugin",
		"docker-compose-plugin")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install docker packages: %w", err)
	}

	// Check if user is in docker group
	user := os.Getenv("USER")
	idCmd := exec.Command("id", "-nG", user)
	groups, err := idCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}

	if !strings.Contains(string(groups), "docker") {
		fmt.Printf("Adding user %s to docker group...\n", user)
		usermodCmd := exec.Command("sudo", "usermod", "-aG", "docker", user)
		usermodCmd.Stdout = os.Stdout
		usermodCmd.Stderr = os.Stderr
		if err := usermodCmd.Run(); err != nil {
			return fmt.Errorf("failed to add user to docker group: %w", err)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
