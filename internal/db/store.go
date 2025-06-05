package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DataStore interface {
	Open(filePath string) (*[]byte, error)
	Decode(data *[]byte) error
	Encode(entries []*Directory) ([]byte, error)
	Add(path string) error // Adds to memory only
	AddAndSave(path string) error // Adds and persists
	Get(path string) (*Directory, error)
	All() ([]*Directory, error)
	Save() error
	Dedup() error
	SortByDirectory() error
	AddUpdate(dir string) error
	Remove(path string) error // Takes string path for consistency
	DetermineFilthy() error
	SwapRemoveIDX(idx int) error // Remove by index, O(1)
	SwapRemove(path string) error // Remove by path, O(1) if found
}

type DirectoryManager struct {
	Entries  []*Directory
	FilePath string
	Dirty    bool
	raw      []byte
	mu       sync.RWMutex
}

// NewDirectoryManager creates a new GobStore instance by accessing reading in data from the given filepath.
func NewDirectoryManagerWithPath(filePath string) (*DirectoryManager, error) {
	dm := &DirectoryManager{
		FilePath: filePath,
		Entries:  []*Directory{},
		Dirty:    false,
	}

	rawgob, err := dm.Open(filePath)
	if err != nil {
		return nil, err
	}
	dm.raw = *rawgob // Store the raw bytes for initial state

	err = dm.Decode(rawgob)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func NewDirectoryManager() (*DirectoryManager, error) {
	filePath := os.Getenv("GOZELLE_DATA_DIR")
	if filePath == "" {
		// Fallback or error handling if GOZELLE_DATA_DIR is not set
		// For now, let's assume it's set by core.SetConfig() before this is called.
		// If not, this could lead to issues. Consider returning an error or having a default.
		// However, core.SetConfig() ensures it's set.
	}
	return NewDirectoryManagerWithPath(filePath)
}

func (dm *DirectoryManager) Open(filePath string) (*[]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("[ERROR] Open: failed to create directories for %s: %v", dir, err)
			return nil, fmt.Errorf("failed to create directories: %w", err)
		}
		if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
			log.Printf("[ERROR] Open: failed to create empty file %s: %v", filePath, err)
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		log.Printf("[DEBUG] Open: Created new empty db file at %s", filePath)
		return &[]byte{}, nil // Return empty byte slice for new file
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("[ERROR] Open: failed to read file %s: %v", filePath, err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return &data, nil
}

// Decode decodes the data from the byte slice into the DirectoryManager's Entries field.
func (dm *DirectoryManager) Decode(data *[]byte) error {
	// dm.mu.RLock() // Decode is usually called during initialization before concurrent access
	// defer dm.mu.RUnlock() // or if called later, lock would be needed. Assuming init context for now.

	if data == nil || len(*data) == 0 {
		dm.Entries = []*Directory{}
		log.Println("[DEBUG] Decode: No data found, initializing empty directory manager.")
		return nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(*data))
	var decodedEntries []*Directory
	err := decoder.Decode(&decodedEntries)
	if err != nil {
		log.Printf("[ERROR] Decode: failed to decode data: %v", err)
		return fmt.Errorf("failed to decode data: %w", err)
	}
	dm.Entries = decodedEntries
	log.Printf("[DEBUG] Decode: Decoded %d entries.", len(dm.Entries))
	return nil
}

// Encode encodes the DirectoryManager's Entries field into a byte slice.
func (dm *DirectoryManager) Encode(entries []*Directory) ([]byte, error) {
	// dm.mu.RLock() // Encode reads Entries, so if called concurrently, RLock is needed.
	// defer dm.mu.RUnlock() // Assuming lock is managed by caller (e.g., saveInternal)
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(entries)
	if err != nil {
		log.Printf("[ERROR] Encode: failed to encode Entries: %v", err)
		return nil, fmt.Errorf("failed to encode Entries: %w", err)
	}
	return buf.Bytes(), nil
}

// Add adds a new directory to the directory manager in memory only.
// It marks the manager as Dirty but does NOT persist changes to disk.
// Call Save() to persist or use AddAndSave for immediate persistence.
// Add adds a new directory to the directory manager in memory only.
// It marks the manager as Dirty but does NOT persist changes to disk.
// Call Save() to persist or use AddAndSave for immediate persistence.
func (dm *DirectoryManager) Add(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	log.Printf("[DEBUG] Add: adding path '%s' to memory", path)
	dir := NewDirectory(path)
	dm.Entries = append(dm.Entries, dir)
	dm.Dirty = true
	return nil
}

// AddAndSave adds a new directory and immediately saves the directory manager to disk.
// This is the legacy behavior of Add prior to v0.2.0.
// AddAndSave adds a new directory and immediately saves the directory manager to disk.
func (dm *DirectoryManager) AddAndSave(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	log.Printf("[DEBUG] AddAndSave: adding path '%s' and saving", path)
	dir := NewDirectory(path)
	dm.Entries = append(dm.Entries, dir)
	dm.Dirty = true // Mark dirty before saveInternal checks it
	if err := dm.saveInternal(); err != nil {
		return fmt.Errorf("AddAndSave: %w", err)
	}
	return nil
}

// gets a directory from the directory manager
func (dm *DirectoryManager) Get(path string) (*Directory, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	for _, dir := range dm.Entries {
		if dir.Path == path {
			return dir, nil
		}
	}
	return nil, fmt.Errorf("directory not found: %s", path)
}

// gets all directories from the directory manager
func (dm *DirectoryManager) All() ([]*Directory, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	if len(dm.Entries) == 0 {
		// It's not an error to have no entries, just return an empty slice.
		return []*Directory{}, nil
	}
	// Return a copy to prevent external modification of the internal slice
	entriesCopy := make([]*Directory, len(dm.Entries))
	copy(entriesCopy, dm.Entries)
	return entriesCopy, nil
}

// saves the directory manager to a file
func (dm *DirectoryManager) Save() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	log.Println("[DEBUG] Save: acquiring lock (public Save called)")
	return dm.saveInternal()
}

// saveInternal is the actual save logic, assumes lock is already held by the caller.
func (dm *DirectoryManager) saveInternal() error {
	if !dm.Dirty {
		log.Println("[DEBUG] saveInternal: Not dirty, skipping save.")
		return nil
	}
	log.Println("[DEBUG] saveInternal: Encoding data.")
	encodedData, err := dm.Encode(dm.Entries) // Encode reads dm.Entries
	if err != nil {
		// Error already logged in Encode
		return fmt.Errorf("saveInternal encoding: %w", err)
	}

	tempFilePath := dm.FilePath + ".tmp"
	log.Printf("[DEBUG] saveInternal: Writing %d bytes to temporary file %s", len(encodedData), tempFilePath)
	err = os.WriteFile(tempFilePath, encodedData, 0644)
	if err != nil {
		log.Printf("[ERROR] saveInternal: failed to write to temp file %s: %v", tempFilePath, err)
		// Attempt to clean up the temporary file if rename fails, even if it's partially written
		_ = os.Remove(tempFilePath)
		return fmt.Errorf("failed to write to temporary file %s: %w", tempFilePath, err)
	}

	log.Printf("[DEBUG] saveInternal: Renaming temporary file %s to %s", tempFilePath, dm.FilePath)
	err = os.Rename(tempFilePath, dm.FilePath)
	if err != nil {
		log.Printf("[ERROR] saveInternal: failed to rename temp file %s to %s: %v", tempFilePath, dm.FilePath, err)
		_ = os.Remove(tempFilePath) // Clean up temp file if rename fails
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tempFilePath, dm.FilePath, err)
	}
	dm.Dirty = false
	dm.raw = make([]byte, len(encodedData)) // Update raw with the successfully saved data
	copy(dm.raw, encodedData)
	log.Println("[DEBUG] saveInternal: Save successful.")
	return nil
}

// Dedup removes duplicate directories from the directory manager
func (dm *DirectoryManager) Dedup() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	log.Println("[DEBUG] Dedup: Performing deduplication.")

	if len(dm.Entries) < 2 {
		log.Println("[DEBUG] Dedup: Not enough entries to dedup.")
		return nil // Nothing to dedup
	}
	dm.SortByDirectory() // SortByDirectory doesn't lock, called under dm.mu

	originalCount := len(dm.Entries)
	newEntries := make([]*Directory, 0, len(dm.Entries))
	if len(dm.Entries) > 0 {
		newEntries = append(newEntries, dm.Entries[0])
		for i := 1; i < len(dm.Entries); i++ {
			lastAdded := newEntries[len(newEntries)-1]
			current := dm.Entries[i]
			if lastAdded.Path == current.Path {
				lastAdded.Score += current.Score
				if lastAdded.LastVisit < current.LastVisit {
					lastAdded.LastVisit = current.LastVisit
				}
				dm.Dirty = true // Mark dirty if merging happened
			} else {
				newEntries = append(newEntries, current)
			}
		}
	}
	dm.Entries = newEntries

	if originalCount != len(dm.Entries) {
		dm.Dirty = true // Also dirty if count changed
		log.Printf("[DEBUG] Dedup: Reduced entries from %d to %d.", originalCount, len(dm.Entries))
	} else if dm.Dirty { // If only scores/visits updated
		log.Println("[DEBUG] Dedup: Updated scores/visits for duplicate paths.")
	} else {
		log.Println("[DEBUG] Dedup: No duplicates found or changes made.")
	}

	// Decide if Dedup should save. If it made changes, it should probably save.
	// if dm.Dirty {
	// 	log.Println("[DEBUG] Dedup: Changes made, saving.")
	// 	return dm.saveInternal()
	// }
	return nil
}

// SortByDirectory sorts the directories in the directory manager by their path
func (dm *DirectoryManager) SortByDirectory() error {
	// This method operates on dm.Entries directly.
	// It should be called when dm.mu is held if concurrent access is possible.
	// It does not modify Dirty status by itself.
	if len(dm.Entries) > 1 {
		quickSort(dm.Entries, 0, len(dm.Entries)-1)
	}
	return nil
}

// AddUpdate adds a directory to the directory manager and immediately persists it to file.
// This is equivalent to AddAndSave and is provided for API compatibility.
// AddUpdate adds a directory to the directory manager and immediately persists it to file.
func (dm *DirectoryManager) AddUpdate(dir string) error {
	log.Printf("[DEBUG] AddUpdate: adding path '%s'", dir)
	return dm.AddAndSave(dir)
}

// RemoveIDX removes an entry by index. Assumes lock is held by caller if needed
// and caller handles Dirty status and saving.
func (dm *DirectoryManager) RemoveIDX(idx int) error {
	if idx < 0 || idx >= len(dm.Entries) {
		return fmt.Errorf("RemoveIDX: index %d out of range for %d entries", idx, len(dm.Entries))
	}
	dm.Entries = append(dm.Entries[:idx], dm.Entries[idx+1:]...)
	// dm.Dirty = true // Caller should set Dirty and save
	return nil
}

// Remove removes a directory from the directory manager
// Remove searches for a directory by path and removes it if found.
// Persists changes immediately.
func (dm *DirectoryManager) Remove(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	log.Printf("[DEBUG] Remove: Attempting to remove path '%s'", path)

	dm.SortByDirectory() // Sort to use binarySearch
	idx := binarySearch(dm.Entries, path)

	if idx == -1 {
		log.Printf("[DEBUG] Remove: Path '%s' not found.", path)
		return fmt.Errorf("directory not found for removal: %s", path)
	}

	dm.Entries = append(dm.Entries[:idx], dm.Entries[idx+1:]...)
	dm.Dirty = true
	log.Printf("[DEBUG] Remove: Path '%s' removed, saving.", path)
	return dm.saveInternal()
}

// DetermineFilthy checks if the directory manager is dirty
func (dm *DirectoryManager) DetermineFilthy() error {
	dm.mu.RLock() // Need RLock to safely call Encode and access dm.raw
	defer dm.mu.RUnlock()

	current, err := dm.Encode(dm.Entries)
	if err != nil {
		return fmt.Errorf("DetermineFilthy encoding: %w", err)
	}
	if bytes.Equal(current, dm.raw) {
		// To be absolutely sure dm.Dirty is correct, we might reset it here.
		// However, dm.Dirty should reflect if an operation *intended* to make a change.
		// This check is more about whether the *current state* matches *last saved state*.
		// For simplicity, let's not modify dm.Dirty here, just report.
		// dm.Dirty = false
		return nil // Not filthy if current encoding matches raw bytes of last save
	}
	// dm.Dirty = true // No, don't set it here. This function just checks.
	return errors.New("DetermineFilthy: current state differs from last saved state (filthy)")
}

// SwapRemove removes a directory from the directory manager and updates the file
// useful because it makes removal 0(1) instead of O(n)
// SwapRemoveIDX removes a directory by index using O(1) swap and persists changes.
func (dm *DirectoryManager) SwapRemoveIDX(idx int) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	log.Printf("[DEBUG] SwapRemoveIDX: Attempting to remove index %d", idx)

	if idx < 0 || idx >= len(dm.Entries) {
		return fmt.Errorf("SwapRemoveIDX: index %d out of range for %d entries", idx, len(dm.Entries))
	}

	lastIdx := len(dm.Entries) - 1
	dm.Entries[idx], dm.Entries[lastIdx] = dm.Entries[lastIdx], dm.Entries[idx]
	dm.Entries = dm.Entries[:lastIdx]
	dm.Dirty = true
	log.Printf("[DEBUG] SwapRemoveIDX: Index %d removed, saving.", idx)
	return dm.saveInternal()
}

// SwapRemove searches for a directory by path, then uses SwapRemoveIDX for O(1) removal if found.
// Persists changes immediately.
func (dm *DirectoryManager) SwapRemove(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	log.Printf("[DEBUG] SwapRemove: Attempting to remove path '%s'", path)

	// This makes SwapRemove O(N log N) or O(N) due to search, not O(1) unless path is pre-indexed.
	// For true O(1) by path, a map[string]int would be needed to store indices.
	// Current implementation: find index first.
	idxToSwap := -1
	for i, dir := range dm.Entries {
		if dir.Path == path {
			idxToSwap = i
			break
		}
	}

	if idxToSwap == -1 {
		log.Printf("[DEBUG] SwapRemove: Path '%s' not found.", path)
		return fmt.Errorf("directory not found for swap-removal: %s", path)
	}
	
	// Now perform the O(1) removal part using the found index
	lastIdx := len(dm.Entries) - 1
	if idxToSwap <= lastIdx {
		dm.Entries[idxToSwap], dm.Entries[lastIdx] = dm.Entries[lastIdx], dm.Entries[idxToSwap]
		dm.Entries = dm.Entries[:lastIdx]
		dm.Dirty = true
		log.Printf("[DEBUG] SwapRemove: Path '%s' (index %d) removed, saving.", path, idxToSwap)
		return dm.saveInternal()
	}
	// Should not happen if idxToSwap was valid and list not empty
	return fmt.Errorf("SwapRemove: inconsistency for path %s", path)
}

func quickSort(arr []*Directory, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func partition(arr []*Directory, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j].Path < pivot.Path {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func binarySearch(arr []*Directory, target string) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := low + (high-low)/2
		midPath := arr[mid].Path

		if midPath == target {
			return mid
		} else if midPath < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1 // Not found
}
