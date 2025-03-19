// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"bytes"
	"fmt"
)

func debugCallbackSetCount(cs *callbackSet) int {
	if cs == nil {
		return 0
	}
	return len(cs.addTime) + len(cs.preCellRenderTime) + len(cs.renderTime) + len(cs.postCellRenderTime)
}

func (t *ATable) GoString() string {
	buf := &bytes.Buffer{}
	errCount := 0
	if t.ErrorContainer != nil {
		errCount = len(t.ErrorContainer.errors_)
	}
	cbCount := (debugCallbackSetCount(&t.tableItselfCallbacks) +
		debugCallbackSetCount(&t.tableCellCallbacks) +
		debugCallbackSetCount(&t.tableRowAdditionCallbacks))

	fmt.Fprintf(buf, "*ATable(%d errors, %d col, %d row, %d cbs)",
		errCount, t.NColumns(), t.NRows(), cbCount)
	if t.propertyImpl.properties != nil {
		fmt.Fprintf(buf, ".Props{%#v}", t.propertyImpl.properties)
	}

	fmt.Fprint(buf, ".Columns{")
	for i := range t.columns {
		if i > 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, "C(%d, %q, %dcbs)", i, t.columns[i].Name,
			(debugCallbackSetCount(&t.columns[i].cellCallbacks) +
				debugCallbackSetCount(&t.columns[i].columnItselfCallbacks)))
		if t.columns[i].propertyImpl.properties != nil {
			fmt.Fprintf(buf, ".Props{%#v}", t.columns[i].propertyImpl.properties)
		}
	}

	if t.headerRow == nil {
		fmt.Fprint(buf, "}.NoHeaders.Body[\n")
	} else {
		fmt.Fprint(buf, "}.HeaderRow{\n")
		fmt.Fprintf(buf, " %#v\n", t.headerRow)
		fmt.Fprint(buf, "}.Body[\n")
	}
	for _, row := range t.rows {
		fmt.Fprintf(buf, " %#v,\n", row)
	}
	fmt.Fprint(buf, "]")

	el := t.Errors()
	if el == nil {
		fmt.Fprint(buf, ".NoErrors☺")
	} else {
		fmt.Fprint(buf, ".Errors[\n")
		for i := range el {
			fmt.Fprintf(buf, "   %q,\n", el[i])
		}
		fmt.Fprint(buf, "]")
	}

	return buf.String()
}

func (row *Row) GoString() string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "R(n%d, %d cbs)",
		row.rowNum,
		(debugCallbackSetCount(&row.rowCellCallbacks) +
			debugCallbackSetCount(&row.rowItselfCallbacks)))
	if row.propertyImpl.properties != nil {
		fmt.Fprintf(buf, ".Props{%#v}", row.propertyImpl.properties)
	}
	if row.isSeparator {
		fmt.Fprint(buf, ".SEP")
		return buf.String()
	}
	fmt.Fprint(buf, ".Cells{")
	for i := range row.cells {
		if i > 0 {
			fmt.Fprint(buf, ", ")
		}
		fmt.Fprintf(buf, "%#v", &row.cells[i])
	}
	fmt.Fprint(buf, "}")

	return buf.String()
}

func (cell *Cell) GoString() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "C(%q, %d cbs)", cell.String(), debugCallbackSetCount(&cell.callbacks))
	if cell.propertyImpl.properties != nil {
		fmt.Fprintf(buf, ".Props{%#v}\n\t", cell.propertyImpl.properties)
	}
	return buf.String()
}
