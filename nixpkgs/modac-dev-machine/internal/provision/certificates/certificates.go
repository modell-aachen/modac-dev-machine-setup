package certificates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run sets up development certificates using mkcert
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	// Get QWIKI_DEVELOPMENT_ROOT_CA from environment
	rootCADir := os.Getenv("QWIKI_DEVELOPMENT_ROOT_CA")
	if rootCADir == "" {
		return fmt.Errorf("QWIKI_DEVELOPMENT_ROOT_CA environment variable not set")
	}

	// Check if root CA exists
	rootCAPem := filepath.Join(rootCADir, "rootCA.pem")
	if !util.FileExists(rootCAPem) {
		out.Step(fmt.Sprintf("Generating root CA in '%s'", rootCADir))
		cmd := exec.Command("mkcert", "-install")
		cmd.Env = append(os.Environ(), fmt.Sprintf("CAROOT=%s", rootCADir))
		if err := out.RunCmd(cmd); err != nil {
			return fmt.Errorf("failed to generate root CA: %w", err)
		}
	} else {
		out.Skipped("Root CA already exists")
	}

	// Generate localhost certificate
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	host := "localhost"
	location := filepath.Join(homeDir, "certs", host)
	certPath := filepath.Join(location, host+".pem")

	if !util.FileExists(certPath) {
		out.Step(fmt.Sprintf("Generating certificate for '%s'", host))

		// Create directory
		if err := os.MkdirAll(location, 0755); err != nil {
			return fmt.Errorf("failed to create certificate directory: %w", err)
		}

		// Generate certificate in that directory
		cmd := exec.Command("mkcert", host)
		cmd.Dir = location
		cmd.Env = append(os.Environ(), fmt.Sprintf("CAROOT=%s", rootCADir))
		if err := out.RunCmd(cmd); err != nil {
			return fmt.Errorf("failed to generate certificate for %s: %w", host, err)
		}
	} else {
		out.Skipped("Certificate for localhost already exists")
	}

	return nil
}
