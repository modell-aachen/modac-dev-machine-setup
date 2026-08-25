package dockerpackages

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

const dockerRepoFile = "/etc/apt/sources.list.d/docker.list"

// Run installs Docker packages
func Run(out *output.Context, plat platform.Platform) error {
	switch plat {
	case platform.Darwin:
		return runDarwin(out)
	case platform.Ubuntu:
		return runUbuntu(out)
	default:
		return fmt.Errorf("unsupported platform: %s", plat)
	}
}

func writeDockerRepo(out *output.Context) error {
	archCmd := exec.Command("dpkg", "--print-architecture")
	arch, err := archCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get dpkg architecture: %w", err)
	}

	codenameCmd := exec.Command("lsb_release", "-cs")
	codename, err := codenameCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Ubuntu codename: %w", err)
	}

	repoLine := dockerRepoLine(strings.TrimSpace(string(arch)), strings.TrimSpace(string(codename)))

	teeCmd := exec.Command("sudo", "tee", dockerRepoFile)
	teeCmd.Stdin = strings.NewReader(repoLine + "\n")
	teeCmd.Stdout = os.Stdout
	if err := teeCmd.Run(); err != nil {
		return fmt.Errorf("failed to add docker repository: %w", err)
	}

	out.Step("Updating apt")
	if err := out.RunCommand("sudo", "apt", "update"); err != nil {
		return fmt.Errorf("failed to update apt: %w", err)
	}

	return nil
}

func dockerRepoLine(arch, codename string) string {
	return fmt.Sprintf("deb [arch=%s signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu %s stable",
		arch, codename)
}

// needsArchPin reports whether a repository entry lets apt query every enabled
// architecture. Docker only publishes some of them, so an unpinned entry makes apt
// report a missing index for the others on each update.
func needsArchPin(repo string) bool {
	return !strings.Contains(repo, "arch=")
}

func runDarwin(out *output.Context) error {
	// Check if docker is installed via brew
	brewListCmd := exec.Command("brew", "list")
	brewOutput, err := brewListCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list brew packages: %w", err)
	}

	lines := strings.Split(string(brewOutput), "\n")
	dockerInstalled := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "docker" {
			dockerInstalled = true
			break
		}
	}

	if !dockerInstalled {
		// Install docker packages
		out.Step("Installing Docker packages")
		if err := out.RunCommand("brew", "install", "docker", "docker-buildx", "docker-compose", "docker-completion"); err != nil {
			return fmt.Errorf("failed to install docker packages: %w", err)
		}

		// Set docker context to orbstack
		out.Step("Setting docker context to orbstack")
		if err := out.RunCommand("docker", "context", "use", "orbstack"); err != nil {
			return fmt.Errorf("failed to set docker context: %w", err)
		}
	} else {
		out.Skipped("Docker packages already installed")
	}

	// Setup docker config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dockerConfigPath := filepath.Join(homeDir, ".docker", "config.json")
	needsConfig := false

	if !util.FileExists(dockerConfigPath) {
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
		out.Step("Setting up Docker config")
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
	} else {
		out.Skipped("Docker config already set up")
	}

	return nil
}

func runUbuntu(out *output.Context) error {
	// Check if running in distrobox
	containerID := os.Getenv("CONTAINER_ID")
	if containerID != "" {
		// Check if docker command exists
		if _, err := exec.LookPath("docker"); err != nil {
			// Docker not found, link it to the host
			out.Step("Running inside a distrobox, linking docker")
			if err := out.RunCommand("sudo", "ln", "-s", "/usr/local/bin/distrobox-host-exec-with-env", "/usr/local/bin/docker"); err != nil {
				return fmt.Errorf("failed to create docker symlink: %w", err)
			}
		} else {
			out.Skipped("Running inside a distrobox, docker already available")
		}

		// Point DOCKER_CONFIG at the host's .docker directory, independent of
		// whether docker itself was just linked.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		user := os.Getenv("USER")
		hostDockerDir := fmt.Sprintf("/run/host/home/%s/.docker", user)
		bashrcPath := filepath.Join(homeDir, ".bashrc")
		exportLine := fmt.Sprintf("export DOCKER_CONFIG=%q", hostDockerDir)

		out.Step("Setting DOCKER_CONFIG in .bashrc")
		if err := appendLineIfMissing(bashrcPath, exportLine); err != nil {
			return fmt.Errorf("failed to update .bashrc: %w", err)
		}

		return nil
	}

	// Not in distrobox, install docker normally

	// Check and install GPG key
	if !util.FileExists("/etc/apt/keyrings/docker.asc") {
		out.Step("Installing docker GPG key")

		if err := out.RunCommand("sudo", "apt-get", "install", "-y", "ca-certificates", "curl"); err != nil {
			return fmt.Errorf("failed to install prerequisites: %w", err)
		}

		mkdirCmd := exec.Command("sudo", "install", "-m", "0755", "-d", "/etc/apt/keyrings")
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("failed to create keyrings directory: %w", err)
		}

		if err := out.RunCommand("sudo", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "-o", "/etc/apt/keyrings/docker.asc"); err != nil {
			return fmt.Errorf("failed to download docker GPG key: %w", err)
		}

		chmodCmd := exec.Command("sudo", "chmod", "a+r", "/etc/apt/keyrings/docker.asc")
		if err := chmodCmd.Run(); err != nil {
			return fmt.Errorf("failed to set key permissions: %w", err)
		}
	} else {
		out.Skipped("Docker GPG key already installed")
	}

	// Check and add docker repository
	existingRepo, err := os.ReadFile(dockerRepoFile)
	switch {
	case os.IsNotExist(err):
		out.Step("Adding docker repository")
		if err := writeDockerRepo(out); err != nil {
			return err
		}
	case err != nil:
		return fmt.Errorf("failed to read docker repository: %w", err)
	case needsArchPin(string(existingRepo)):
		out.Step("Pinning docker repository to the host architecture")
		if err := writeDockerRepo(out); err != nil {
			return err
		}
	default:
		out.Skipped("Docker repository already added")
	}

	// Install docker packages
	out.Step("Installing Docker packages")
	if err := out.RunCommand("sudo", "apt-get", "install", "-y",
		"docker-ce", "docker-ce-cli", "containerd.io",
		"docker-buildx-plugin", "docker-compose-plugin"); err != nil {
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
		out.Step(fmt.Sprintf("Adding user %s to docker group", user))
		if err := out.RunCommand("sudo", "usermod", "-aG", "docker", user); err != nil {
			return fmt.Errorf("failed to add user to docker group: %w", err)
		}
	} else {
		out.Skipped("User already in docker group")
	}

	// Check if docker can be run without sudo
	if !canAccessDockerWithoutSudo() {
		out.Info("Please logout and login again to use docker without sudo")
		return fmt.Errorf("logout required - please logout and login again")
	}

	return nil
}

// appendLineIfMissing appends line to the file at path (creating it if needed),
// unless the line is already present. A trailing newline is ensured.
func appendLineIfMissing(path, line string) error {
	if util.FileExists(path) {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, existing := range strings.Split(string(content), "\n") {
			if strings.TrimSpace(existing) == line {
				return nil
			}
		}
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s\n", line); err != nil {
		return err
	}
	return nil
}

func canAccessDockerWithoutSudo() bool {
	// Try to run a simple docker command without sudo
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	return err == nil
}
