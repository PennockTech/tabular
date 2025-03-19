// Copyright © 2018,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package json // import "go.pennock.tech/tabular/json"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/properties"
)

// A JSONTable wraps a tabular.Table to act as a render control for JSON output.
type JSONTable struct {
	tabular.Table
}

// Wrap returns a JSONTable rendering object for the given tabular.Table.
func Wrap(t tabular.Table) *JSONTable {
	return &JSONTable{
		Table: t,
	}
}

// New returns a JSONTable with a new Table inside it, access via .Table
// or just use the interface methods on the JSONTable.
func New() *JSONTable {
	return Wrap(tabular.New())
}

// Render takes a tabular.Table and creates a default options JSONTable object
// and then calls the Render method upon it.
func Render(t tabular.Table) (string, error) {
	return Wrap(t).Render()
}

// RenderTo takes a tabular.Table and creates a default options JSONTable object
// and calls the RenderTo method upon it.
func RenderTo(t tabular.Table, w io.Writer) error {
	return Wrap(t).RenderTo(w)
}

// Render takes a tabular Table and returns a string representing the fully
// rendered table or an error.
func (jt *JSONTable) Render() (string, error) {
	b := &bytes.Buffer{}
	err := jt.RenderTo(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderTo writes the table to the provided writer, stopping if it encounters
// an error.
func (jt *JSONTable) RenderTo(w io.Writer) error {
	jt.InvokeRenderCallbacks()
	var err error
	columnCount := jt.NColumns()
	if columnCount < 1 {
		return fmt.Errorf("json:RenderTo: can't emit a table with %d columns", columnCount)
	}

	var defaultSkipable, defaultOmit bool
	if defaultSkipable, err = properties.ExpectBoolPropertyOrNil(
		properties.Skipable,
		jt.Column(0).GetProperty(properties.Skipable),
		"json:RenderTo", "default column", 0); err != nil {
		return err
	}
	if defaultOmit, err = properties.ExpectBoolPropertyOrNil(
		properties.Omit,
		jt.Column(0).GetProperty(properties.Omit),
		"json:RenderTo", "default column", 0); err != nil {
		return err
	}

	skipableColumns := make([]bool, columnCount)
	omitColumns := make([]bool, columnCount)
	keys := make([][]byte, columnCount)
	seen := make(map[string]int, columnCount)
	headers := jt.Headers()
	if headers == nil {
		return fmt.Errorf("json:RenderTo: require headers for JSON rendering to provide keys")
	}
	if len(headers) < columnCount {
		return fmt.Errorf("json:RenderTo: require %d headers for keys, only found %d", columnCount, len(headers))
	}
	for i := range columnCount {
		s := headers[i].String()
		if s == "" {
			return fmt.Errorf("json:RenderTo: column %d has an empty header, unusable as a key", i+1)
		}
		if previous, already := seen[s]; already {
			return fmt.Errorf("json:RenderTo: column %d header matches previous column %d: %q", i+1, previous, s)
		}
		seen[s] = i
		t, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("json:RenderTo: column %d header JSON encoding failure: %s", i+1, err)
		}
		keys[i] = append(t, byte(':'), byte(' '))

		c := jt.Column(i + 1)
		sk := c.GetProperty(properties.Skipable)
		if sk != nil {
			if skipableColumns[i], err = properties.ExpectBoolPropertyOrNil(properties.Skipable, sk, "json:RenderTo", "column", i+1); err != nil {
				return err
			}
		} else {
			skipableColumns[i] = defaultSkipable
		}
		omit := c.GetProperty(properties.Omit)
		if omit != nil {
			if omitColumns[i], err = properties.ExpectBoolPropertyOrNil(properties.Omit, omit, "json:RenderTo", "column", i+1); err != nil {
				return err
			}
		} else {
			omitColumns[i] = defaultOmit
		}
	}

	if _, err = io.WriteString(w, "[\n"); err != nil {
		return err
	}
	var skipRow bool
	needComma := false
	for rowNum, r := range jt.AllRows() {
		if skipRow, err = properties.ExpectBoolPropertyOrNil(properties.Omit, r.GetProperty(properties.Omit), "text:renderTo", "row", rowNum+1); err != nil {
			return err
		}
		if skipRow {
			continue
		}
		if needComma {
			if _, err = io.WriteString(w, ",\n"); err != nil {
				return err
			}
			needComma = false
		}
		if r.IsSeparator() {
			if _, err = io.WriteString(w, "\n"); err != nil {
				return err
			}
			continue
		}
		if err = jt.emitRowAsJSONObject(w, skipableColumns, omitColumns, keys, r.Cells()); err != nil {
			return err
		}
		needComma = true
	}
	// We assume need newline prefix because no comma+newline from new row,
	// but if the table is empty, this will result in "[\n\n]\n" which is
	// slightly ugly.  But valid.  So live with it.
	if _, err = io.WriteString(w, "\n]\n"); err != nil {
		return err
	}
	return nil
}

// emitRowAsJSONObject handles just one row, as a JSON object, it does not handle
// any trailing commas outside the object, separating it from the next.
func (jt *JSONTable) emitRowAsJSONObject(w io.Writer, skipableColumns []bool, omitColumns []bool, keys [][]byte, cells []tabular.Cell) error {
	var (
		i, max int
		err    error
	)
	max = len(cells)
	if len(keys) < max {
		return fmt.Errorf("structural bug, %d headers but %d cells", len(keys), max)
	}

	separator := "{"

	for i = range max {
		if omitColumns[i] {
			continue
		}
		if skipableColumns[i] && cells[i].Empty() {
			continue
		}
		if _, err = io.WriteString(w, separator); err != nil {
			return err
		}
		separator = ", "

		if _, err = w.Write(keys[i]); err != nil {
			return err
		}

		// We call .String, so that updateCache stuff is done (workaround for
		// no update-cache in .Item), and so that we have a fallback for when
		// the JSON marshalling returns an empty struct: our callers only have
		// to set a String() method, but json.Marshal doesn't use that as a
		// marshalling method.  If we rework our API, then we can suggest that
		// cell data types have MarshalText() method.
		fallback := cells[i].String()
		t, err := json.Marshal(cells[i].Item())
		if err != nil {
			return fmt.Errorf("json:RenderTo: column %d header JSON encoding failure: %s", i+1, err)
		}
		if bytes.Equal(t, []byte("{}")) && fallback != "" {
			t, err = json.Marshal(fallback)
			if err != nil {
				return fmt.Errorf("json:RenderTo: column %d header JSON encoding of text fallback failure: %s", i+1, err)
			}
		}
		if _, err = w.Write(t); err != nil {
			return err
		}
	}

	if separator == "{" {
		// never printed the opening brace, all fields skipable
		if _, err = io.WriteString(w, "{}"); err != nil {
			return err
		}
		return nil
	}
	if _, err = io.WriteString(w, "}"); err != nil {
		return err
	}
	return nil
}
