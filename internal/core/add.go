package core

import (
	"fmt"

	"github.com/atliod/gozelle/internal/db"
)

func Add(path string) error {
	database, err := db.NewDirectoryManager()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		panic(err)
	}

	err = database.Add(path)
	if err != nil {
		fmt.Println("Error adding path:", err)
		panic(err)
	}

	err = database.Save()
	if err != nil {
		fmt.Println("Error saving database:", err)
		panic(err)
	}
	fmt.Print("Path added successfully: ", path, "\n")
	return nil
}
