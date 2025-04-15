package core

import (
	"Gozelle/internal/db"
	"math"
	"time"
)

// WeighFrecency calculates the frecency score of a Directory instance based on its LastVisit time and Score.
func WeighFrecency(dir *db.Directory) float64 {
	const halfLifeDecay = 0.693

	lastVisitTime := time.Unix(int64(dir.LastVisit), 0) // Convert age to int64 then to time.Time

	elapsedTime := time.Since(lastVisitTime)

	elapsedDays := elapsedTime.Hours() / 24.0

	decayFactor := math.Exp(-halfLifeDecay * elapsedDays)

	// Return the frecency score
	return float64(dir.Score) * decayFactor
}
