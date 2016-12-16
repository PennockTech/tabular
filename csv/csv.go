// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package csv // import "go.pennock.tech/tabular/csv"

import (
	"bytes"
	"fmt"
	"io"

	"go.pennock.tech/tabular"
)

// A CSVTable wraps a tabular.Table to act as a render control for CSV output.
type CSVTable struct {
	tabular.Table

	fieldSeparator string
	// TODO: any output style controls here, to deviate from RFC4180 (eg,
	// tab-output, only-quote-if-needed, other-escaping.
}

// Wrap returns a CSVTable rendering object for the given tabular.Table.
func Wrap(t tabular.Table) *CSVTable {
	return &CSVTable{
		Table:          t,
		fieldSeparator: ",",
	}
}

// New returns a CSVTable with a new Table inside it, access via .Table
// or just use the interface methods on the CSVTable.
func New() *CSVTable {
	return Wrap(tabular.New())
}

// Render takes a tabular.Table and creates a default options CSVTable object
// and then calls the Render method upon it.
func Render(t tabular.Table) (string, error) {
	return Wrap(t).Render()
}

// RenderTo takes a tabular.Table and creates a default options CSVTable object
// and calls the RenderTo method upon it.
func RenderTo(t tabular.Table, w io.Writer) error {
	return Wrap(t).RenderTo(w)
}

// Render takes a tabular Table and returns a string representing the fully
// rendered table or an error.
func (ct *CSVTable) Render() (string, error) {
	b := &bytes.Buffer{}
	err := ct.RenderTo(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderTo writes the table to the provided writer, stopping if it encounters
// an error.
func (ct *CSVTable) RenderTo(w io.Writer) error {
	ct.InvokeRenderCallbacks()
	var err error
	columnCount := ct.NColumns()
	if columnCount < 1 {
		return fmt.Errorf("csv:RenderTo: can't emit a table with %d columns", columnCount)
	}
	headers := ct.Headers()
	if headers != nil {
		if err = ct.emitRow(w, columnCount, headers); err != nil {
			return err
		}
	}
	for _, r := range ct.AllRows() {
		if r.IsSeparator() {
			continue
		}
		if err = ct.emitRow(w, columnCount, r.Cells()); err != nil {
			return err
		}
	}
	return nil
}

// emitRow handles just one row, whether from headers or body.
// It needs to know how many columns should be in the row, so that it can add extras,
// or error out, as needed.
func (ct *CSVTable) emitRow(w io.Writer, columnCount int, cells []tabular.Cell) error {
	// RFC4180 §2 item 7 states newlines are allowed _within_ a virtual row of
	// CSV, as long as enclosed within double-quotes.  Wikipedia says _some_ variants
	// use backslash-escaping instead of doubling for quotes, but it seems that
	// \n is not well understood.
	// Alas, or we coult just use fmt.Sprintf("%#v", str) to convert.
	// We're only a generator, don't need to handle everything, so we just follow
	// what spec exists and allow newlines within a field.  We quote _all_ fields.
	var i int
	var max int = len(cells)
	if columnCount < max {
		return fmt.Errorf("structural bug, columnCount %d but %d cells", columnCount, max)
	}
	// Game-plan:
	// 1. repeatedly print all-but-last available column with trailing separator
	// 2. print last column, no separator
	// 3. if too few columns in this row, repeatedly add leading separator and extra column
	// if too many columns in this row, should have errored out above
	// if only one column, the first repeated print should be skipped
	for i = 0; i < max-1; i++ {
		if _, err := fmt.Fprint(w, ct.csvEscape(cells[i].String()), ct.fieldSeparator); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, ct.csvEscape(cells[i].String())); err != nil {
		return err
	}
	i++
	for ; i < columnCount; i++ {
		if _, err := fmt.Fprint(w, ct.fieldSeparator, "\"\""); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

// csvEscape handles producing the output for one field, with surrounding
// quotes, escaping the contents as needed.  It's a method even though at time
// of writing the object isn't referenced, so that future options can change
// the escaping style (eg, newlines as escape sequence instead of raw).
func (ct *CSVTable) csvEscape(in string) string {
	// We can ignore Unicode by just matching literally on the quote character.
	b := make([]byte, len(in)*2+2)
	b[0] = '"'
	j := 1
	i := 0
	max := len(in)
	for ; i < max; i++ {
		b[j] = in[i]
		j++
		if in[i] == '"' {
			b[j] = '"'
			j++
		}
	}
	b[j] = '"'
	return string(b[:j+1])
}
