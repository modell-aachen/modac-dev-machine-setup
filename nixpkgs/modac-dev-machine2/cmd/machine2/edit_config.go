package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/modell-aachen/machine2/internal/config"
)

var editConfigCmd = &cobra.Command{
	Use:   "edit-config",
	Short: "Edit devbox.json configuration",
	Long:  "Open devbox.json in your editor (respects $EDITOR environment variable)",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.DevboxPath()
		if err != nil {
			return fmt.Errorf("failed to get devbox config path: %w", err)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi" // fallback to vi
		}

		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
	},
}
