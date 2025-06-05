package core

import (
	"path/filepath"
	"strings"
)

// MatchByKeywords checks if the path contains all the keywords in order.
func MatchByKeywords(path string, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}

	path = strings.ToLower(path)

	lastKeyword := keywords[len(keywords)-1]
	rest := keywords[:len(keywords)-1]

	idx := strings.LastIndex(path, strings.ToLower(lastKeyword))
	if idx == -1 {
		return false
	}

	after := path[idx+len(lastKeyword):]
	if strings.ContainsAny(after, string(filepath.Separator)) {
		return false
	}

	path = path[:idx]

	for i := len(rest) - 1; i >= 0; i-- {
		k := strings.ToLower(rest[i])
		idx = strings.LastIndex(path, k)
		if idx == -1 {
			return false
		}
		path = path[:idx]
	}

	return true
}
