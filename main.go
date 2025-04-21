package main

import (
	"Gozelle/cmd"
	"os"
)

func main() {
	keywords := os.Args[2:]

	switch os.Args[1] {
	case "query":
		// Call the query function
		cmd.QueryTop(keywords)
	case "add":
		// Call the add function
	case "remove":
		// Call the remove function
	case "list":
		// Call the list function
	case "help":
		// Call the help function
	}
}
