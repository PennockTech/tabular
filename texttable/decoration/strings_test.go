// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration

import (
	"testing"
	"unicode/utf8"

	"github.com/liquidgecka/testlib"
)

func TestWithinWidth(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	for testI, tuple := range []struct {
		In   string
		InW  int
		OutW int
		Want string
	}{
		// Items labelled `non-spec` are recording the current behaviour,
		// ensuring that stupid/buggy inputs are handled sanely without
		// crashing, but are not guarantees and output may change in the
		// future to "something else vaguely sane(r)".
		{"a", 0, 5, "a    "},
		{"aa", 0, 5, "aa   "},
		{"aaaa", 0, 5, "aaaa "},
		{"aaaaa", 0, 5, "aaaaa"},
		{"bbbbbb", 0, 5, "bbbbbb"}, // non-spec
		{"", 0, 5, "     "},
		{"cc", 2, 5, "cc   "},
		{"cc", 1, 5, "cc    "}, // non-spec
		{"x", -5, 5, "     "},  // non-spec
		//
		{"£", 0, 2, "£"},  // non-spec
		{"£", 2, 2, "£"},  // non-spec, matches byte len
		{"£", 1, 2, "£ "}, // spec, coerce length
		{"£", -1, 2, "£ "},
		// TODO: test 2-cellwidth glyphs here, when we want to support them
	} {
		ws := WidthString{
			S: tuple.In,
			W: tuple.InW,
		}
		switch ws.W {
		case 0:
			ws.W = len(ws.S)
		case -1:
			ws.W = utf8.RuneCountInString(ws.S)
		}
		have := ws.WithinWidth(tuple.OutW)
		T.Equalf(have, tuple.Want, "[%d] right-padded string equality (from %v)", testI, ws)
	}
}
