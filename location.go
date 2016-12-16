// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// We measure from 1:1 so that 0:0 means "not initialized" and if either x or y
// is zero in a cell, we know the data is invalid.  For some contexts, one or
// the other might be zero, to describe an entire row or column or just unknown.
type CellLocation struct {
	Row    int
	Column int
}

// Location determines the current location of a cell
func (c Cell) Location() (loc CellLocation) {
	loc.Column = c.columnNum
	if c.inRow != nil {
		loc.Row = c.inRow.rowNum
	}
	return
}

// Location returns a row's CellLocation where the column is 0
func (r *Row) Location() (loc CellLocation) {
	loc.Column = 0
	loc.Row = r.rowNum
	return
}
