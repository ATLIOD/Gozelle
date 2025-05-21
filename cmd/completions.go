package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CompletionsCmd = &cobra.Command{
	Use:   "completions [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for gozelle.

Example:
  gozelle completions bash > /etc/bash_completion.d/gozelle
  gozelle completions zsh > "${fpath[1]}/_gozelle"`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{
		"bash",
		"zsh",
		"fish",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			// Needed for Zsh completion to work
			fmt.Fprintln(os.Stdout, "autoload -U compinit; compinit")
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}
