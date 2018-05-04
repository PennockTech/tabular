// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// An ATable is the top-level container for a grid in tabular.
// You should not be declaring fields to be of type ATable; instead,
// use the Table interface.
//
// An ATable consists of rows of cells; there may also be a header row
// which defines names for columns.
// An ATable, a Column, a Row and a Cell can all contain both Properties,
// and callbacks for updating properties.  The callbacks can each individually
// be registered to activate upon object addition or upon object rendering.
// Property Callbacks can take various child objects.
// For instance, a table-level callback upon cells might update a namespaced
// Property which holds the cell's rendered width when emitted as characters
// for terminal display.
type ATable struct {
	*ErrorContainer
	propertyImpl
	headerRow                 *Row
	rows                      []*Row
	nColumns                  int
	columnNames               map[string]int
	columns                   []column    // has nColumns+1 entries
	tableItselfCallbacks      callbackSet // only useful for render-time
	tableCellCallbacks        callbackSet
	tableRowAdditionCallbacks callbackSet
}

type column struct {
	Name                  string
	ofTable               *ATable
	cellCallbacks         callbackSet
	columnItselfCallbacks callbackSet
	propertyImpl
}

// New creates a new empty ATable, which satisfies Table.
func New() *ATable {
	t := &ATable{
		ErrorContainer: NewErrorContainer(),
		rows:           make([]*Row, 0, 50),
		columns:        make([]column, 1, 10),
	}
	t.columns[0].ofTable = t
	return t
}

func (t *ATable) resizeColumnsAtLeast(newCount int) {
	if newCount <= t.nColumns {
		return
	}
	// include space for column 0 as well
	// this could be optimized to reduce copying while len<cap
	extraColumns := make([]column, newCount+1-len(t.columns))
	for i := range extraColumns {
		extraColumns[i].ofTable = t
	}
	t.columns = append(t.columns, extraColumns...)
	t.nColumns = newCount
}

// Column returns a representation of a given column in the table.
// Column counting starts at 1.  Providing an invalid column count
// returns nil.  Column 0 also exists but is the implicit defaults
// column, letting you set default column properties and then
// override for other columns.
func (t *ATable) Column(n int) *column {
	if n < 0 || n > t.nColumns {
		return nil
	}
	return &t.columns[n]
}

// NewTableWithHeaders() might allow auto-sizing?

// AddRow adds a *Row to the *ATable, returning the table to allow for chaining.
// Any errors accumulate in the table.
// Any existing errors in the row become table errors.
func (t *ATable) AddRow(row *Row) Table {
	t.rows = append(t.rows, row)
	row.inTable = t
	row.rowNum = len(t.rows)
	t.resizeColumnsAtLeast(len(row.cells))
	// swallow existing errors
	es := row.Errors()
	if es != nil {
		t.AddErrorList(es)
		// TODO: do we want an AddContextualErrorList which prepends a row-id to each error in the table?
	}
	// divert new errors to come directly to us
	row.ErrorContainer = t.ErrorContainer

	// Invoke property-updating callbacks
	invokePropertyCallbacks(row.rowItselfCallbacks, CB_AT_ADD, row, t.ErrorContainer)
	invokePropertyCallbacks(t.tableRowAdditionCallbacks, CB_AT_ADD, row, t.ErrorContainer)
	for i := range row.cells {
		// using row.EC just in case do switch to something which prepends row-ids later
		ptr := &row.cells[i]
		if col := ptr.columnOfTable(); col != nil {
			invokePropertyCallbacks(col.cellCallbacks, CB_AT_ADD, ptr, row.ErrorContainer)
		}
		invokePropertyCallbacks(t.tableCellCallbacks, CB_AT_ADD, ptr, row.ErrorContainer)
	}

	return t
}

// AddSeparator adds a rule to the table.
func (t *ATable) AddSeparator() Table {
	sep := newSeparator()
	t.rows = append(t.rows, sep)
	sep.inTable = t
	sep.rowNum = len(t.rows)
	return t
}

// NColumns says how many columns are in the table
func (t *ATable) NColumns() int {
	// Because we maintain columns as we add rows, this is easy
	return t.nColumns
}

// NRows says how many rows are in the table; separators count
func (t *ATable) NRows() int {
	return len(t.rows)
}

// Headers returns the headers of a table, as Cells
// TODO v2: what if we have multiple header rows?  How does this change?
func (t *ATable) Headers() []Cell {
	if t.headerRow == nil {
		return nil
	}
	return t.headerRow.Cells()
}

// AddHeaders creates a header-row from the passed items and sets it
// as the table's header row.  The table is returned.
func (t *ATable) AddHeaders(items ...interface{}) Table {
	t.resizeColumnsAtLeast(len(items))
	hr := NewRowWithCapacity(len(items))
	hr.ErrorContainer = t.ErrorContainer
	for i := range items {
		hr.Add(NewCell(items[i]))
	}
	t.headerRow = hr

	invokePropertyCallbacks(t.tableRowAdditionCallbacks, CB_AT_ADD, hr, t.ErrorContainer)
	for i := range hr.cells {
		ptr := &hr.cells[i]
		if col := ptr.columnOfTable(); col != nil {
			invokePropertyCallbacks(col.cellCallbacks, CB_AT_ADD, ptr, t.ErrorContainer)
		}
		invokePropertyCallbacks(t.tableCellCallbacks, CB_AT_ADD, ptr, t.ErrorContainer)
	}
	return t
}

// AllRows returns an iterable of rows.
func (t *ATable) AllRows() []*Row {
	// We don't want callers using this to mess with the table bypassing the API;
	// they can still mess with cells via this but we at least reduce the risk
	// of row-reordering etc.  If there's a compelling argument, we can drop this
	// and just return the table's own t.rows
	rr := make([]*Row, len(t.rows))
	copy(rr, t.rows)
	return rr
}

// CellAt returns a pointer to the cell found at the given row and column
// coordinates, where the top-left item is 1,1.
func (t *ATable) CellAt(loc CellLocation) (*Cell, error) {
	if loc.Row < 1 || loc.Column < 1 || loc.Row > len(t.rows) {
		return nil, NoSuchCellError{Location: loc}
	}
	r := t.rows[loc.Row-1]
	if r.cells == nil || loc.Column > len(r.cells) {
		return nil, NoSuchCellError{Location: loc}
	}
	return &r.cells[loc.Column-1], nil
}
