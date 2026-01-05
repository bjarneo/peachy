package shared

import (
	"strings"
)

// IsValidImage checks if a file path has a valid image extension
func IsValidImage(path string) bool {
	lower := strings.ToLower(path)
	for _, ext := range ValidImageExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
