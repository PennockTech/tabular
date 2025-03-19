// Copyright © 2016,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package csv // import "go.pennock.tech/tabular/csv"

import (
	"bytes"
	"fmt"
	"io"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/properties"
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
	var (
		err          error
		defaultOmit  bool
		omittedCount int
		omitColumns  []bool
		skipRow      bool
	)

	displayColumnCount := ct.NColumns()
	if displayColumnCount < 1 {
		return fmt.Errorf("csv:RenderTo: can't emit a table with %d columns", displayColumnCount)
	}

	if defaultOmit, err = properties.ExpectBoolPropertyOrNil(
		properties.Omit,
		ct.Column(0).GetProperty(properties.Omit),
		"csv:RenderTo", "default column", 0); err != nil {
		return err
	}
	omitColumns = make([]bool, displayColumnCount)
	for i := range displayColumnCount {
		omit := ct.Column(i + 1).GetProperty(properties.Omit)
		if omit != nil {
			if omitColumns[i], err = properties.ExpectBoolPropertyOrNil(properties.Omit, omit, "csv:RenderTo", "column", i+1); err != nil {
				return err
			}
		} else {
			omitColumns[i] = defaultOmit
		}
		if omitColumns[i] {
			omittedCount++
		}
	}
	if omittedCount == displayColumnCount {
		return fmt.Errorf("csv:RenderTo: can't emit a table with all columns omitted")
	}
	displayColumnCount -= omittedCount

	headers := ct.Headers()
	if headers != nil {
		if err = ct.emitRow(w, displayColumnCount, omitColumns, headers); err != nil {
			return err
		}
	}

	for rowNum, r := range ct.AllRows() {
		if skipRow, err = properties.ExpectBoolPropertyOrNil(properties.Omit, r.GetProperty(properties.Omit), "text:renderTo", "row", rowNum+1); err != nil {
			return err
		}
		if skipRow {
			continue
		}
		if r.IsSeparator() {
			continue
		}
		if err = ct.emitRow(w, displayColumnCount, omitColumns, r.Cells()); err != nil {
			return err
		}
	}
	return nil
}

// emitRow handles just one row, whether from headers or body.
// It needs to know how many columns should be in the row, so that it can add extras,
// or error out, as needed.  The displayColumnCount should be "how many to actually show",
// pre-adjusted for columns being omitted, and is used so that "short" rows get
// extra columns added to make up the difference.
func (ct *CSVTable) emitRow(w io.Writer, displayColumnCount int, omitColumns []bool, cells []tabular.Cell) error {
	// RFC4180 §2 item 7 states newlines are allowed _within_ a virtual row of
	// CSV, as long as enclosed within double-quotes.  Wikipedia says _some_ variants
	// use backslash-escaping instead of doubling for quotes, but it seems that
	// \n is not well understood.
	// Alas, or we coult just use fmt.Sprintf("%#v", str) to convert.
	// We're only a generator, don't need to handle everything, so we just follow
	// what spec exists and allow newlines within a field.  We quote _all_ fields.

	var (
		i     int
		err   error
		shown bool
	)

	// line.Write() is guaranteed to never return an error (it will panic instead)
	var line bytes.Buffer
	line.Grow(256)

	// We shouldn't have more cells than columns, but we do in the regression
	// tests to make sure we cleanly handle a broken table with an error instead of a panic.
	if len(cells) > len(omitColumns) {
		return fmt.Errorf("structural bug, max_columns %d but this row has %d cells", len(omitColumns), len(cells))
	}

	// Game-plan:
	// 1. Track once we've printed a column, show a leading separator before subsequent
	// 2. Print all columns we have
	// 3. if too few columns in this row, repeatedly add empty extra columns
	// if too many columns in this row, should have errored out above
	// if only one column, the first repeated print should be skipped
	shown = false
	for i = range len(cells) {
		if omitColumns[i] {
			continue
		}
		if shown {
			line.WriteString(ct.fieldSeparator)
		}
		line.WriteString(ct.csvEscape(cells[i].String()))
		shown = true
	}
	for i++; i < displayColumnCount; i++ {
		if shown {
			line.WriteString(ct.fieldSeparator)
		}
		line.WriteString("\"\"")
		shown = true
	}
	line.WriteRune('\n')
	if _, err = w.Write(line.Bytes()); err != nil {
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
