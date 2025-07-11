package db

import (
	"fmt"
	"os"
	"testing"
)

// add dummy data to the store
func (dm *DirectoryManager) dummyData() {
	dm.Add("/test/path1")
	dm.Add("/test/path2")
	dm.Add("/test/path3")
	dm.Add("/test/path4")

	dm.Dirty = true

	dm.Save()
}

func TestOpen(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	_, err = dm.Open(dm.FilePath)
	if err != nil {
		t.Fatalf("failed to open test store: %v", err)
	}
}

func TestAdd(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

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
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

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
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

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
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()
	dm2, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm2.DeleteTestStore()

	dm.dummyData()
	dm2.dummyData()

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
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dir := "/test/path"
	err = dm.Add(dir)
	if err != nil {
		t.Fatalf("failed to add directory: %v", err)
	}

	entry, err := dm.Get(dir)
	if err != nil {
		t.Fatalf("failed to get directory: %v", err)
	}

	if entry.Path != dir {
		t.Fatalf("expected path %s, got %s", dir, entry.Path)
	}
	if entry.Score != 1 {
		t.Fatalf("expected score 1, got %f", entry.Score)
	}
	if entry.LastVisit == 0 {
		t.Fatal("expected non-zero last visit time")
	}
}

func TestAll(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dm.dummyData()

	entries, err := dm.All()
	if err != nil {
		t.Fatalf("failed to get all directories: %v", err)
	}

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	for i, entry := range entries {
		if entry.Path != fmt.Sprintf("/test/path%d", i+1) {
			t.Fatalf("expected path /test/path%d, got %s", i+1, entry.Path)
		}
		if entry.Score != 1 {
			t.Fatalf("expected score 1, got %f", entry.Score)
		}
		if entry.LastVisit == 0 {
			t.Fatal("expected non-zero last visit time")
		}
	}
}

func TestSortByDirectory(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dm.dummyData()

	dm.Add("/stest/path")

	err = dm.SortByDirectory()
	if err != nil {
		t.Fatalf("failed to sort directories: %v", err)
	}

	if len(dm.Entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(dm.Entries))
	}

	for i := 0; i < len(dm.Entries)-1; i++ {
		if dm.Entries[i].Path > dm.Entries[i+1].Path {
			t.Fatalf("expected sorted order, got %s > %s", dm.Entries[i].Path, dm.Entries[i+1].Path)
		}
	}
}

func TestDedup(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dm.dummyData()

	// Add duplicate entries
	dm.Add("/test/path1")
	dm.Add("/test/path1")
	dm.Add("/test/path1")

	dm.Add("/test/path2")
	dm.Add("/test/path2")

	dm.Add("/test/path3")

	err = dm.Dedup()
	if err != nil {
		t.Fatalf("failed to deduplicate directories: %v", err)
	}

	if len(dm.Entries) != 4 {
		t.Fatalf("expected 4 entries after deduplication, got %d", len(dm.Entries))
	}

	if dm.Entries[0].Path != "/test/path1" {
		t.Fatalf("expected path /test/path1, got %s", dm.Entries[0].Path)
	}
	if dm.Entries[1].Path != "/test/path2" {
		t.Fatalf("expected path /test/path2, got %s", dm.Entries[1].Path)
	}
	if dm.Entries[2].Path != "/test/path3" {
		t.Fatalf("expected path /test/path3, got %s", dm.Entries[2].Path)
	}
	if dm.Entries[3].Path != "/test/path4" {
		t.Fatalf("expected path /test/path4, got %s", dm.Entries[3].Path)
	}
	if dm.Entries[0].Score != 4 {
		t.Fatalf("expected score 4, got %f", dm.Entries[0].Score)
	}
	if dm.Entries[1].Score != 3 {
		t.Fatalf("expected score 3, got %f", dm.Entries[1].Score)
	}
	if dm.Entries[2].Score != 2 {
		t.Fatalf("expected score 2, got %f", dm.Entries[2].Score)
	}
	if dm.Entries[3].Score != 1 {
		t.Fatalf("expected score 1, got %f", dm.Entries[3].Score)
	}

	for i := range dm.Entries {
		if dm.Entries[i].LastVisit == 0 {
			t.Fatal("expected non-zero last visit time")
		}
	}
}

func TestAddUpdate(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dir := "/test/path5"
	err = dm.AddUpdate(dir)
	if err != nil {
		t.Fatalf("failed to update directory: %v", err)
	}

	entry, err := dm.Get(dir)
	if err != nil {
		t.Fatalf("failed to get directory: %v", err)
	}

	if entry.Path != dir {
		t.Fatalf("expected path %s, got %s", dir, entry.Path)
	}

	if entry.Score != 1 {
		t.Fatalf("expected score 1, got %f", entry.Score)
	}
}

func TestSwapRemove(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dm.dummyData()

	if len(dm.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(dm.Entries))
	}

	dm.SwapRemoveIDX(2)

	if len(dm.Entries) != 3 {
		t.Fatalf("expected 3 entries after swap remove, got %d", len(dm.Entries))
	}

	if dm.Entries[2].Path == "/test/path3" {
		t.Fatal("expected path /test/path3 to be removed")
	}
}

func TestRemove(t *testing.T) {
	dm, err := CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dir := "/test/path"
	err = dm.Add(dir)
	if err != nil {
		t.Fatalf("failed to add directory: %v", err)
	}

	err = dm.Remove(dir)
	if err != nil {
		t.Fatalf("failed to remove directory: %v", err)
	}

	_, err = dm.Get(dir)
	if err == nil {
		t.Fatalf("expected error getting removed directory, got nil")
	}
}

// func TestDetermineFilthy(t *testing.T) {
// 	dm, err := CreateTestStore()
// 	if err != nil {
// 		t.Fatalf("failed to create test store: %v", err)
// 	}
// 	defer dm.DeleteTestStore()
//
// 	dm.dummyData()
// 	dm.Entries = append(dm.Entries, &Directory{Path: "/test/path5", Score: 1, LastVisit: 0})
//
// 	dm.DetermineFilthy()
// 	if !dm.Dirty {
// 		t.Fatal("expected dirty flag to be true")
// 	}
//
// 	dm.Save()
// 	dm.DetermineFilthy()
// 	if dm.Dirty {
// 		t.Fatal("expected dirty flag to be false")
// 	}
// }
