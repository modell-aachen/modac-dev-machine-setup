package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/modell-aachen/machine/internal/provision"
)

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision this machine",
	Long:  "Provision a development machine with required tools and configurations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, err := cmd.Flags().GetString("filter")
		if err != nil {
			return fmt.Errorf("failed to get filter flag: %w", err)
		}

		skipInstall, err := cmd.Flags().GetBool("skip-install")
		if err != nil {
			return fmt.Errorf("failed to get skip-install flag: %w", err)
		}

		opts := &provision.Options{
			Filter:      filter,
			SkipInstall: skipInstall,
		}

		return provision.Execute(opts)
	},
}

var listModulesCmd = &cobra.Command{
	Use:   "list-modules",
	Short: "List all available modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		return provision.ListModules()
	},
}

func init() {
	provisionCmd.Flags().StringP("filter", "f", "", "Comma-separated list of modules to run (implies --skip-install)")
	provisionCmd.Flags().Bool("skip-install", false, "Skip running install, devbox shellenv, and op signin")
	provisionCmd.AddCommand(listModulesCmd)
}
