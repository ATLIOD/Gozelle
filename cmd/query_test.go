package cmd

import (
	"gozelle/internal/db"
	"testing"
)

// Test for queryTop
func TestQueryTop(t *testing.T) {
	dm, err := db.CreateTestStore()
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	defer dm.DeleteTestStore()

	dm.QueryDummyData()

	keywords := []string{"test"}
	bestMatch := QueryTop(keywords, dm.FilePath)
	if bestMatch.Path == nil {
		t.Fatalf("expected a match, got nil")
	}
	if bestMatch.Path.Path != "/path1/test" {
		t.Fatalf("expected path /path1/test, got %s", bestMatch.Path.Path)
	}
	if bestMatch.Frecency <= 0 {
		t.Fatalf("expected positive frecency, got %f", bestMatch.Frecency)
	}
	if bestMatch.Path.Score != 4*1.05 {
		t.Fatalf("expected score 4.2, got %f", bestMatch.Path.Score)
	}
	if bestMatch.Path.LastVisit == 0 {
		t.Fatalf("expected non-zero last visit, got %d", bestMatch.Path.LastVisit)
	}

	dm.Entries[1].Score = 5

	dm.Dirty = true
	dm.Save()

	bestMatchPath2 := QueryTop(keywords, dm.FilePath)
	if bestMatchPath2.Path == nil {
		t.Fatalf("expected a match, got nil")
	}
	if bestMatchPath2.Path.Path != "/path2/test" {
		t.Fatalf("expected path /path2/test, got %s", bestMatch.Path.Path)
	}
	if bestMatchPath2.Frecency <= 0 {
		t.Fatalf("expected positive frecency, got %f", bestMatch.Frecency)
	}
	if bestMatchPath2.Path.Score != 5*1.05 {
		t.Fatalf("expected score 5.25, got %f", bestMatch.Path.Score)
	}
	if bestMatchPath2.Path.LastVisit == 0 {
		t.Fatalf("expected non-zero last visit, got %d", bestMatch.Path.LastVisit)
	}

	keywords = []string{"different"}
	bestMatchDifferent := QueryTop(keywords, dm.FilePath)
	if bestMatch.Path == nil {
		t.Fatalf("expected a match, got nil")
	}
	if bestMatchDifferent.Path.Path != "/different/test" {
		t.Fatalf("expected path /different/test, got %s", bestMatch.Path.Path)
	}
	if bestMatchDifferent.Frecency <= 0 {
		t.Fatalf("expected positive frecency, got %f", bestMatch.Frecency)
	}
	if bestMatchDifferent.Path.Score != 1*1.05 {
		t.Fatalf("expected score 1.05, got %f", bestMatch.Path.Score)
	}
	if bestMatchDifferent.Path.LastVisit == 0 {
		t.Fatalf("expected non-zero last visit, got %d", bestMatch.Path.LastVisit)
	}

	keywords = []string{"nonexistent"}
	bestMatchNonexistent := QueryTop(keywords, dm.FilePath)
	if bestMatchNonexistent.Path != nil {
		t.Fatalf("expected no match, got %s", bestMatch.Path.Path)
	}
	if bestMatchNonexistent.Frecency != 0 {
		t.Fatalf("expected frecency 0, got %f", bestMatch.Frecency)
	}

	keywords = []string{}
	bestMatchEmpty := QueryTop(keywords, dm.FilePath)
	if bestMatchEmpty.Path != nil {
		t.Fatalf("expected no match, got %s", bestMatch.Path.Path)
	}
	if bestMatchEmpty.Frecency != 0 {
		t.Fatalf("expected frecency 0, got %f", bestMatch.Frecency)
	}
}
