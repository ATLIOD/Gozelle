package db

import (
	"testing"
	"time"
)

func TestNewDirectory(t *testing.T) {
	path := "/path/to/directory"
	dir := NewDirectory(path)

	if dir.Path != path {
		t.Errorf("Expected Path to be %s, got %s", path, dir.Path)
	}

	if dir.LastVisit == 0 {
		t.Error("Expected LastVisit to be set to current time, got 0")
	}

	if dir.Score != 1 {
		t.Errorf("Expected Score to be 1, got %f", dir.Score)
	}
}

func TestUpdateLastVisit(t *testing.T) {
	path := "/path/to/directory"
	dir := NewDirectory(path)

	time.Sleep(1 * time.Second) // Ensure some time has passed

	originalLastVisit := dir.LastVisit
	dir.UpdateLastVisit()

	if dir.LastVisit == originalLastVisit {
		t.Error("Expected LastVisit to be updated, got same value")
	}
}

func TestUpdateScore(t *testing.T) {
	path := "/path/to/directory"
	dir := NewDirectory(path)

	originalScore := dir.Score
	dir.UpdateScore()

	if dir.Score <= originalScore {
		t.Errorf("Expected Score to be updated, got %f", dir.Score)
	}
}
