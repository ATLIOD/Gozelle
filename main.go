package main

import (
	"Gozelle/cmd"
	"fmt"
	"os"
)

func main() {
	keywords := os.Args[2:]

	switch os.Args[1] {
	case "query":
		// Call the query function
		result := cmd.QueryTop(keywords)
		fmt.Printf("Best match: %s\n", result.Path.Path)
		fmt.Printf("Frecency: %f\n", result.Frecency)
		fmt.Printf("Last visit: %d\n", result.Path.LastVisit)
		fmt.Printf("Score: %f\n", result.Path.Score)

	case "add":
		// Call the add function
		cmd.Add(os.Args[2])
	case "remove":
		// Call the remove function
	case "list":
		// Call the list function
	case "help":
		// Call the help function
	}
}
