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

func (dm *DirectoryManager) createTestStore() (*DirectoryManager, error) {
	database, err := NewDirectoryManagerWithPath("./")
	if err != nil {
		return nil, fmt.Errorf("failed to create test store: %w", err)
	}
	return database, nil
}

func TestOpen(t *testing.T) {
}

func TestDecode(t *testing.T) {
}

func TestEncode(t *testing.T) {
}

func TestAdd(t *testing.T) {
}

func TestGet(t *testing.T) {
}

func TestAll(t *testing.T) {
}

func TestSave(t *testing.T) {
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
