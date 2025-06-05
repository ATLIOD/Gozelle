package test_atomicdb

import (
	"os"
	"testing"
	"path/filepath"

	"github.com/stretchr/testify/require"
	"github.com/atliod/gozelle/internal/db"
)

func TestAtomicSave_NoCorruptionOnInterruption(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	// Add an entry and save
	dm.Add("/test/dir1")
	dm.Save()
	before, err := os.ReadFile(dbFile)
	require.NoError(t, err)
	require.NotEmpty(t, before)

	// Simulate interruption: write partial data to temp file, then crash
	tempFile := dbFile + ".tmp"
	partial := before[:len(before)/2]
	err = os.WriteFile(tempFile, partial, 0644)
	require.NoError(t, err)

	// Now call Save again (should atomically replace with full data)
	dm.Add("/test/dir2")
	dm.Save()
	after, err := os.ReadFile(dbFile)
	require.NoError(t, err)
	require.NotEmpty(t, after)
	// Should not be equal to partial, should be valid gob
	require.NotEqual(t, partial, after)

	// Try to decode
	dm2, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)
	entries, err := dm2.All()
	require.NoError(t, err)
	require.Len(t, entries, 2)
}

func TestAtomicSave_ConcurrentSaves(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "db.gob")
	dm, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)

	N := 10
	ch := make(chan struct{}, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			_ = dm.Add("/test/dir" + string(rune('A'+i)))
			_ = dm.Save()
			ch <- struct{}{}
		}(i)
	}
	for i := 0; i < N; i++ {
		<-ch
	}
	// Should not crash or corrupt
	dm2, err := db.NewDirectoryManagerWithPath(dbFile)
	require.NoError(t, err)
	entries, err := dm2.All()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 1)
}
