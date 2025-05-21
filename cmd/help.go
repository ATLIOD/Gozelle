package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help information for Gozelle",
	Long: `Gozelle is a lightning-fast, minimal directory-jumping tool written in Go.
It helps you jump to frequently used directories with just a keyword,
powered by frecency scoring, fuzzy matching, and shell integration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
Gozelle - Lightning-fast directory jumping tool

USAGE:
  gozelle <command> [arguments]

COMMANDS:
  query <keyword> Show matching directories without jumping
  add <path>      Add a directory to the index
  remove <path>   Remove a directory from the index
  list           List all indexed directories
  help           Show this help message

EXAMPLES:
  # Initialize shell integration
  gozelle init <shell>  # e.g., bash, zsh, fish

  # Jump to a directory using a keyword
  gz <keyword> # Jump to the best match directory, e.g., 'gz projects' jumps to ~/Documents/projects
  # Show top match without jumping
  gozelle query <keyword>

  # Add a directory manually
  gozelle add /some/path/to/add

  # Remove a directory
  gozelle remove /some/path/to/remove

  # List all indexed directories
  gozelle list

Environment Variables:
  GOZELLE_ECHO           Whether the top match is printed before navigation or no(false or true)
  GOZELLE_DATA_DIR           The path where the data is stored (default: $XDG_DATA_HOME/gozelle/db.gob or ~/.local/share/gozelle/db.gob)

For more information, visit the project repository.
Github.com/atliod/gozelle
`)
	},
}
