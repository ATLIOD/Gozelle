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
	Add(path string) error
	Get(path string) (*Directory, error)
	All() ([]Directory, error)
	Save() error
	Dedup() error
	SortByDirectory() error
	AddUpdate(dir string) error
	Remove(dir *Directory) error
	DetermineFilthy() error
	SwapRemove(idx int) error
}

type DirectoryManager struct {
	Entries  []*Directory
	FilePath string
	Dirty    bool
	raw      []byte
	mu       sync.RWMutex
}

// NewDirectoryManager creates a new GobStore instance by accessing  reading in data from the given filepath.
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

	dm.raw = *rawgob

	err = dm.Decode(rawgob)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func NewDirectoryManager() (*DirectoryManager, error) {
	filePath := os.Getenv("GOZELLE_DATA_DIR")

	dm := &DirectoryManager{
		FilePath: filePath,
		Entries:  []*Directory{},
		Dirty:    false,
	}

	rawgob, err := dm.Open(filePath)
	if err != nil {
		return nil, err
	}

	dm.raw = *rawgob

	err = dm.Decode(rawgob)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func (dm *DirectoryManager) Open(filePath string) (*[]byte, error) {
	// check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// check directory exists
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println("Error creating directory:", filePath)
			return nil, fmt.Errorf("failed to create directories: %w", err)
		}

		// create an empty file
		if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
	}

	// read the file's contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &data, nil
}

// Decode decodes the data from the byte slice into the DirectoryManager's Entries field.
func (dm *DirectoryManager) Decode(data *[]byte) error {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if data == nil || len(*data) == 0 {
		dm.Entries = []*Directory{}
		// log.Println("No data found, initializing empty directory manager.")
		return nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(*data))

	// decode datainto directory slice
	var decodedEntries []*Directory
	err := decoder.Decode(&decodedEntries)
	if err != nil {
		return errors.New("failed to decode data: " + err.Error())
	}

	dm.Entries = decodedEntries

	// log.Println("Decoded entries:", dm.Entries)
	return nil
}

// Encode encodes the DirectoryManager's Entries field into a byte slice.
func (dm *DirectoryManager) Encode(entries []*Directory) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Encode the slice of Directory pointers into the buffer
	err := encoder.Encode(entries)
	if err != nil {
		return nil, errors.New("failed to encode Entries: " + err.Error())
	}

	return buf.Bytes(), nil
}

// add new directory to directory manager and updates the file
func (dm *DirectoryManager) Add(path string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dir := NewDirectory(path)
	dm.Entries = append(dm.Entries, dir)
	dm.Dirty = true
	if err := dm.Save(); err != nil {
		return fmt.Errorf("failed to save directory manager: %w", err)
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
	if len(dm.Entries) > 0 {
		return dm.Entries, nil
	}
	// if no entries are found, return an empty slice
	return []*Directory{}, fmt.Errorf("no directories found")
}

// saves the directory manager to a file
func (dm *DirectoryManager) Save() error {
	if !dm.Dirty {
		return nil
	}
	encodedData, err := dm.Encode(dm.Entries)
	if err != nil {
		return fmt.Errorf("failed to encode directory manager: %w", err)
	}
	err = os.WriteFile(dm.FilePath, encodedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	dm.Dirty = false
	// update the raw data after saving
	dm.raw, _ = dm.Encode(dm.Entries)

	return nil
}

// Dedup removes duplicate directories from the directory manager
func (dm *DirectoryManager) Dedup() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.SortByDirectory()

	for i := 0; i < len(dm.Entries)-1; {
		if dm.Entries[i].Path == dm.Entries[i+1].Path {
			dm.Entries[i].Score += dm.Entries[i+1].Score
			if dm.Entries[i].LastVisit < dm.Entries[i+1].LastVisit {
				dm.Entries[i].LastVisit = dm.Entries[i+1].LastVisit
			}
			// remove duplicate entry
			dm.RemoveIDX(i + 1)
			if !dm.Dirty {
				dm.Dirty = true
			}
			// don't increment i, check the new i+1 again
		} else {
			i++
		}
	}
	return nil
}

// SortByDirectory sorts the directories in the directory manager by their path
func (dm *DirectoryManager) SortByDirectory() error {
	quickSort(dm.Entries, 0, len(dm.Entries)-1)
	return nil
}

// AddUpdate adds a directory to the directory manager and then saves it to file
func (dm *DirectoryManager) AddUpdate(dir string) error {
	dm.Add(dir)
	dm.Save()
	return nil
}

func (dm *DirectoryManager) RemoveIDX(idx int) error {
	if idx < 0 || idx >= len(dm.Entries) {
		return fmt.Errorf("index out of range: %d", idx)
	}

	dm.Entries = append(dm.Entries[:idx], dm.Entries[idx+1:]...)
	return nil
}

// Remove removes a directory from the directory manager
func (dm *DirectoryManager) Remove(dir string) error {
	dm.SortByDirectory()

	idx := binarySearch(dm.Entries, dir)

	if idx == -1 {
		return fmt.Errorf("directory not found: %s", dir)
	}

	dm.Entries = append(dm.Entries[:idx], dm.Entries[idx+1:]...)
	return nil
}

// DetermineFilthy checks if the directory manager is dirty
func (dm *DirectoryManager) DetermineFilthy() error {
	current, err := dm.Encode(dm.Entries)
	if err != nil {
		return fmt.Errorf("failed to encode directory manager: %w", err)
	}
	if bytes.Equal(current, dm.raw) {
		dm.Dirty = false
		return nil
	}
	dm.Dirty = true
	return nil
}

// SwapRemove removes a directory from the directory manager and updates the file
// useful because it makes removal 0(1) instead of O(n)
func (dm *DirectoryManager) SwapRemoveIDX(idx int) error {
	if idx < 0 || idx >= len(dm.Entries) {
		return fmt.Errorf("index out of range: %d", idx)
	}
	if idx == -1 {
		return fmt.Errorf("directory not found")
	}
	// Swap the entry with the last entry and then remove the last entry
	dm.Entries[idx], dm.Entries[len(dm.Entries)-1] = dm.Entries[len(dm.Entries)-1], dm.Entries[idx]
	dm.Entries = dm.Entries[:len(dm.Entries)-1]
	dm.Dirty = true
	return nil
}

func (dm *DirectoryManager) SwapRemove(dir string) error {
	dm.SortByDirectory()

	idx := binarySearch(dm.Entries, dir)

	if idx == -1 {
		return fmt.Errorf("directory not found: %s", dir)
	}

	// Swap the entry with the last entry and then remove the last entry
	dm.Entries[idx], dm.Entries[len(dm.Entries)-1] = dm.Entries[len(dm.Entries)-1], dm.Entries[idx]
	dm.Entries = dm.Entries[:len(dm.Entries)-1]
	dm.Dirty = true
	return nil
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
