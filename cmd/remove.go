package cmd

import (
	"log"

	"github.com/atliod/gozelle/internal/core"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove [path]",
	Short: "Remove a directory from the store",
	Long:  `Remove a directory from the store.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if err := core.Remove(path); err != nil {
			log.Println("Error removing path:", err)
			return
		}
		core.Prune()
	},
}
