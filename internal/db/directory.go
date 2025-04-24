package db

import (
	"time"
)

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
func (d *Directory) UpdateLastVisit() {
	d.LastVisit = Age(time.Now().Unix())
}

// UpdateScore updates the Score field of the Directory instance by multiplying it by 2.
func (d *Directory) UpdateScore() {
	d.Score *= 1.05
}
