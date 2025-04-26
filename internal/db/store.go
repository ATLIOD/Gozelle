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
	// new methods
	Dedup() error
	SortByDirectory() error
	AddUpdate(dir *Directory) error
	Remove(dir *Directory) error
	DetermineFilthy() error
	SwapRemove(dir *Directory) error
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
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not get user home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	filePath := filepath.Join(dataDir, "Gozelle", "db.gob")
	// log.Println("Using file path:", filePath)
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
		log.Println("No data found, initializing empty directory manager.")
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
	return nil, fmt.Errorf("no directories found")
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
			dm.Entries = append(dm.Entries[:i+1], dm.Entries[i+2:]...)
			// don't increment i, check the new i+1 again
		} else {
			i++
		}
	}
	// set dirty flag to true
	dm.Dirty = true

	return nil
}

// SortByDirectory sorts the directories in the directory manager by their path
func (dm *DirectoryManager) SortByDirectory() error {
	quickSort(dm.Entries, 0, len(dm.Entries)-1)
	return nil
}

// AddUpdate adds a directory to the directory manager and then saves it to file
func (dm *DirectoryManager) AddUpdate(dir *Directory) error {
	dm.Add(dir.Path)
	dm.Save()
	return nil
}

// Remove removes a directory from the directory manager
func (dm *DirectoryManager) Remove(dir *Directory) error {
	return nil
}

// DetermineFilthy checks if the directory manager is dirty
func (dm *DirectoryManager) DetermineFilthy() error {
	return nil
}

// SwapRemove removes a directory from the directory manager and updates the file
func (dm *DirectoryManager) SwapRemove(dir *Directory) error {
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
