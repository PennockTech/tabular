// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration

import "strings"

// DividerSet provides the characters to use for vertical line-drawing when
// putting content into cells.
type DividerSet struct {
	Left  string
	Inner string
	Right string
}

type emitter struct {
	colWidths  []int
	totalWidth int
	decor      *Decoration
	eol        string
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
		totalWidth += w + 2
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

func (e emitter) commonTemplateLine(left, horiz, cross, right string) string {
	if e.decor.isBoxless {
		return ""
	}
	fields := make([]string, 0, len(e.colWidths)*2+2)
	fields = append(fields, left)
	if len(e.colWidths) > 0 {
		for i := range e.colWidths {
			fields = append(fields, strings.Repeat(horiz, 2+e.colWidths[i]))
			fields = append(fields, cross)
		}
		fields[len(fields)-1] = right
	} else {
		fields = append(fields, right)
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

func (e emitter) commonRenderedLine(ds DividerSet, cellStrs []WidthString) string {
	fields := make([]string, 0, len(e.colWidths)*2+1)
	if ds.Left != "" {
		fields = append(fields, ds.Left)
	}
	for i := range e.colWidths {
		fields = append(fields, cellStrs[i].WithinWidth(e.colWidths[i]))
		if ds.Inner != "" {
			fields = append(fields, ds.Inner)
		}
	}
	if ds.Right != "" && ds.Inner != "" {
		fields[len(fields)-1] = ds.Right
	} else if ds.Right != "" {
		fields = append(fields, ds.Right)
	} else if ds.Inner != "" {
		fields = fields[:len(fields)-1]
	}
	return strings.Join(fields, " ") + e.eol
}

func (e emitter) HeaderLineRendered(cellStrs []WidthString) string {
	return e.commonRenderedLine(e.HeaderDividers(), cellStrs)
}

func (e emitter) BodyLineRendered(cellStrs []WidthString) string {
	return e.commonRenderedLine(e.BodyDividers(), cellStrs)
}
