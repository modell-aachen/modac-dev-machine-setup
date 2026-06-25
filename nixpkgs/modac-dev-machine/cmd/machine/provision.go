package main

import (
	"fmt"
	"strings"

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

		opts := &provision.Options{
			Filter: filter,
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
	provisionCmd.Flags().StringP("filter", "f", "", "Comma-separated list of modules to run (tab-completable)")
	provisionCmd.AddCommand(listModulesCmd)

	provisionCmd.Long += "\n\nUse --filter to run only specific modules (comma-separated). Available modules:\n  " +
		strings.Join(provision.GetAllModuleNames(), "\n  ")

	_ = provisionCmd.RegisterFlagCompletionFunc("filter", completeFilter)
}

func completeFilter(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	segments := strings.Split(toComplete, ",")
	chosen := make(map[string]bool)
	for _, s := range segments[:len(segments)-1] {
		chosen[s] = true
	}

	prefix := strings.Join(segments[:len(segments)-1], ",")
	if prefix != "" {
		prefix += ","
	}

	var comps []string
	for _, name := range provision.GetAllModuleNames() {
		if chosen[name] {
			continue
		}
		if candidate := prefix + name; strings.HasPrefix(candidate, toComplete) {
			comps = append(comps, candidate)
		}
	}
	return comps, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
