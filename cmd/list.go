package cmd

import (
	"log"

	"github.com/atliod/gozelle/internal/core"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all directories in the store",
	Long:  `List all directories in the store.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the list function from core package
		err := core.List()
		if err != nil {
			log.Println("Error listing directories:", err)
			return
		}
	},
}
