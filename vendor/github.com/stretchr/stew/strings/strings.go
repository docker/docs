package strings

import (
	"bytes"
	"strings"
)

// MergeStrings merges many strings together.
func MergeStrings(stringArray ...string) string {

	var buffer bytes.Buffer
	for _, v := range stringArray {
		buffer.WriteString(v)
	}
	return buffer.String()

}

// MergeStringsReversed merges many strings together backwards.
func MergeStringsReversed(stringArray ...string) string {

	var buffer bytes.Buffer
	for vi := len(stringArray); vi >= 0; vi-- {
		buffer.WriteString(stringArray[vi])
	}
	return buffer.String()

}

// JoinStrings joins many strings together separated by the specified separator.
func JoinStrings(separator string, stringArray ...string) string {

	var buffer bytes.Buffer
	var max int = len(stringArray) - 1
	for vi, v := range stringArray {
		buffer.WriteString(v)
		if vi < max {
			buffer.WriteString(separator)
		}
	}
	return buffer.String()

}

// JoinStringsReversed joins many strings together backwards separated by the specified separator.
func JoinStringsReversed(separator string, stringArray ...string) string {

	var buffer bytes.Buffer

	for vi := len(stringArray) - 1; vi >= 0; vi-- {
		buffer.WriteString(stringArray[vi])
		if vi > 0 {
			buffer.WriteString(separator)
		}
	}

	return buffer.String()

}

// SplitBy splits a string to segments based on the return of the decider function.
func SplitBy(s string, decider func(r rune) bool) []string {

	// split by caps
	var segments []string
	var currentSeg []rune

	for rIndex, r := range s {

		if decider(r) {
			// new word
			if len(currentSeg) > 0 {
				segments = append(segments, string(currentSeg))
				currentSeg = nil
			}
		}

		currentSeg = append(currentSeg, r)

		// is this the last one?
		if rIndex == len(s)-1 {
			segments = append(segments, string(currentSeg))
		}

	}

	return segments
}

// SplitByCamelCase splits the string up by each capital character.
func SplitByCamelCase(s string) []string {
	return SplitBy(s, func(r rune) bool {
		return strings.ToUpper(string(r)) == string(r)
	})
}
