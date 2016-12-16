// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"fmt"
	"strings"

	"go.pennock.tech/tabular/length"
)

// A Cell is one item in a table; it holds an object and fields calculated
// from it.  If the object added is mutated after addition, it is the mutator's
// responsibility to call Update.
type Cell struct {
	raw    interface{}
	str    string
	width  int
	height int
	empty  bool

	propertyImpl
	callbacks callbackSet

	// 0 is "unknown", first column is 1; rowNum via inRow below
	columnNum int

	// When within a row, holds row
	inRow *Row

	// Do we need to hold something saying that we're in a row, like a row
	// holds a table, so that if a .SetItem() is done, we can re-trigger row
	// (and column/table) callbacks?

	// NB: at present, this just means "calculate the contained raw, then use the values"
	// but that doesn't need to be the case;
	mustCalc bool
	// Also: how to we want to handle mutating contained objects and staleness?
	// Interface which must have a IsUpdatedSince(yourStateId) function?
	// Or just a .Refresh() which redoes things?
}

// NewCell creates a Cell, handling object rendering at init time.
// TODO: handle rune as rune, or as int?  Any special flags to use?
func NewCell(object interface{}) Cell {
	c := Cell{raw: object, empty: false}
	c.Update()
	return c
}

// Update changes metadata to reflect the current state of the object stored in
// a cell.
func (c *Cell) Update() {
	c.empty = false
	c.mustCalc = false

	switch o := c.raw.(type) {
	case nil:
		c.empty = true
		c.str = ""
		c.width = 0
		c.height = 0
		return
	case Cell:
		c.str = o.str
		c.width = o.width
		c.height = o.height
		c.empty = o.empty
		return

	// After this point, MUST set .str
	case string:
		c.str = o
	case rune:
	case Stringer:
		c.str = o.String()
	case GoStringer:
		c.str = o.GoString()
	case error:
		c.str = o.Error()
	default:
		c.str = fmt.Sprintf("%v", o)
	}

	overrideOnly := false
	if c.str == "" {
		c.width = 0
		c.height = 0
		c.empty = true
		overrideOnly = true
	}

	if h, ok := c.raw.(Heighter); ok {
		c.height = h.Height()
	} else if !overrideOnly {
		c.height = 1 + strings.Count(c.str, "\n")
		if strings.HasSuffix(c.str, "\n") {
			c.height -= 1
		}
	}
	if w, ok := c.raw.(TerminalCellWidther); ok {
		c.width = w.TerminalCellWidth()
	} else if !overrideOnly {
		c.width = length.LongestLineCells(c.str)
	}
}

// updateCache does an update and mutates fields destructively; it should only be used on a transient copy
// of a persistent object, not the persistent object itself.
func (c *Cell) updateCache() {
	c.Update()
	c.mustCalc = false
}

// Item returns the object stored inside a cell.
func (c Cell) Item() interface{} {
	return c.raw
}

// String returns some string representation of the content of a cell.
func (c Cell) String() string {
	if c.mustCalc {
		(&c).updateCache()
	}
	return c.str
}

// Lines returns the string representation of the content of a cell, as
// a splice of strings, one per line, without newlines; if the string has a
// final \n then there will NOT be an extra empty in the result for the
// "empty" final segment.
func (c Cell) Lines() []string {
	return length.Lines(c.String())
}

// TerminalCellWidth returns the number of terminal cells which we believe
// are necessary to render the contents of the object stored in the cell.
// This is overriden by a TerminalCellWidth method on the object being stored.
// To a first approximation, this is how many runes are in a cell, but we
// handle combining characters, wide characters, etc.
func (c Cell) TerminalCellWidth() int {
	if c.mustCalc {
		(&c).updateCache()
	}
	if c.width < 0 {
		return 0
	}
	return c.width
}

// Height returns the height of a cell; usually this is the number of lines in
// the string representation of the object stored in the cell, as returned by
// Lines, but an object which has a Height method will override this.
func (c Cell) Height() int {
	if c.mustCalc {
		(&c).updateCache()
	}
	if c.height < 1 {
		if c.TerminalCellWidth() > 0 {
			return 1
		}
		return 0
	}
	return c.height
}

func (c Cell) columnOfTable() *column {
	if c.columnNum < 1 {
		return nil
	}
	if c.inRow == nil || c.inRow.inTable == nil {
		return nil
	}
	t := c.inRow.inTable
	// CellLocation 1-indexed, t []column 0-indexed
	if c.columnNum > t.nColumns {
		return nil
	}
	return &t.columns[c.columnNum-1]
}
