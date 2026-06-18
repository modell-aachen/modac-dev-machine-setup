package nixconf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

const (
	experimentalFeaturesKey  = "experimental-features"
	experimentalFeaturesLine = "experimental-features = nix-command flakes"
)

// requiredFeatures are the features that must be enabled for flakes and the
// modern nix CLI to work.
var requiredFeatures = []string{"nix-command", "flakes"}

// Run ensures ~/.config/nix/nix.conf enables the nix-command and flakes
// experimental features.
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	confPath := filepath.Join(homeDir, ".config", "nix", "nix.conf")

	if !util.FileExists(confPath) {
		return createConf(out, confPath)
	}

	return updateConf(out, confPath)
}

func createConf(out *output.Context, confPath string) error {
	if err := os.MkdirAll(filepath.Dir(confPath), 0755); err != nil {
		return fmt.Errorf("failed to create nix config directory: %w", err)
	}

	out.Step("Creating nix.conf with experimental features")
	content := experimentalFeaturesLine + "\n"
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write nix.conf: %w", err)
	}

	return nil
}

func updateConf(out *output.Context, confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return fmt.Errorf("failed to read nix.conf: %w", err)
	}

	if hasRequiredFeatures(string(data)) {
		out.Skipped("nix.conf already enables nix-command and flakes")
		return nil
	}

	out.Step("Appending experimental features to nix.conf")
	content := string(data)
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += experimentalFeaturesLine + "\n"

	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write nix.conf: %w", err)
	}

	return nil
}

// hasRequiredFeatures reports whether an existing experimental-features setting
// already enables all required features.
func hasRequiredFeatures(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, found := strings.Cut(trimmed, "=")
		if !found || strings.TrimSpace(key) != experimentalFeaturesKey {
			continue
		}

		enabled := strings.Fields(value)
		if containsAll(enabled, requiredFeatures) {
			return true
		}
	}

	return false
}

func containsAll(haystack, needles []string) bool {
	set := make(map[string]bool, len(haystack))
	for _, h := range haystack {
		set[h] = true
	}
	for _, n := range needles {
		if !set[n] {
			return false
		}
	}
	return true
}
