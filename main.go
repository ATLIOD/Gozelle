package main

import (
	"fmt"
	"gozelle/cmd"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: Gozelle <command> [arguments]")
		fmt.Println("Commands: init, query, add, remove, list, help")
		os.Exit(1)
	}

	keywords := os.Args[2:]

	switch os.Args[1] {
	case "init":
		cmd.InitCmd.Run(cmd.InitCmd, keywords)
	case "query":
		// Call the query function
		result := cmd.QueryTop(keywords)
		log.Println("Best match:", result.Path.Path)
		log.Println("Frecency:", result.Frecency)
		log.Println("Last visit:", result.Path.LastVisit)
		log.Println("Score:", result.Path.Score)
	case "add":
		// Call the add function
		cmd.Add(os.Args[2])
	case "remove":
		// Call the remove function
	case "list":
		// Call the list function
	case "help":
		// Call the help function
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
