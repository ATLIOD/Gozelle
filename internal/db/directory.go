package db

import (
	"math"
	"sync"
	"time"
)

const halfLifeDecay = 0.693

type (
	Score float64
	Age   int64
)

type Directory struct {
	Path      string
	LastVisit Age
	Score     Score
}

// NewDirectory creates a new Directory instance with the given path, current time as LastVisit, and a default frecency score.
func NewDirectory(path string) *Directory {
	return &Directory{
		path,
		Age(time.Now().Unix()),
		Score(1), // NOTE: evaluate optimal default weight
	}
}

// UpdateLastVisit updates the LastVisit field of the Directory instance to the current time.
func UpdateLastVisit(dir *Directory) {
	dir.LastVisit = Age(time.Now().Unix())
}

// UpdateScore updates the Score field of the Directory instance by multiplying it by 2.
func UpdateScore(dir *Directory) {
	dir.Score *= 2
}

// WeighFrecency calculates the frecency score of a Directory instance based on its LastVisit time and Score.
func WeighFrecency(dir *Directory) float64 {
	lastVisitTime := time.Unix(int64(dir.LastVisit), 0) // Convert age to int64 then to time.Time

	elapsedTime := time.Since(lastVisitTime)

	elapsedDays := elapsedTime.Hours() / 24.0

	decayFactor := math.Exp(-halfLifeDecay * elapsedDays)

	// Return the frecency score
	return float64(dir.Score) * decayFactor
}

// DirectoryManager manages a list of Directory instances.
type DirectoryManager struct {
	dirs  []Directory
	dirty bool
	mu    sync.RWMutex // For concurrent access
}

// NewDirManager creates a new DirManager instance.
func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{
		dirs:  make([]Directory, 0),
		dirty: false,
	}
}

// AddDirectory adds a new Directory to the DirManager.
func (dm *DirectoryManager) AddDirectory(path string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.dirs = append(dm.dirs, *NewDirectory(path))
	dm.dirty = true
}
