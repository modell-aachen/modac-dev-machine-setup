package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "machine",
	Short: "Modac development machine provisioner",
	Long: `Machine is a CLI tool for provisioning and managing
development environments for Modac projects.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	rootCmd.AddCommand(provisionCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(editConfigCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
