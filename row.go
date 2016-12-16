// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"errors"
)

// A Row represents a row in a Table.
type Row struct {
	*ErrorContainer
	propertyImpl
	cells              []Cell
	rowCellCallbacks   callbackSet
	rowItselfCallbacks callbackSet
	inTable            *ATable
	isSeparator        bool
	rowNum             int
}

// AddError records that an error has happened when dealing with a row.
func (r *Row) AddError(e error) {
	if r.ErrorContainer == nil {
		r.ErrorContainer = NewErrorContainer()
	}
	r.ErrorContainer.AddError(e)
}

// NewRow creates a new Row.
func NewRow() *Row {
	return NewRowWithCapacity(10)
}

// NewRowWithCapacity create a new row pre-sized to contain the given
// number of Cells.
func NewRowWithCapacity(c int) *Row {
	return &Row{cells: make([]Cell, 0, c)}
}

// NewRowSizedFor creates a new Row sized for the given table which
// it is called upon, assuming that NColumns() is usefully available.
func (t *ATable) NewRowSizedFor() *Row {
	return &Row{cells: make([]Cell, 0, t.NColumns())}
}

// AppendNewRow creates a new row sized for the table (per NewRowSizedFor)
// and adds it to the Table, returning the row.
func (t *ATable) AppendNewRow() *Row {
	r := &Row{
		cells: make([]Cell, 0, t.NColumns()),
	}
	t.AddRow(r) // will also set the location
	return r
}

// Add adds one cell to this row, and returns the row for chaining, thus
// r.Add(c1).Add(c2).Add(c3)
func (r *Row) Add(c Cell) *Row {
	if r.cells == nil {
		r.AddError(errors.New("can't add cells to a non-cell row"))
		return r
	}
	r.cells = append(r.cells, c)
	column := len(r.cells)
	ptr := &r.cells[column-1]
	ptr.inRow = r
	ptr.columnNum = column
	invokePropertyCallbacks(r.rowCellCallbacks, CB_AT_ADD, ptr, r.ErrorContainer)
	return r
}

// Cells returns an iterable of the cells in a row.  If it returns nil
// then you have a non-cell row (probably a separator).
func (r *Row) Cells() []Cell {
	// Possible v2 change: should we dup this as we do for t.AllRows?
	// This isn't array-of-pointer so gets more expensive.
	return r.cells
}

// AddRowItems creates a row from the passed items and adds it to the table, returning
// the table for chaining.
func (t *ATable) AddRowItems(items ...interface{}) Table {
	r := NewRowWithCapacity(len(items))
	for i := range items {
		r.Add(NewCell(items[i]))
	}
	return t.AddRow(r)
}

// newSeparator returns a separator row.
// Open Question for v2:
//   object model might need clean-up, to have interfaces satisfied by rows and separator be a distinct type?
func newSeparator() *Row { return &Row{isSeparator: true} }

// IsSeparator is true iff a row is a "separator".
func (r *Row) IsSeparator() bool {
	return r.isSeparator
}
