package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gozelle",
	Short: "Gozelle - smart directory jumper",
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(QueryCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(RemoveCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(InteractiveCmd)
	RootCmd.AddCommand(CompletionsCmd)
	// RootCmd.AddCommand(HelpCmd)
}
