package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
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
}

type DirectoryManager struct {
	Entries  []*Directory
	FilePath string
	dirty    bool
	raw      []byte
	mu       sync.RWMutex
}

// NewDirectoryManager creates a new GobStore instance by accessing  reading in data from the given filepath.
func NewDirectoryManager(filePath string) (*DirectoryManager, error) {
	dm := &DirectoryManager{
		FilePath: filePath,
		Entries:  []*Directory{},
		dirty:    false,
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
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Read the file's contents
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
	decoder := gob.NewDecoder(bytes.NewReader(*data))

	// decode datainto directory slice
	var decodedEntries []*Directory
	err := decoder.Decode(&decodedEntries)
	if err != nil {
		return errors.New("failed to decode data: " + err.Error())
	}

	dm.Entries = decodedEntries

	return nil
}

// Encode encodes the DirectoryManager's Entries field into a byte slice.
func (dm *DirectoryManager) Encode(entries []*Directory) ([]byte, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
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
	dm.dirty = true
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
	dm.mu.RLock()
	defer dm.mu.Unlock()
	if !dm.dirty {
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
	dm.dirty = false
	// update the raw data after saving
	dm.raw, _ = dm.Encode(dm.Entries)

	return nil
}
