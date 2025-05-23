package core

import "github.com/atliod/gozelle/internal/db"

func Remove(path string) error {
	database, err := db.NewDirectoryManager()
	if err != nil {
		panic(err)
	}

	// find the path in the database
	for i, dir := range database.Entries {
		if dir.Path == path {
			database.SwapRemoveIDX(i)
			break
		}
	}
	return database.Save()
}
