package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine2/internal/config"
)

const (
	DefaultVault      = "Employee"
	ItemTitlePrefix   = "Devbox backup: "
	DevboxItemName    = "devbox-config"
)

func Create(filterName string) error {
	client := NewOpClient(DefaultVault)

	if err := client.CheckInstalled(); err != nil {
		return err
	}
	if err := client.CheckSignedIn(); err != nil {
		return err
	}

	backups, err := buildEffectiveBackups()
	if err != nil {
		return err
	}

	processed := 0
	for _, backup := range backups {
		if filterName != "" && backup.Name != filterName {
			continue
		}

		processed++
		if err := createBackup(client, backup); err != nil {
			return err
		}
	}

	if filterName != "" && processed == 0 {
		return fmt.Errorf("no backup entry with name '%s' found", filterName)
	}

	fmt.Println("Backup completed.")
	return nil
}

func Restore(filterName string) error {
	client := NewOpClient(DefaultVault)

	if err := client.CheckInstalled(); err != nil {
		return err
	}
	if err := client.CheckSignedIn(); err != nil {
		return err
	}

	backups, err := buildEffectiveBackupsForRestore()
	if err != nil {
		return err
	}

	processed := 0
	for _, backup := range backups {
		if filterName != "" && backup.Name != filterName {
			continue
		}

		processed++
		if err := restoreBackup(client, backup); err != nil {
			return err
		}
	}

	if filterName != "" && processed == 0 {
		return fmt.Errorf("no backup entry with name '%s' found", filterName)
	}

	fmt.Println("Restore completed.")
	return nil
}

func buildEffectiveBackups() ([]config.BackupConfig, error) {
	devboxPath, err := config.DevboxPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(devboxPath); err != nil {
		return nil, fmt.Errorf("devbox.json not found at: %s", devboxPath)
	}

	cfg, err := config.LoadDevbox()
	if err != nil {
		return nil, err
	}

	// Add synthetic devbox-config entry
	backups := []config.BackupConfig{
		{
			Name:  DevboxItemName,
			Path:  devboxPath,
			Vault: DefaultVault,
			Type:  "file",
		},
	}

	backups = append(backups, cfg.Backups...)
	return backups, nil
}

func buildEffectiveBackupsForRestore() ([]config.BackupConfig, error) {
	devboxPath, err := config.DevboxPath()
	if err != nil {
		return nil, err
	}

	// Always include devbox-config
	backups := []config.BackupConfig{
		{
			Name:  DevboxItemName,
			Path:  devboxPath,
			Vault: DefaultVault,
			Type:  "file",
		},
	}

	// Try to load existing config, but don't fail if it doesn't exist
	if _, err := os.Stat(devboxPath); err == nil {
		cfg, err := config.LoadDevbox()
		if err == nil {
			backups = append(backups, cfg.Backups...)
		}
	}

	return backups, nil
}

func createBackup(client *OpClient, backup config.BackupConfig) error {
	path, err := config.ExpandPath(backup.Path)
	if err != nil {
		return fmt.Errorf("failed to expand path for backup '%s': %w", backup.Name, err)
	}
	title := ItemTitlePrefix + backup.Name
	vault := backup.Vault
	if vault == "" {
		vault = DefaultVault
	}

	backupType := backup.Type
	if backupType == "" {
		backupType = "file"
	}

	var uploadPath string
	var cleanup func()

	if backupType == "directory" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Skipping '%s': directory not found: %s\n", backup.Name, path)
			return nil
		}

		fmt.Printf("Backing up directory '%s' as '%s' to vault '%s'...\n", path, title, vault)

		tmpFile, err := os.CreateTemp("", backup.Name+".tar.gz.*")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()

		if err := CreateTarGz(path, tmpPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("failed to create archive: %w", err)
		}

		uploadPath = tmpPath
		cleanup = func() { os.Remove(tmpPath) }
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Skipping '%s': file not found: %s\n", backup.Name, path)
			return nil
		}

		fmt.Printf("Backing up file '%s' as '%s' to vault '%s'...\n", path, title, vault)
		uploadPath = path
		cleanup = func() {}
	}

	defer cleanup()

	// Delete existing item if it exists
	existingID, err := client.FindItemByTitle(title)
	if err != nil {
		return fmt.Errorf("failed to check for existing item: %w", err)
	}

	if existingID != "" {
		if err := client.DeleteItem(existingID); err != nil {
			return fmt.Errorf("failed to delete existing item: %w", err)
		}
	}

	if err := client.CreateDocument(uploadPath, title); err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

func restoreBackup(client *OpClient, backup config.BackupConfig) error {
	path, err := config.ExpandPath(backup.Path)
	if err != nil {
		return fmt.Errorf("failed to expand path for backup '%s': %w", backup.Name, err)
	}
	title := ItemTitlePrefix + backup.Name
	vault := backup.Vault
	if vault == "" {
		vault = DefaultVault
	}

	backupType := backup.Type
	if backupType == "" {
		backupType = "file"
	}

	itemID, err := client.FindItemByTitle(title)
	if err != nil {
		return fmt.Errorf("failed to find item: %w", err)
	}

	if itemID == "" {
		fmt.Printf("Skipping '%s': no 1Password item titled '%s' in vault '%s'\n", backup.Name, title, vault)
		return nil
	}

	fmt.Printf("Restoring '%s' from '%s' in vault '%s'...\n", backup.Name, title, vault)

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if backupType == "directory" {
		tmpFile, err := os.CreateTemp("", backup.Name+".tar.gz.*")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		if err := client.GetDocument(itemID, tmpPath); err != nil {
			return fmt.Errorf("failed to download document: %w", err)
		}

		if err := ExtractTarGz(tmpPath, parentDir); err != nil {
			return fmt.Errorf("failed to extract archive: %w", err)
		}
	} else {
		if err := client.GetDocument(itemID, path); err != nil {
			return fmt.Errorf("failed to download document: %w", err)
		}
	}

	return nil
}
