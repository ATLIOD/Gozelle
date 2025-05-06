package core

import (
	"gozelle/internal/db"
)

// cleanstore loads in a datastore and calls the deduplicate function
func CleanStore() {
	database, err := db.NewDirectoryManager()
	if err != nil {
		panic(err)
	}

	database.Dedup()

	// NOTE: later will also call pruning

	if database.Dirty {
		err = database.Save()
		if err != nil {
			panic(err)
		}
	}
}
