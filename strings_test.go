package iohelper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexWithEsc(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    string
		delim    string
		esc      *rune
		expected int
	}{
		// All edge cases:
		{
			name:     "delim empty",
			input:    "abc",
			delim:    "",
			esc:      RunePtr(rune('宇')),
			expected: 0,
		},
		{
			name:     "esc empty",
			input:    "abc",
			delim:    "bc",
			esc:      nil,
			expected: 1,
		},
		{
			name:     "input empty, delim non empty, esc non empty",
			input:    "",
			delim:    "abc",
			esc:      RunePtr(rune('宙')),
			expected: -1,
		},
		// normal non empty cases:
		{
			name:     "len(input) < len(delim)",
			input:    "a",
			delim:    "abc",
			esc:      RunePtr(rune('洪')),
			expected: -1,
		},
		{
			name:     "len(input) == len(delim), esc not present",
			input:    "abc",
			delim:    "abc",
			esc:      RunePtr(rune('荒')),
			expected: 0,
		},
		{
			name:     "len(input) > len(delim), esc not present",
			input:    "мир во всем мире",
			delim:    "мире",
			esc:      RunePtr(rune('Ф')),
			expected: len("мир во всем "),
		},
		{
			name:     "len(input) > len(delim), esc present",
			input:    "мир во всем /мире",
			delim:    "мире",
			esc:      RunePtr(rune('/')),
			expected: -1,
		},
		{
			name:     "len(input) > len(delim), esc present",
			input:    "мир во всем ξξмире",
			delim:    "мире",
			esc:      RunePtr(rune('ξ')),
			expected: len("мир во всем ξξ"),
		},
		{
			name:     "len(input) > len(delim), consecutive esc present",
			input:    "мир во вξξξξξсем ξξмире",
			delim:    "ире",
			esc:      RunePtr(rune('ξ')),
			expected: len("мир во вξξξξξсем ξξм"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IndexWithEsc(test.input, test.delim, test.esc))
			if test.expected >= 0 {
				assert.True(t, strings.HasPrefix(string([]byte(test.input)[test.expected:]), test.delim))
			}
		})
	}
}

func TestSplitWithEsc(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    string
		delim    string
		esc      *rune
		expected []string
	}{
		{
			name:     "delim empty",
			input:    "abc",
			delim:    "",
			esc:      RunePtr(rune('宇')),
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "esc not set",
			input:    "",
			delim:    "abc",
			esc:      nil,
			expected: []string{""},
		},
		{
			name:     "esc set, delim not found",
			input:    "?xyz",
			delim:    "xyz",
			esc:      RunePtr(rune('?')),
			expected: []string{"?xyz"},
		},
		{
			name:     "esc set, delim found",
			input:    "a*bc/*d*efg",
			delim:    "*",
			esc:      RunePtr(rune('/')),
			expected: []string{"a", "bc/*d", "efg"},
		},
		{
			name:     "esc set, delim not empty, input empty",
			input:    "",
			delim:    "*",
			esc:      RunePtr(rune('/')),
			expected: []string{""},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, SplitWithEsc(test.input, test.delim, test.esc))
		})
	}
}

func TestUnescape(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    string
		esc      *rune
		expected string
	}{
		{
			name:     "esc not set",
			input:    "abc",
			esc:      nil,
			expected: "abc",
		},
		{
			name:     "esc set, input empty",
			input:    "",
			esc:      RunePtr(rune('宇')),
			expected: "",
		},
		{
			name:     "esc set, input non empty",
			input:    "ξξabcξdξ",
			esc:      RunePtr(rune('ξ')),
			expected: "ξabcd",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, Unescape(test.input, test.esc))
		})
	}
}
