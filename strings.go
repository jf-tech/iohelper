package iohelper

import (
	"strings"
)

// RunePtr returns a pointer to a rune.
func RunePtr(r rune) *rune {
	return &r
}

// IndexWithEsc is similar to strings.Index but taking escape sequnce into consideration.
// For example, IndexWithEsc("abc%|efg|xyz", "|", RunePtr("%")) would return 8, not 4.
func IndexWithEsc(s, delim string, esc *rune) int {
	if len(delim) == 0 {
		return 0
	}
	if len(s) == 0 {
		return -1
	}
	if esc == nil {
		return strings.Index(s, delim)
	}

	sRunes := []rune(s)
	delimRunes := []rune(delim)
	escRune := *esc

	// Yes this old dumb double loop isn't the most efficient algo but it's super easy and simple to understand
	// and bug free compared with fancy strings.Index or bytes.Index which could potentially lead to index errors
	// and/or rune/utf-8 bugs. Plus for vast majority of use cases, delim will be of a single rune, so effectively
	// not much perf penalty at all.
	for i := 0; i < len(sRunes)-len(delimRunes)+1; i++ {
		if sRunes[i] == escRune {
			// skip the escaped rune (aka the rune after the escape rune)
			i++
			continue
		}
		delimFound := true
		for j := 0; j < len(delimRunes); j++ {
			if sRunes[i+j] != delimRunes[j] {
				delimFound = false
				break
			}
		}
		if delimFound {
			return len(string(sRunes[:i]))
		}
	}

	return -1
}

// SplitWithEsc is similar to strings.Split but taking escape sequence into consideration.
// For example, SplitWithEsc("abc%|efg|xyz", "|", RunePtr("%")) would return []string{"abc%|efg", "xyz"}.
func SplitWithEsc(s, delim string, esc *rune) []string {
	if len(delim) == 0 || esc == nil {
		return strings.Split(s, delim)
	}
	// From here on, delim != empty **and** esc is set.
	var split []string
	for delimIndex := IndexWithEsc(s, delim, esc); delimIndex >= 0; delimIndex = IndexWithEsc(s, delim, esc) {
		split = append(split, s[:delimIndex])
		s = s[delimIndex+len(delim):]
	}
	split = append(split, s)
	return split
}

// Unescape unescapes a string with escape sequence.
// For example, SplitWithEsc("abc%|efg", RunePtr("%")) would return "abc|efg".
func Unescape(s string, esc *rune) string {
	if esc == nil {
		return s
	}
	sRunes := []rune(s)
	escRune := *esc
	for i := 0; i < len(sRunes); i++ {
		if sRunes[i] != escRune {
			continue
		}
		copy(sRunes[i:], sRunes[i+1:])
		sRunes = sRunes[:len(sRunes)-1]
	}
	return string(sRunes)
}
