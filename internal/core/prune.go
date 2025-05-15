package core

import (
	"os"

	"github.com/atliod/gozelle/internal/db"
)

func Prune() {
	database, err := db.NewDirectoryManagerWithPath(os.Getenv("GOZELLE_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	database.Dedup()
	database.Save()
}
