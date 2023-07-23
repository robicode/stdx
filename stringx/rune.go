package stringx

import (
	"unicode/utf8"
)

// RuneIndexFromStartByte returns the index of the rune that starts
// at the byte index index, or -1 if index is not a rune boundary.
func RuneIndexFromStartByte(str string, index int) int {
	runeMap := RuneIndicies(str)
	if runeMap == nil {
		return -1
	}
	for idx, pos := range runeMap {
		if pos == index {
			return idx
		}
	}
	return -1
}

// RuneIndicies returns the starting byte indicies of all
// runes in str as a map[int]int, or nil if str is empty.
func RuneIndicies(str string) map[int]int {
	if len(str) == 0 {
		return nil
	}

	runes := utf8.RuneCountInString(str)
	var pos int
	runeIndexes := make(map[int]int, runes)

	for i := 0; i < runes; i++ {
		_, rsz := utf8.DecodeRuneInString(str[pos:])
		runeIndexes[i] = pos
		pos += rsz
	}
	return runeIndexes
}
