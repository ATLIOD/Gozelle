package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atliod/gozelle/internal/core"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell integration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := "bash"
		if len(args) > 0 {
			shell = args[0]
		}

		switch shell {
		case "bash":
			go core.CleanStore()
			fmt.Println(bashInitScript())
		case "zsh":
			fmt.Println(zshInitScript())
		case "fish":
			fmt.Println(fishInitScript())
		default:
			fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
			os.Exit(1)
		}
	},
}

// eval "$(gozelle init bash)"
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
    if [ $# -eq 0 ]; then
        cd ~
    elif [ $# -eq 1 ] && [ "$1" = "-" ]; then
        cd "${OLDPWD}"
    elif [ $# -eq 1 ] && [ -d "$1" ]; then
        cd "$1"
    elif [ $# -eq 2 ] && [ "$1" = "--" ]; then
        cd "$2"
    else
        target="$(command gozelle query "$@")" && cd "$target"
    fi
}
gi() {
	target="$(command gozelle interactive "$@")" && cd "$target"
}
`)
}

// eval "$(gozelle init zsh)"
func zshInitScript() string {
	return strings.TrimSpace(`
# Gozelle Zsh init script (recommended chpwd hook version)

__gozelle_on_cd() {
    command gozelle add "$PWD" >/dev/null 2>&1
}

autoload -Uz add-zsh-hook
add-zsh-hook chpwd __gozelle_on_cd

gz() {
    if [ $# -eq 0 ]; then
        cd ~
    elif [ $# -eq 1 ] && [ "$1" = "-" ]; then
        cd "${OLDPWD}"
    elif [ $# -eq 1 ] && [ -d "$1" ]; then
        cd "$1"
    elif [ $# -eq 2 ] && [ "$1" = "--" ]; then
        cd "$2"
    else
        target="$(command gozelle query "$@")" && cd "$target"
    fi
}
gi() {
	target="$(command gozelle interactive "$@")" && cd "$target"
}
`)
}

// eval (gozelle init fish)
func fishInitScript() string {
	return strings.TrimSpace(`
# Gozelle Fish init script

function __gozelle_prompt_hook --on-event fish_prompt
    if not set -q __gozelle_oldpwd
        set -g __gozelle_oldpwd $PWD
    end

    if test "$__gozelle_oldpwd" != "$PWD"
        set -g __gozelle_oldpwd $PWD
        gozelle add "$PWD" > /dev/null 2>&1
    end
end

function gz
    if test (count $argv) -eq 0
        cd ~
    else if test (count $argv) -eq 1 -a "$argv[1]" = "-"
        cd "$OLDPWD"
    else if test (count $argv) -eq 1 -a -d "$argv[1]"
        cd "$argv[1]"
    else if test (count $argv) -eq 2 -a "$argv[1]" = "--"
        cd "$argv[2]"
    else
        set target (gozelle query $argv)
        if test -n "$target"
            cd "$target"
        end
    end
end

function gi
    set target (gozelle interactive $argv)
    if test -n "$target"
        cd "$target"
    end
end
`)
}
