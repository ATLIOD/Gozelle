package db

import (
	"fmt"
	"os"
	"testing"
)

func (dm *DirectoryManager) deleteTestStore() error {
	err := os.Remove(dm.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete test store: %w", err)
	}
	return nil
}

func createTestStore() (*DirectoryManager, error) {
	database, err := NewDirectoryManagerWithPath("./")
	if err != nil {
		return nil, fmt.Errorf("failed to create test store: %w", err)
	}
	return database, nil
}

// add dummy data to the store
func (dm *DirectoryManager) dummyData() {
	dm.Add("/test/path1")
	dm.Add("/test/path2")
	dm.Add("/test/path3")
	dm.Add("/test/path4")

	dm.Save()
}

func TestOpen(t *testing.T) {
	dm, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.deleteTestStore()

	_, err = dm.Open(dm.FilePath)
	if err != nil {
		t.Fatalf("failed to open test store: %v", err)
	}
}

func TestAdd(t *testing.T) {
	dm, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.deleteTestStore()

	dir := "/test/path"
	err = dm.Add(dir)
	if err != nil {
		t.Fatalf("failed to add directory: %v", err)
	}

	if len(dm.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(dm.Entries))
	}
	if dm.Entries[0].Path != dir {
		t.Fatalf("expected path %s, got %s", dir, dm.Entries[0].Path)
	}
	if dm.Entries[0].Score != 1 {
		t.Fatalf("expected score 1, got %f", dm.Entries[0].Score)
	}
}

func TestEncode(t *testing.T) {
	dm, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.deleteTestStore()

	dm.dummyData()

	encodedData, err := dm.Encode(dm.Entries)
	if err != nil {
		t.Fatalf("failed to encode data: %v", err)
	}

	if len(encodedData) == 0 {
		t.Fatal("encoded data is empty")
	}
}

func TestSave(t *testing.T) {
	dm, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.deleteTestStore()

	dm.dummyData()

	err = dm.Save()
	if err != nil {
		t.Fatalf("failed to save data: %v", err)
	}

	fileInfo, err := os.Stat(dm.FilePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatal("saved file is empty")
	}
}

func TestDecode(t *testing.T) {
	dm, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.deleteTestStore()
	dm2, err := createTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm2.deleteTestStore()

	dm.dummyData()

	encodedData, err := dm.Encode(dm.Entries)
	if err != nil {
		t.Fatalf("failed to encode data: %v", err)
	}

	err = dm.Decode(&encodedData)
	if err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if len(dm.Entries) != len(dm2.Entries) {
		t.Fatalf("expected %d entries, got %d", len(dm2.Entries), len(dm.Entries))
	}

	for i, entry := range dm.Entries {
		if entry.Path != dm2.Entries[i].Path {
			t.Fatalf("expected path %s, got %s", dm2.Entries[i].Path, entry.Path)
		}
		if entry.Score != dm2.Entries[i].Score {
			t.Fatalf("expected score %f, got %f", dm2.Entries[i].Score, entry.Score)
		}
	}
}

func TestGet(t *testing.T) {
}

func TestAll(t *testing.T) {
}

func TestDedup(t *testing.T) {
}

func TestSortByDirectory(t *testing.T) {
}

func TestAddUpdate(t *testing.T) {
}

func TestRemove(t *testing.T) {
}

func TestDetermineFilthy(t *testing.T) {
}

func TestSwapRemove(t *testing.T) {
}
