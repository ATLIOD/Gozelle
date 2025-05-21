package cmd

import (
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
			return RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return RootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return RootCmd.GenFishCompletion(os.Stdout, true)
		default:
			return nil
		}
	},
	Hidden: true,
}
