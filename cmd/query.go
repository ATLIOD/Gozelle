package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/atliod/gozelle/internal/core"
	"github.com/atliod/gozelle/internal/db"
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

	database, err := db.NewDirectoryManagerWithPath(path)
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

func QueryInteractive(path string, multi bool) (string, error) {
	dm, err := db.NewDirectoryManagerWithPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to load directory manager: %w", err)
	}

	if len(dm.Entries) == 0 {
		return "", fmt.Errorf("no directories found in datastore")
	}

	lines := make([]string, len(dm.Entries))
	for i, dir := range dm.Entries {
		lines[i] = dir.Path
	}

	args := []string{"--ansi"}
	if multi {
		args = append(args, "--multi")
	}

	cmd := exec.Command("fzf", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		for _, line := range lines {
			fmt.Fprintln(stdin, line)
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("fzf exited with error or no selection: %w", err)
	}

	selectedLines := strings.Split(strings.TrimSpace(string(output)), "\n")

	var sel string
	// Print selected directories to stdout
	for _, sel = range selectedLines {
		fmt.Print(sel, "\n")
	}

	return sel, nil
}
