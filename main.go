package main

import (
	"fmt"
	"os"

	"github.com/atliod/gozelle/cmd"
	"github.com/atliod/gozelle/internal/core"
)

func main() {
	core.SetConfig()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
