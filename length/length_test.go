// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package length_test

import (
	"testing"

	"github.com/liquidgecka/testlib"

	"github.com/PennockTech/tabular/length"
)

func TestStringLengths(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	for i, tuple := range []struct {
		s string
		b int // bytes length
		r int // runes length
		c int // cells length
	}{
		{"", 0, 0, 0},
		{"a", 1, 1, 1},
		{" b", 2, 2, 2},
		{"£", 2, 1, 1},              // %C2%A3
		{"⌘€", 6, 2, 2},             // %E2%8C%98 %E2%82%AC
		{"á", 2, 1, 1},              // %C3%A1
		{"\xCC\x81a", 3, 2, 1},      // %CC%81 = COMBINING ACUTE ACCENT
		{"ａ", 3, 1, 2},              // %EF%BD%81 = FULLWIDTH LATIN SMALL LETTER A
		{"\xE2\x80\xAA", 3, 1, 0},   //  %E2%80%AA = 202a = LEFT-TO-RIGHT EMBEDDING (a control, zero-width)
		{"a\xC2\xA0b", 4, 3, 3},     // %C2%A0 = a0 = NO-BREAK SPACE
		{"a\xE2\x80\x82b", 5, 3, 3}, // \xE2\x80\x82 = 2002 = EN SPACE
		{"a\xE2\x80\x83b", 5, 3, 3}, // \xE2\x80\x83 = 2003 = EM SPACE
		{"a\xE2\x80\x8bb", 5, 3, 2}, // \xE2\x80\x8B = 200b = ZERO WIDTH SPACE
		{"a\xE3\x80\x80b", 5, 3, 4}, // \xE3\x80\x80 = 3000 = IDEOGRAPHIC SPACE (2 wide)
		// Broken FIXME items (presence here is not API guarantee):
		{"💪", 4, 1, 1 /* want: 2 */}, // %F0%9F%92%AA = FLEXED BICEPS, followed by a space

	} {
		T.Equalf(length.StringBytes(tuple.s), tuple.b, "bytes length test [%d] string %q", i, tuple.s)
		T.Equalf(length.StringRunes(tuple.s), tuple.r, "runes length test [%d] string %q", i, tuple.s)
		T.Equalf(length.StringCells(tuple.s), tuple.c, "cells length test [%d] string %q", i, tuple.s)
	}

}

func TestMultiLineStringLengths(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	for i, tuple := range []struct {
		s string
		l int // lines count
		b int // bytes length
		r int // runes length
		c int // cells length
	}{
		{"", 0, 0, 0, 0},
		{"\n", 1, 0, 0, 0}, // "" on first line, drop second line
		{"\n\n", 2, 0, 0, 0},
		{"a", 1, 1, 1, 1},
		{"a\n", 1, 1, 1, 1},
		{"a\n\n", 2, 1, 1, 1},
		{"\nbbb", 2, 3, 3, 3},
		{"\nb\n", 2, 1, 1, 1},
		{"a\nbb\nc", 3, 2, 2, 2},
		{"a\nbb", 2, 2, 2, 2},
		{"aa\nb", 2, 2, 2, 2},
		{"\xCC\x81a\nbb", 2, 3, 2, 2}, // 1st line: 3 bytes, 2 runes, 1 cell
		{"\xCC\x81a\nb", 2, 3, 2, 1},
		{"\xCC\x81a\n\n", 2, 3, 2, 1},
		{"\xCC\x81a\nｂ", 2, 3, 2, 2}, // 2nd line FULLWIDTH LATIN SMALL LETTER B
		{"💪\nbb", 2, 4, 2, 2},        // 1st line: 4 bytes, 1 rune, currently 1 cell but should be 2
		{"💪\nb", 2, 4, 1, 1 /* want: 2 */},
		{"💪\n\n", 2, 4, 1, 1 /* want: 2 */},
		{"💪\nｂ", 2, 4, 1, 2}, // 2nd line is fullwidth, 2 cells (1st line theoretically 2, but currently 1)
	} {
		T.Equalf(len(length.Lines(tuple.s)), tuple.l, "lines count test [%d] string %q", i, tuple.s)
		T.Equalf(length.LongestLineBytes(tuple.s), tuple.b, "bytes line length test [%d] string %q", i, tuple.s)
		T.Equalf(length.LongestLineRunes(tuple.s), tuple.r, "runes line length test [%d] string %q", i, tuple.s)
		T.Equalf(length.LongestLineCells(tuple.s), tuple.c, "cells line length test [%d] string %q", i, tuple.s)
	}
}
