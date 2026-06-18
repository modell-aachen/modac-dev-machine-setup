package nssdb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run initializes a per-user NSS database at ~/.pki/nssdb on Linux so that
// mkcert can make the development root CA available to Chrome and Firefox.
// It runs before the certificates module and is a no-op on macOS, where the
// browsers trust the system keychain that mkcert already populates.
func Run(out *output.Context, plat platform.Platform) error {
	if plat != platform.Ubuntu {
		out.Skipped("NSS database only required on Linux")
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nssDir := filepath.Join(homeDir, ".pki", "nssdb")

	// cert9.db is created when the database is initialized; its presence marks
	// an already-initialized NSS database.
	if util.FileExists(filepath.Join(nssDir, "cert9.db")) {
		out.Skipped("NSS database already exists")
		return nil
	}

	out.Step(fmt.Sprintf("Initializing NSS database in '%s'", nssDir))
	if err := os.MkdirAll(nssDir, 0755); err != nil {
		return fmt.Errorf("failed to create NSS database directory: %w", err)
	}

	cmd := exec.Command("certutil", "-d", "sql:"+nssDir, "-N", "--empty-password")
	if err := out.RunCmd(cmd); err != nil {
		return fmt.Errorf("failed to initialize NSS database: %w", err)
	}

	// When a root CA already exists, register it in the freshly created database
	// so existing certificates are trusted by the browsers. If no root CA exists
	// yet, the certificates module creates and installs it into this database.
	rootCADir := os.Getenv("QWIKI_DEVELOPMENT_ROOT_CA")
	if rootCADir != "" && util.FileExists(filepath.Join(rootCADir, "rootCA.pem")) {
		out.Step("Registering existing root CA in the new NSS database")
		install := exec.Command("mkcert", "-install")
		install.Env = append(os.Environ(), fmt.Sprintf("CAROOT=%s", rootCADir))
		if err := out.RunCmd(install); err != nil {
			return fmt.Errorf("failed to register root CA in NSS database: %w", err)
		}
	}

	return nil
}
