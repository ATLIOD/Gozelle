package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell integration (bash only for now)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := "bash"
		if len(args) > 0 {
			shell = args[0]
		}

		switch shell {
		case "bash":
			fmt.Println(bashInitScript())
		default:
			fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
			os.Exit(1)
		}
	},
}

func bashInitScript() string {
	return strings.TrimSpace(`
# Gozelle Bash init script

__gozelle_oldpwd="$(pwd)"

__gozelle_hook() {
    local retval=$?
    local pwd_now="$(pwd)"
    if [[ "$__gozelle_oldpwd" != "$pwd_now" ]]; then
        __gozelle_oldpwd="$pwd_now"
        command gozelle add "$pwd_now" >/dev/null 2>&1
    fi
    return $retval
}

if [[ "$PROMPT_COMMAND" != *"__gozelle_hook"* ]]; then
    PROMPT_COMMAND="__gozelle_hook;${PROMPT_COMMAND#;}"
fi

gz() {
        target="$(command gozelle query "$@")" && cd "$target"
}
`)
}
