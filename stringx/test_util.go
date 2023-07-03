package stringx

import (
	"fmt"
	"strings"
)

// quoteSliceElements is a helper function that makes []string slices easier to
// read by wrapping each element in quotes.
func quoteSliceElements(arr []string) string {
	var s string = "["

	for _, item := range arr {
		s = s + fmt.Sprintf("\"%s\", ", item)
	}

	if strings.HasSuffix(s, ", ") {
		s = s[:len(s)-2]
	}
	s = s + "]"

	return s
}

// equalSlices checks if two []string slices have equal
// content.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for n, item := range a {
		if b[n] != item {
			return false
		}
	}
	return true
}
