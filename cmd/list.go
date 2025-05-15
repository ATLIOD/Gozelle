package cmd

import (
	"fmt"

	"github.com/atliod/gozelle/internal/db"
)

func List() {
	database, err := db.NewDirectoryManager()
	if err != nil {
		panic(err)
	}

	// Print the list of directories
	for _, dir := range database.Entries {
		fmt.Println("Path: ", dir.Path, "|Frequency Score:", dir.Score)
		fmt.Println("-------------------------------------------------------------------------------")
	}
}
