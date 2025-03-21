// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import (
	"strings"

	"go.pennock.tech/tabular/properties/align"
)

type WidthString struct {
	S string
	W int
}

func (ws WidthString) WithinWidth(available int) string { return ws.WithinWidthAligned(available, nil) }

func (ws WidthString) WithinWidthAligned(available int, howAlign align.Alignment) string {
	if ws.W < 0 {
		return strings.Repeat(" ", available)
	}
	if howAlign == nil {
		howAlign = align.Left
	}

	// this will need to change when we support more than basic left/right alignment
	pad := max(available-ws.W, 0)
	switch howAlign {
	case align.Left:
		return ws.S + strings.Repeat(" ", pad)
	case align.Right:
		return strings.Repeat(" ", pad) + ws.S
	case align.Center:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + ws.S + strings.Repeat(" ", right)
	default:
		panic("unhandled alignment")
	}
}
