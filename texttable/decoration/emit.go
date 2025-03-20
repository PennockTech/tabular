// Copyright © 2016,2018,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import (
	"strings"

	"go.pennock.tech/tabular/properties/align"
)

// DividerSet provides the characters to use for vertical line-drawing when
// putting content into cells.
type DividerSet struct {
	Left  string
	Inner string
	Right string
}

type emitter struct {
	colWidths    []int
	totalWidth   int
	decor        *Decoration
	eol          string
	escStart     string
	escStop      string
	escCellStart string
	escClearEOL  string
	noResetEOL   bool
}

// ForColumnWidths returns an emitter object with methods for getting
// table lines as strings.
func (d *Decoration) ForColumnWidths(widths []int) emitter {
	// The totalWidth is so that buffers can be sized correctly.
	// We need the width of each cell, plus padding for each cell, plus dividers.
	// Possibly also with newlines.
	// However, separators can be Unicode and not fit in bytes.
	// So might not use this.
	totalWidth := 1 + len(widths)
	for _, w := range widths {
		if w >= 0 {
			totalWidth += w + 2
		}
	}
	return emitter{
		colWidths:  widths,
		decor:      d,
		totalWidth: totalWidth,
	}
}

func (e *emitter) SetEOL(eol string) {
	e.totalWidth = e.totalWidth - len(e.eol) + len(eol)
	e.eol = eol
}

func (e *emitter) SetANSIEscapes(escStart, escStop, escCellStart string) {
	e.escStart = escStart
	e.escStop = escStop
	if escCellStart != "" {
		e.escCellStart = escCellStart
	} else if e.escStart != "" {
		e.escCellStart = e.escStart
	}
	// clr_eol / el
	// without this, terminals are inconsistent in whether or not to stop the color at newline;
	// eg, in gnome-terminal, we get color to right margin on all lines except the first.
	e.escClearEOL = "\x1B[K"
}

func (e *emitter) SetNoResetEOL(onoff bool) {
	e.noResetEOL = onoff
}

func (e emitter) commonTemplateLine(left, horiz, cross, right string) string {
	if e.decor.isBoxless {
		return ""
	}
	fields := make([]string, 0, len(e.colWidths)*2+4)
	fields = append(fields, e.escStart)
	fields = append(fields, left)
	if len(e.colWidths) > 0 {
		for i := range e.colWidths {
			if e.colWidths[i] >= 0 {
				fields = append(fields, strings.Repeat(horiz, 2+e.colWidths[i]))
				fields = append(fields, cross)
			}
		}
		fields[len(fields)-1] = right
	} else {
		fields = append(fields, right)
	}
	if e.noResetEOL {
		// Without e.escStop here too, and with LineBottom appending e.escStop when e.noResetEOL,
		// we get perfection in xterm but gnome-terminal carries the background over for one more line, which seems like a bug.
		// So force a stop at the end of the lines.
		fields = append(fields, e.escClearEOL+e.escStop)
	} else {
		fields = append(fields, e.escStop)
	}
	fields = append(fields, e.eol)
	return strings.Join(fields, "")
}

func (e emitter) LineHeaderTop() string {
	return e.commonTemplateLine(e.decor.TopLeft, e.decor.HOuter, e.decor.HTopDown, e.decor.TopRight)
}

func (e emitter) LineHeaderBodySep() string {
	return e.commonTemplateLine(e.decor.HBLeft, e.decor.HOuter, e.decor.HBCross, e.decor.HBRight)
}

func (e emitter) LineBodyTop() string {
	return e.commonTemplateLine(e.decor.TopLeft, e.decor.HOuter, e.decor.BTopDown, e.decor.TopRight)
}

func (e emitter) LineBottom() string {
	return e.commonTemplateLine(e.decor.BottomLeft, e.decor.HOuter, e.decor.BBottomUp, e.decor.BottomRight)
}

func (e emitter) LineSeparator() string {
	return e.commonTemplateLine(e.decor.LeftBodyRule, e.decor.HRule, e.decor.CrossPiece, e.decor.RightBodyRule)
}

func (e emitter) LineHeaderBlanks() string {
	return e.commonTemplateLine(e.decor.VHeader, " ", e.decor.VHeader, e.decor.VHeader)
}

func (e emitter) LineBodyBlanks() string {
	return e.commonTemplateLine(e.decor.VBodyBorder, " ", e.decor.VBodyInner, e.decor.VBodyBorder)
}

func (e emitter) HeaderDividers() DividerSet {
	return DividerSet{
		Left:  e.decor.VHeader,
		Inner: e.decor.VHeader,
		Right: e.decor.VHeader,
	}
}

func (e emitter) BodyDividers() DividerSet {
	return DividerSet{
		Left:  e.decor.VBodyBorder,
		Inner: e.decor.VBodyInner,
		Right: e.decor.VBodyBorder,
	}
}

func (e emitter) commonRenderedLine(ds DividerSet, cellStrs []WidthString, colAligns []align.Alignment) string {
	fields := make([]string, 0, len(e.colWidths)*2+1)
	eolReset := e.escStop
	if e.noResetEOL {
		eolReset = e.escClearEOL + e.escStop
	}
	if ds.Left != "" {
		fields = append(fields, e.escStart+ds.Left+e.escCellStart)
	}
	for i := range e.colWidths {
		if e.colWidths[i] >= 0 {
			fields = append(fields, cellStrs[i].WithinWidthAligned(e.colWidths[i], colAligns[i]))
			if ds.Inner != "" {
				fields = append(fields, e.escStart+ds.Inner+e.escCellStart)
			}
		}
	}
	if ds.Right != "" && ds.Inner != "" {
		fields[len(fields)-1] = e.escStart + ds.Right + eolReset
	} else if ds.Right != "" {
		fields = append(fields, e.escStart+ds.Right+eolReset)
	} else if ds.Inner != "" {
		fields = fields[:len(fields)-1]
	}
	return strings.Join(fields, " ") + e.eol
}

// HeaderLineRendered is internal to tabular, package predates 'internal' else
// this would be hidden.  It's used to render a single header line of a
// texttable.
func (e emitter) HeaderLineRendered(cellStrs []WidthString, colAligns []align.Alignment) string {
	return e.commonRenderedLine(e.HeaderDividers(), cellStrs, colAligns)
}

// BodyLineRendered is internal to tabular, package predates 'internal' else
// this would be hidden.  It's used to render a single line of a texttable.
func (e emitter) BodyLineRendered(cellStrs []WidthString, colAligns []align.Alignment) string {
	return e.commonRenderedLine(e.BodyDividers(), cellStrs, colAligns)
}
