package cmd

import (
	"log"
	"os"

	"github.com/atliod/gozelle/internal/core"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query [keywords]",
	Short: "Query for directories",
	Long:  `Query for directories based on keywords.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keywords := args
		path := os.Getenv("GOZELLE_DATA_DIR")
		result := core.QueryTop(keywords, path)
		if result.Path == nil {
			log.Println("No match found")
			return
		}
		if os.Getenv("GOZELLE_ECHO") == "true" {
			log.Println("jumped to:", result.Path.Path)
		}
		core.Prune()
	},
}
