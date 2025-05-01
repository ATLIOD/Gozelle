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
  gozelle init bash

  # Jump to a directory using a keyword
  gz proj       # jumps to the best match (e.g., ~/projects)

  # Show top match without jumping
  gozelle query proj

  # Add a directory manually
  gozelle add /some/path/to/add

  # Remove a directory
  gozelle remove /some/path/to/remove

  # List all indexed directories
  gozelle list

FEATURES:
  - Frecency Scoring: Jump history is ranked by frequency and recency
  - Fuzzy Matching: Jump with just a keyword or part of a directory name
  - Smart Ranking: Most relevant paths surface first
  - Manual Add: Add directories to the index yourself
  - Query Mode: List matching directories without jumping
  - Compact Storage: Gob-encoded data stored locally
  - Shell Integration: Bash command-line hooks for seamless tracking

For more information, visit the project repository.
`)
	},
}
