package core

import (
	"gozelle/internal/db"
	"os"
)

func Prune() {
	database, err := db.NewDirectoryManagerWithPath(os.Getenv("GOZELLE_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	database.Dedup()
	database.Save()
}
