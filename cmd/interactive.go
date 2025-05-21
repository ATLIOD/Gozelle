package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/atliod/gozelle/internal/core"
	"github.com/spf13/cobra"
)

var InteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode for Gozelle",
	Long: `Gozelle interactive mode allows you to jump to directories using fuzzy matching.
You can use this command to quickly navigate to your frequently used directories without needing to remember their exact paths.
This command is particularly useful for users who prefer a more visual and interactive way to select directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if fzf is installed
		core.CheckFzfInstalled()

		// Call the interactive query function
		result, err := core.QueryInteractive(os.Getenv("GOZELLE_DATA_DIR"), false)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if os.Getenv("GOZELLE_ECHO") == "true" {
			log.Println("jumped to:", result)
		}

		core.Prune()
	},
}
