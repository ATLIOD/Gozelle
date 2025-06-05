package test_atomicdb

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/atliod/gozelle/internal/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestAtomicSave_TempFileCleanup(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	// Add entry and save
	err = dm.Add("/test/dir1")
	require.NoError(t, err)
	err = dm.Save()
	require.NoError(t, err)

	// Verify temp file doesn't exist after successful save
	tempFile := dbFile + ".tmp"
	_, err = os.Stat(tempFile)
	require.True(t, os.IsNotExist(err), "Temp file should be cleaned up after successful save")

	// Verify main file exists and is valid
	_, err = os.Stat(dbFile)
	require.NoError(t, err)

	// Verify we can reload the data
	dm2, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)
	entries, err := dm2.All()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "/test/dir1", entries[0].Path)
}

func TestAtomicSave_PreservesOriginalOnTempFileFailure(t *testing.T) {
	dir := t.TempDir()

	// Use a subdirectory for the DB file
	dbDir := filepath.Join(dir, "dbdir")
	err := os.Mkdir(dbDir, 0755)
	require.NoError(t, err)
	dbFile := filepath.Join(dbDir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	// Create initial valid state
	err = dm.Add("/test/dir1")
	require.NoError(t, err)
	err = dm.Save()
	require.NoError(t, err)

	// Verify initial state
	originalData, err := os.ReadFile(dbFile)
	require.NoError(t, err)
	require.NotEmpty(t, originalData)

	// Try to add another entry
	err = dm.Add("/test/dir2")
	require.NoError(t, err)

	// Make the DB subdirectory read-only to simulate write failure
	err = os.Chmod(dbDir, 0555) // Read and execute only, no write
	require.NoError(t, err)

	// Ensure we restore permissions even if test fails
	defer func() {
		_ = os.Chmod(dbDir, 0755)
	}()

	// Try to save again - should fail due to directory permissions
	err = dm.Save()
	require.Error(t, err, "Save should fail when temp file can't be created")

	// Restore permissions before reading file
	err = os.Chmod(dbDir, 0755)
	require.NoError(t, err)

	// Original file should be unchanged
	currentData, err := os.ReadFile(dbFile)
	require.NoError(t, err)
	if !assert.Equal(t, originalData, currentData, "Original file should be preserved on save failure") {
		t.Logf("[DEBUG] originalData: %x", originalData)
		t.Logf("[DEBUG] currentData:  %x", currentData)
	}
}

func TestAtomicSave_ConsistentStateUnderConcurrency(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	const numGoroutines = 10
	const opsPerGoroutine = 5

	var wg sync.WaitGroup

	// Launch concurrent goroutines that add entries
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				path := fmt.Sprintf("/test/dir_%d_%d", id, j)
				err := dm.Add(path)
				if err != nil {
					t.Errorf("Add failed: %v", err)
					return
				}
				// Small delay to increase chance of interleaving
				time.Sleep(time.Millisecond)
				err = dm.Save()
				if err != nil {
					t.Errorf("Save failed: %v", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is consistent
	dm2, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)
	entries, err := dm2.All()
	require.NoError(t, err)

	// Should have exactly numGoroutines * opsPerGoroutine entries
	require.Len(t, entries, numGoroutines*opsPerGoroutine)

	// Verify all expected paths are present
	pathSet := make(map[string]bool)
	for _, entry := range entries {
		pathSet[entry.Path] = true
	}

	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < opsPerGoroutine; j++ {
			expectedPath := fmt.Sprintf("/test/dir_%d_%d", i, j)
			require.True(t, pathSet[expectedPath], "Expected path %s not found", expectedPath)
		}
	}
}

func TestAtomicSave_FileIntegrityAfterMultipleSaves(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	// Perform multiple saves and verify integrity each time
	for i := 0; i < 20; i++ {
		path := fmt.Sprintf("/test/dir%d", i)
		err = dm.Add(path)
		require.NoError(t, err)

		err = dm.Save()
		require.NoError(t, err)

		// Verify file can be read and decoded after each save
		dm2, err := db.NewDirectoryManagerWithPath(dbFile)
		require.NoError(t, err)
		entries, err := dm2.All()
		require.NoError(t, err)
		require.Len(t, entries, i+1)

		// Verify last added entry exists
		found := false
		for _, entry := range entries {
			if entry.Path == path {
				found = true
				break
			}
		}
		require.True(t, found, "Recently added path %s not found", path)
	}
}
