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
	Load() error
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
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Encode the slice of Directory pointers into the buffer
	err := encoder.Encode(entries)
	if err != nil {
		return nil, errors.New("failed to encode Entries: " + err.Error())
	}

	return buf.Bytes(), nil
}

func (dm *DirectoryManager) Add(path string) error {
	return nil
}

func (dm *DirectoryManager) Get(path string) (*Directory, error) {
	return nil, nil
}

func (dm *DirectoryManager) All() ([]Directory, error) {
	return nil, nil
}

func (dm *DirectoryManager) Save() error {
	return nil
}

func (dm *DirectoryManager) Load() error {
	return nil
}
