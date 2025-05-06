package db

import (
	"fmt"
	"os"
	"time"
)

func (dm *DirectoryManager) DeleteTestStore() error {
	err := os.Remove(dm.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete test store: %w", err)
	}
	return nil
}

func CreateTestStore() (*DirectoryManager, error) {
	database, err := NewDirectoryManagerWithPath("./test")
	if err != nil {
		return nil, fmt.Errorf("failed to create test store: %w", err)
	}
	return database, nil
}

// add dummy data to the store
func (dm *DirectoryManager) QueryDummyData() {
	dm.Add("/path1/test")
	dm.Add("/path2/test")
	dm.Add("/path3/test")
	dm.Add("/path4/test")
	dm.Add("/different/test")

	// increase the score of the first path
	dm.Entries[0].Score = 4
	dm.Entries[0].LastVisit = Age(time.Now().Unix())

	dm.Dirty = true

	dm.Save()
}
