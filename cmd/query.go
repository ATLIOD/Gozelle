package cmd

import (
	"Gozelle/internal/core"
	"Gozelle/internal/db"
	"runtime"
	"sync"
	"time"
)

type ScoredMatch struct {
	Path     *db.Directory
	Frecency float64
}

// QueryTop searches for the best match in the directories based on keywords.
func QueryTop(keywords []string) ScoredMatch {
	database, err := db.NewDirectoryManager()
	if err != nil {
		panic(err)
	}
	jobs := make(chan *db.Directory)
	results := make(chan ScoredMatch)
	var wg sync.WaitGroup

	// start workers
	numWorkers := runtime.NumCPU()
	for range numWorkers {
		wg.Add(1)
		go worker(jobs, results, keywords, &wg)
	}

	// feed jobs
	go func() {
		for _, dir := range database.Entries {
			jobs <- dir
		}
		close(jobs)
	}()

	// close results when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// find best match
	var bestMatch ScoredMatch
	for match := range results {
		if match.Frecency > bestMatch.Frecency {
			bestMatch = match
		}
	}
	bestMatch.Path.LastVisit = db.Age(time.Now().Unix())
	bestMatch.Path.Score += 1
	return bestMatch
}

func worker(jobs <-chan *db.Directory, results chan<- ScoredMatch, keywords []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for dir := range jobs {
		if core.MatchByKeywords(dir.Path, keywords) {
			score := core.WeighFrecency(dir)
			results <- ScoredMatch{Path: dir, Frecency: score}
		}
	}
}
