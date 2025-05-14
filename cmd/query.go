package cmd

import (
	"fmt"
	"gozelle/internal/core"
	"gozelle/internal/db"
	"log"
	"runtime"
	"sync"
)

type ScoredMatch struct {
	Path     *db.Directory
	Frecency float64
}

// QueryTop searches for the best match in the directories based on keywords.
func QueryTop(keywords []string, path string) ScoredMatch {
	if len(keywords) == 0 {
		fmt.Print("./")
		return ScoredMatch{}
	}

	var database *db.DirectoryManager
	var err error

	// TODO: changes tests to use env variable so this comparison is not neeeded
	database, err = db.NewDirectoryManagerWithPath(path)
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
	if bestMatch.Path == nil {
		fmt.Print("./")
		return bestMatch
	}
	bestMatch.Path.UpdateLastVisit()
	bestMatch.Path.UpdateScore()
	database.Dirty = true
	if err := database.Save(); err != nil {
		log.Println("Error saving database:", err)
		panic(err)
	}
	fmt.Print(bestMatch.Path.Path)
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
