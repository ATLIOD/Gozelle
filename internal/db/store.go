package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type DataStore interface {
	Open(filePath string) (*[]byte, error)
	Decode(data *[]byte) (map[string]*Directory, error)
	Encode() ([]byte, error)
	Add(path string) error
	Get(path string) (*Directory, error)
	All() ([]Directory, error)
	Save() error
	Load() error
}

type directoryManager struct {
	Entries  map[string]*Directory
	FilePath string
	dirty    bool
	raw      []byte
}

// NewDirectoryManager creates a new GobStore instance by accessing  reading in data from the given filepath.
func NewDirectoryManager(filePath string) (*directoryManager, error) {
	dm := &directoryManager{
		FilePath: filePath,
		Entries:  make(map[string]*Directory),
		dirty:    false,
	}

	rawgob, err := dm.Open(filePath)
	if err != nil {
		return nil, err
	}

	dm.raw = *rawgob

	dm.Entries, err = dm.Decode(rawgob)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func (dm *directoryManager) Open(filePath string) (*[]byte, error) {
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

func (dm *directoryManager) Decode(data *[]byte) (map[string]*Directory, error) {
	// buffer to read into
	buf := bytes.NewReader(*data)
	decoder := gob.NewDecoder(buf)

	var decodedData map[string]*Directory

	// decode data
	if err := decoder.Decode(&decodedData); err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return decodedData, nil
}

func (dm *directoryManager) Encode(data map[string]*Directory) ([]byte, error) {
	// make buffer to put encoded data in
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// encode data
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	// return encoded data
	return buf.Bytes(), nil
}

func (dm *directoryManager) Add(path string) error {
	return nil
}

func (dm *directoryManager) Get(path string) (*Directory, error) {
	return nil, nil
}

func (dm *directoryManager) All() ([]Directory, error) {
	return nil, nil
}

func (dm *directoryManager) Save() error {
	return nil
}

func (dm *directoryManager) Load() error {
	return nil
}
