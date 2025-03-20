// Copyright © 2016,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable // import "go.pennock.tech/tabular/texttable"

import (
	"fmt"

	"go.pennock.tech/tabular/color"
	"go.pennock.tech/tabular/texttable/decoration"
)

// SetDecoration sets a Decoration type for rendering a table.  The caller
// must provide the decoration object.
func (t *TextTable) SetDecoration(decor decoration.Decoration) *TextTable {
	t.decor = decor
	return t
}

// SetDecorationNamed selects a decoration by name, from the decoration package.
// It returns an error if the name is not known.
// If the name is not known, this is still registered in the TextTable, so that
// later attempts to render it will fail, instead of succeeding with unexpected
// output.
func (t *TextTable) SetDecorationNamed(n string) (*TextTable, error) {
	d := decoration.Named(n)
	t.decor = d
	if d == decoration.EmptyDecoration {
		return t, fmt.Errorf("unknown decoration name %q", n)
	}
	return t, nil
}

// SetFGColorNamed selects a color by name for the table borders.
func (t *TextTable) SetFGColorNamed(n string) (*TextTable, error) {
	c, err := color.ByHTMLNamedColor(n)
	if err != nil {
		return t, err
	}
	t.fgcolor = &c
	return t, nil
}

// SetFGColor directly assigns a color to the table borders.
func (t *TextTable) SetFGColor(c color.Color) *TextTable {
	t.fgcolor = &c
	return t
}

// RemoveFGColor removes the color for the table borders.
func (t *TextTable) RemoveFGColor() *TextTable {
	t.fgcolor = nil
	return t
}

// SetBGColorNamed selects a color by name for the table borders.
func (t *TextTable) SetBGColorNamed(n string) (*TextTable, error) {
	c, err := color.ByHTMLNamedColor(n)
	if err != nil {
		return t, err
	}
	t.bgcolor = &c
	return t, nil
}

// SetBGColor directly assigns a color to the table borders.
func (t *TextTable) SetBGColor(c color.Color) *TextTable {
	t.bgcolor = &c
	return t
}

// RemoveBGColor removes the color for the table borders.
func (t *TextTable) RemoveBGColor() *TextTable {
	t.bgcolor = nil
	return t
}

// SetBGSolid causes the table to not reset the background color for cell content
func (t *TextTable) SetBGSolid(onoff bool) *TextTable {
	if onoff {
		t.bgflags |= colorBGSolid
	} else {
		t.bgflags &= ^colorBGSolid
	}
	return t
}

// SetColorToEOL causes the table to not reset color at the end of the table lines
func (t *TextTable) SetColorToEOL(onoff bool) *TextTable {
	if onoff {
		t.bgflags |= colorToEOL
	} else {
		t.bgflags &= ^colorToEOL
	}
	return t
}

// SetCellFGColor sets the color for the foreground color of cells
func (t *TextTable) SetCellFGColor(c color.Color) *TextTable {
	t.cellfgcolor = &c
	return t
}

// SetCellBGColor sets the color for the background color of cells
func (t *TextTable) SetCellBGColor(c color.Color) *TextTable {
	t.cellbgcolor = &c
	return t
}

// RemoveCellFGColor removes the color for the foreground of cells
func (t *TextTable) RemoveCellFGColor(c color.Color) *TextTable {
	t.cellfgcolor = nil
	return t
}

// RemoveCellBGColor removes the color for the background of cells
func (t *TextTable) RemoveCellBGColor(c color.Color) *TextTable {
	t.cellbgcolor = nil
	return t
}
