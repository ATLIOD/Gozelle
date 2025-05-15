package core

import (
	"fmt"
	"os"
	"os/exec"
)

func CheckFzfInstalled() {
	if _, err := exec.LookPath("fzf"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: fzf is not installed. Install it to use interactive mode (-i).")
		os.Exit(1)
	}
}
