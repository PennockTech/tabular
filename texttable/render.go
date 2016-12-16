// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable // import "go.pennock.tech/tabular/texttable"

import (
	"bytes"
	"errors"
	"io"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/texttable/decoration"
)

// Render takes a tabular.Table and creates a default options TextTable object
// and then calls the Render method upon it.
func Render(t tabular.Table) (string, error) {
	return Wrap(t).Render()
}

// Render returns a string representing the fully rendered table, or an error.
func (t *TextTable) Render() (string, error) {
	b := &bytes.Buffer{}
	err := t.RenderTo(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderTo takes a tabular.Table and creates a default options TextTable object
// and then calls the RenderTo method upon it.
func RenderTo(t tabular.Table, w io.Writer) error {
	return Wrap(t).RenderTo(w)
}

// RenderTo writes the TextTable to the provided writer, stopping if it
// encounters an error.
func (t *TextTable) RenderTo(w io.Writer) error {
	if t.decor == decoration.EmptyDecoration {
		return errors.New("table has no decoration at all, can't render")
	}

	t.InvokeRenderCallbacks()

	columnCount := t.NColumns()
	headers := t.Headers() // may be nil

	columnWidths := make([]int, columnCount)
	if headers != nil {
		for i := range columnWidths {
			columnWidths[i] = CellPropertyExtractDimensions(&headers[i]).cellWidth
		}
	}
	for _, row := range t.AllRows() {
		if row.IsSeparator() {
			continue
		}
		for i, cell := range row.Cells() {
			if i > columnCount {
				break
			}
			d := CellPropertyExtractDimensions(&cell)
			if d.cellWidth > columnWidths[i] {
				columnWidths[i] = d.cellWidth
			}
		}
	}

	emitter := t.decor.ForColumnWidths(columnWidths)
	emitter.SetEOL("\n")

	if headers != nil {
		if _, err := io.WriteString(w, emitter.LineHeaderTop()); err != nil {
			return err
		}
		for _, lineParts := range t.RowToLinesOfWidthStrings(headers, columnCount) {
			if _, err := io.WriteString(w, emitter.HeaderLineRendered(lineParts)); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, emitter.LineHeaderBodySep()); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w, emitter.LineBodyTop()); err != nil {
			return err
		}
	}
	for _, row := range t.AllRows() {
		// do we want a channel returning rows instead?  I don't _think_ so
		if row.IsSeparator() {
			if _, err := io.WriteString(w, emitter.LineSeparator()); err != nil {
				return err
			}
			continue
		}
		for _, lineParts := range t.RowToLinesOfWidthStrings(row.Cells(), columnCount) {
			if _, err := io.WriteString(w, emitter.BodyLineRendered(lineParts)); err != nil {
				return err
			}
		}
	}
	if _, err := io.WriteString(w, emitter.LineBottom()); err != nil {
		return err
	}
	// do _not_ try to close the writer, that's not ours
	return nil
}

func (t *TextTable) RowToLinesOfWidthStrings(
	cells []tabular.Cell,
	columnCount int,
) [][]decoration.WidthString {
	max := len(cells)
	if columnCount < max {
		max = columnCount
	}
	columns := make([][]decoration.WidthString, max)
	lineCount := 1
	for i := 0; i < max; i++ {
		columns[i] = CellPropertyExtractLinesWidths(&cells[i])
		if len(columns[i]) > lineCount {
			lineCount = len(columns[i])
		}
	}

	lines := make([][]decoration.WidthString, lineCount)
	for l := 0; l < lineCount; l++ {
		lines[l] = make([]decoration.WidthString, columnCount)
		for c := 0; c < max; c++ {
			if l >= len(columns[c]) {
				lines[l][c] = decoration.WidthString{}
			} else {
				lines[l][c] = columns[c][l]
			}
		}
	}

	return lines
}
