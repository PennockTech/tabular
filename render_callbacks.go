// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// InvokeRenderCallbacks is used by rendering packages to trigger a table-wide
// invocation of render-time callbacks.
//
// We first invoke pre-cell callbacks going: table->column->row->cell
// We then invoke regular render callbacks cell->row->column->table
func (t *ATable) InvokeRenderCallbacks() {
	ec := t.ErrorContainer
	invokePropertyCallbacks(t.tableItselfCallbacks, CB_AT_RENDER_PRECELL, t, ec)
	for _, col := range t.columns {
		invokePropertyCallbacks(col.columnItselfCallbacks, CB_AT_RENDER_PRECELL, &col, ec)
	}
	if t.headerRow != nil {
		t.headerRow.invokeRenderCallbacks(t, ec)
	}
	for _, row := range t.rows {
		row.invokeRenderCallbacks(t, ec)
	}
	for _, col := range t.columns {
		invokePropertyCallbacks(col.columnItselfCallbacks, CB_AT_RENDER_POSTCELL, &col, ec)
	}
	invokePropertyCallbacks(t.tableItselfCallbacks, CB_AT_RENDER_POSTCELL, t, ec)
}

func (row *Row) invokeRenderCallbacks(t *ATable, ec ErrorReceiver) {
	invokePropertyCallbacks(row.rowItselfCallbacks, CB_AT_RENDER_PRECELL, row, ec)
	for i := range row.cells {
		ptr := &row.cells[i]
		col := ptr.columnOfTable()
		invokePropertyCallbacks(t.tableCellCallbacks, CB_AT_RENDER_PRECELL, ptr, ec)
		if col != nil {
			invokePropertyCallbacks(col.cellCallbacks, CB_AT_RENDER_PRECELL, ptr, row.ErrorContainer)
		}
		invokePropertyCallbacks(row.rowCellCallbacks, CB_AT_RENDER_PRECELL, ptr, ec)

		invokePropertyCallbacks(t.tableCellCallbacks, CB_AT_RENDER, ptr, ec)
		invokePropertyCallbacks(ptr.callbacks, CB_AT_RENDER, ptr, ec)

		invokePropertyCallbacks(row.rowCellCallbacks, CB_AT_RENDER_POSTCELL, ptr, ec)
		if col != nil {
			invokePropertyCallbacks(col.cellCallbacks, CB_AT_RENDER_POSTCELL, ptr, row.ErrorContainer)
		}

	}
	invokePropertyCallbacks(row.rowItselfCallbacks, CB_AT_RENDER_POSTCELL, row, ec)

}
