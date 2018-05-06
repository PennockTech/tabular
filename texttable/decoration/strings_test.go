// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import (
	"testing"
	"unicode/utf8"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular/properties/align"
)

func TestWithinWidth(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	var (
		ourLeft   align.Alignment = align.Left
		ourCenter align.Alignment = align.Center
		ourRight  align.Alignment = align.Right
	)

	for testI, tuple := range []struct {
		In      string
		InW     int
		OutW    int
		Want    string
		aligned *align.Alignment
	}{
		// Items labelled `non-spec` are recording the current behaviour,
		// ensuring that stupid/buggy inputs are handled sanely without
		// crashing, but are not guarantees and output may change in the
		// future to "something else vaguely sane(r)".
		{"a", 0, 5, "a    ", nil},
		{"a", 0, 5, "a    ", &ourLeft},
		{"a", 0, 5, "    a", &ourRight},
		{"a", 0, 5, "  a  ", &ourCenter},
		{"aa", 0, 5, "aa   ", nil},
		{"aaaa", 0, 5, "aaaa ", nil},
		{"aaaaa", 0, 5, "aaaaa", nil},
		{"bbbbbb", 0, 5, "bbbbbb", nil}, // non-spec
		{"", 0, 5, "     ", nil},
		{"cc", 2, 5, "cc   ", nil},
		{"cc", 1, 5, "cc    ", nil}, // non-spec
		{"x", -5, 5, "     ", nil},  // non-spec
		//
		{"£", 0, 2, "£", nil},  // non-spec
		{"£", 2, 2, "£", nil},  // non-spec, matches byte len
		{"£", 1, 2, "£ ", nil}, // spec, coerce length
		{"£", -1, 2, "£ ", nil},
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
		var have string
		if tuple.aligned == nil {
			have = ws.WithinWidth(tuple.OutW)
		} else {
			have = ws.WithinWidthAligned(tuple.OutW, *tuple.aligned)
		}
		T.Equalf(have, tuple.Want, "[%d] right-padded string equality (from %v)", testI, ws)
	}

	T.ExpectPanic(func() {
		ws := WidthString{S: "x", W: 1}
		have := ws.WithinWidthAligned(5, align.TestingInvalidAlignment())
		panic("did not panic, got: " + have)
	}, "unhandled alignment", "should panic with invalid alignment")
}
