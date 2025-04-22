package core

import "strings"

// MatchByKeywords checks if the path contains all the keywords in order.
func MatchByKeywords(path string, keywords []string) bool {
	path = strings.ToLower(path)
	lastIdx := -1
	for _, keyword := range keywords {
		idx := strings.Index(path[lastIdx+1:], strings.ToLower(keyword))
		if idx == -1 {
			return false
		}
		lastIdx += idx + len(keyword)
	}
	return true
}
