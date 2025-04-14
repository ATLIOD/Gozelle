package models

import "time"

type Directory struct {
	Path          string
	LastVisit     time.Time
	FrecencyScore float32
}

// NewDirecotry creates a new Directory instance with the given path, current time as LastVisit, and a default frecency score.
func NewDirecotry(path *string) *Directory {
	return &Directory{
		*path,
		time.Now(), // NOTE: evaluate optimal default weight
		5,
	}
}
