package core

import "testing"

func TestMatchByKeywords(t *testing.T) {
	path := "/path/to/directory"
	keywords := []string{"path", "directory"}

	if !MatchByKeywords(path, keywords) {
		t.Errorf("Expected MatchByKeywords to return true for path: %s and keywords: %v", path, keywords)
	}

	keywords = []string{"not", "in", "path"}
	if MatchByKeywords(path, keywords) {
		t.Errorf("Expected MatchByKeywords to return false for path: %s and keywords: %v", path, keywords)
	}

	keywords = []string{"path", "not", "directory"}
	if MatchByKeywords(path, keywords) {
		t.Errorf("Expected MatchByKeywords to return false for path: %s and keywords: %v", path, keywords)
	}
	keywords = []string{}
	if MatchByKeywords(path, keywords) {
		t.Errorf("Expected MatchByKeywords to return false for empty path and empty keywords")
	}

	path = ""
	keywords = []string{"path", "directory"}
	if MatchByKeywords(path, keywords) {
		t.Errorf("Expected MatchByKeywords to return false for empty path and keywords: %v", keywords)
	}
}
