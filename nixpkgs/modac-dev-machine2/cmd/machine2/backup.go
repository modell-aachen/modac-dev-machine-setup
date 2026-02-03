package main

import (
	"github.com/spf13/cobra"

	"github.com/modell-aachen/machine2/internal/backup"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create or restore backups to/from 1Password",
	Long: `Create or restore backup files to/from 1Password configured in devbox.json.

The backup command manages files and directories by storing them in 1Password.
Configuration is read from devbox.json's "backups" array.

Examples:
  machine backup create                    # Backup everything
  machine backup restore                   # Restore everything
  machine backup create devbox-config      # Backup only devbox.json
  machine backup restore ssh-config        # Restore only ssh-config`,
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create backup to 1Password",
	Long:  "Backup devbox.json and all configured items, or a specific item by name",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filterName := ""
		if len(args) > 0 {
			filterName = args[0]
		}
		return backup.Create(filterName)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [name]",
	Short: "Restore backup from 1Password",
	Long:  "Restore devbox.json and all configured items, or a specific item by name",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filterName := ""
		if len(args) > 0 {
			filterName = args[0]
		}
		return backup.Restore(filterName)
	},
}

func init() {
	backupCmd.AddCommand(createCmd)
	backupCmd.AddCommand(restoreCmd)
}
