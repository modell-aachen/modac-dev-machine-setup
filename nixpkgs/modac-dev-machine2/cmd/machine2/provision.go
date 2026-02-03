package main

import (
	"github.com/spf13/cobra"

	"github.com/modell-aachen/machine2/internal/provision"
)

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision this machine2",
	Long:  "Provision a development machine2 with required tools and configurations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, _ := cmd.Flags().GetString("filter")
		skipInstall, _ := cmd.Flags().GetBool("skip-install")

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
