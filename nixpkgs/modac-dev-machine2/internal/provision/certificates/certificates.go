package certificates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/util"
)

// Run sets up development certificates using mkcert
func Run() error {
	// Get QWIKI_DEVELOPMENT_ROOT_CA from environment
	rootCADir := os.Getenv("QWIKI_DEVELOPMENT_ROOT_CA")
	if rootCADir == "" {
		return fmt.Errorf("QWIKI_DEVELOPMENT_ROOT_CA environment variable not set")
	}

	// Check if root CA exists
	rootCAPem := filepath.Join(rootCADir, "rootCA.pem")
	if !util.FileExists(rootCAPem) {
		fmt.Printf("Generating root CA in '%s'\n", rootCADir)
		cmd := exec.Command("mkcert", "-install")
		cmd.Env = append(os.Environ(), fmt.Sprintf("CAROOT=%s", rootCADir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate root CA: %w", err)
		}
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
		fmt.Printf("Generating certificate for '%s' in '%s'\n", host, location)

		// Create directory
		if err := os.MkdirAll(location, 0755); err != nil {
			return fmt.Errorf("failed to create certificate directory: %w", err)
		}

		// Generate certificate in that directory
		cmd := exec.Command("mkcert", host)
		cmd.Dir = location
		cmd.Env = append(os.Environ(), fmt.Sprintf("CAROOT=%s", rootCADir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate certificate for %s: %w", host, err)
		}
	}

	return nil
}
