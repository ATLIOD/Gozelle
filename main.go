package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atliod/gozelle/cmd"
	"github.com/atliod/gozelle/internal/core"
)

func main() {
	core.SetConfig()

	if len(os.Args) < 2 {
		cmd.HelpCmd.Run(cmd.HelpCmd, []string{})
		os.Exit(1)
	}

	keywords := os.Args[2:]

	switch os.Args[1] {
	case "init":
		cmd.InitCmd.Run(cmd.InitCmd, keywords)
	case "query":
		// Call the query function
		result := cmd.QueryTop(keywords, os.Getenv("GOZELLE_DATA_DIR"))
		if result.Path == nil {
			log.Println("No match found")
			return
		}
		if os.Getenv("GOZELLE_ECHO") == "true" {
			log.Println("jumped to:", result.Path.Path)
		}
		core.Prune()
		// log.Println("Frecency:", result.Frecency)
		// log.Println("Last visit:", result.Path.LastVisit)
		// log.Println("Score:", result.Path.Score)
	case "add":
		// Call the add function
		cmd.Add(os.Args[2])
		core.Prune()

	case "remove":
		// Call the remove function
		cmd.Remove(os.Args[2])
		core.Prune()
	case "list":
		// Call the list function
		cmd.List()
		core.Prune()
	case "help":
		// Call the help function
		cmd.HelpCmd.Run(cmd.HelpCmd, keywords)
	case "interactive":
		core.CheckFzfInstalled()
		result, err := cmd.QueryInteractive(os.Getenv("GOZELLE_DATA_DIR"), false)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if os.Getenv("GOZELLE_ECHO") == "true" {
			log.Println("jumped to:", result)
		}

		core.Prune()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		cmd.HelpCmd.Run(cmd.HelpCmd, keywords)
		os.Exit(1)
	}
}
