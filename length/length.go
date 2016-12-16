// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package length // import "go.pennock.tech/tabular/length"

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// StringBytes is the number of bytes in a string
func StringBytes(s string) int {
	return len(s)
}

// StringRunes is the number of runes in a string
func StringRunes(s string) int {
	return utf8.RuneCountInString(s)
}

// StringCells is an attempt to guess the number of display cells in a
// fixed-grid terminal window system of cells required for the characters in
// the string.
//
// It does handle full-width and combining, but doesn't handle wide emojis (at
// this time).
//
// The implementation of this function is subject to change as we try to get
// closer.
func StringCells(s string) int {
	return runewidth.StringWidth(s)
}

// Lines breaks a string apart into lines; a final newline in the string does
// not add a blank final line, but two or more final newlines will add
// one-less-than-count blank lines.
func Lines(s string) []string {
	ss := strings.Split(s, "\n")
	if ss[len(ss)-1] == "" {
		ss = ss[:len(ss)-1]
	}
	return ss
}

// LongestLineBytes returns the length of the longest virtual line in a string
// containing embedded newlines, measuring length per StringBytes.
func LongestLineBytes(s string) int {
	ss := Lines(s)
	switch len(ss) {
	case 0:
		return 0
	case 1:
		return StringBytes(ss[0])
	}
	max := 0
	for i := range ss {
		t := StringBytes(ss[i])
		if t > max {
			max = t
		}
	}
	return max
}

// LongestLineRunes returns the length of the longest virtual line in a string
// containing embedded newlines, measuring length per StringRunes.
func LongestLineRunes(s string) int {
	ss := Lines(s)
	switch len(ss) {
	case 0:
		return 0
	case 1:
		return StringRunes(ss[0])
	}
	max := 0
	for i := range ss {
		t := StringRunes(ss[i])
		if t > max {
			max = t
		}
	}
	return max
}

// LongestLineCells returns the length of the longest virtual line in a string
// containing embedded newlines, measuring length per StringCells.
func LongestLineCells(s string) int {
	ss := Lines(s)
	switch len(ss) {
	case 0:
		return 0
	case 1:
		return StringCells(ss[0])
	}
	max := 0
	for i := range ss {
		t := StringCells(ss[i])
		if t > max {
			max = t
		}
	}
	return max
}
