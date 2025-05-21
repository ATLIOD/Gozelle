package cmd

import (
	"log"

	"github.com/atliod/gozelle/internal/core"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a directory to the index",
	Long:  `Add a directory to the index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if err := core.Add(path); err != nil {
			log.Println("Error adding path:", err)
			return
		}
		core.Prune()
	},
}
