package cmd

import (
	"fmt"
	"gozelle/internal/core"
	"gozelle/internal/db"
)

func List() {
	database, err := db.NewDirectoryManager()
	if err != nil {
		panic(err)
	}

	// Print the list of directories
	for _, dir := range database.Entries {
		fmt.Println("Path: ", dir.Path, "|Frecency:", core.WeighFrecency(dir))
		fmt.Println("-------------------------------------------------------------------------------")
	}
}
