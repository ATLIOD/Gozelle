package core

import (
	"math"
	"time"

	"github.com/atliod/gozelle/internal/db"
)

// WeighFrecency calculates the frecency score of a Directory instance based on its LastVisit time and Score.
func WeighFrecency(dir *db.Directory) float64 {
	const halfLifeDecay float64 = 0.693

	lastVisitTime := time.Unix(int64(dir.LastVisit), 0) // Convert age to int64 then to time.Time

	elapsedTime := time.Since(lastVisitTime)

	decayFactor := math.Exp(-halfLifeDecay * elapsedTime.Hours())

	// Return the frecency score
	return float64(dir.Score) * decayFactor
}
