// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// The Table interface is a thin wrapper around the actual *ATable struct, so that
// methods can all be on the interface and objects which embed an unnamed table
// can be used as tables.
//
// If you want to create and use a table for pre-canned rendering then look at
// sub-packages for TextTable, HTMLTable, etc wrappers.
type Table interface {
	AddRow(row *Row) Table
	AddSeparator() Table
	NColumns() int
	NRows() int
	Headers() []Cell
	AddHeaders(items ...interface{}) Table
	AllRows() []*Row
	NewRowSizedFor() *Row
	AppendNewRow() *Row
	AddRowItems(items ...interface{}) Table
	CellAt(location CellLocation) (*Cell, error)
	Column(int) *column

	RegisterPropertyCallback(PropertyOwner, callbackTime, cbTarget, PropertyCallback) error
	InvokeRenderCallbacks()

	// Also the other interfaces embedded in ATable:
	ErrorReceiver
	ErrorSource
	PropertyOwner
}

// Mutator style for setting table options?

// Do we have a more concise name than TerminalCellWidth for cells and tables?
// Length : Bytes?
// RuneCount : doesn't handle combining characters, etc
// Glyphs : doesn't handle variable-width characters (combining, wide)
//
// Unicode TR11 uses "cells" but we want that terminology for table cells.

// Inclining towards height as intrinsic based upon rows, but width as a per-backend thing.

//TODO: func (t *Table) Select(columns []int) *Table          {}
//TODO: func (t *Table) SelectByName(columns []string) *Table {}
//TODO: func (t *Table) Height() int                          {} // multi-line rows; separators; headers, box-lines

// Backend-specific:
//TODO: func (t *Table) TerminalCellWidth() int                          {}
//TODO: func (t *Table) ColumnTerminalCellWidth(column int) int          {}
//TODO: func (t *Table) ColumnTerminalCellWidthByName(column string) int {}

// What sort of streaming input do we want?
//TODO: func (t *Table) InputWriter(fieldSplitter func(string) []string) io.Writer

// Do we use options and one common Render function, or per-backend Render functions, pulled in via modules?

// TODO:
// Set a table property of rendering alignment-functions, registered per-type.
// Have a default set for new tables.  In a utility sub-package, used by relevant output sub-packages.
// Can then also have one for "alignment offset" and use a function for float type, registered as default on the table.
// Be able to set, get existing functions; set `nil` to delete.
// If can get and set, then can encapsulate.
// Should properties be in core tabular or in texttable?  Stored where?  Move, keep as now?
