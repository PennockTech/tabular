// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

/*
Markdown core does not support tables; for the most portable support,
assuming HTML is the final target, use the `html` sub-package instead.

Support for GitHub-Flavored Markdown's Tables is fairly widespread in
markdown packages.  That's the dialect which we speak at this time.

Note that the data model supported by GFM tables is very limited and is
likely to break on a number of valid inputs.  We attempt to deal
"safely" with this, to meet the security model, but it's likely that
there are inputs which break the integrity of the markdown table for
various renderers.  Again, the solution is to use the HTMLTable wrapper
instead.

So why have this sub-package?  In my experience, people working on
documentation like to be able to use the same input data and ask the
tools to generate markdown, which can then be included as-is, for output
in various formats.  Thus we have _two_ target audiences: a human reviewing
our output as it goes into .md documents, and people viewing the rendered
tables later.  We need to look "decent" for both, but can defer sanitization
to the human review step.
*/
package markdown // import "go.pennock.tech/tabular/markdown"

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/length"
	"go.pennock.tech/tabular/properties/align"
)

// A MarkdownTable wraps a tabular.Table to act as a render control for Markdown output.
type MarkdownTable struct {
	tabular.Table
}

// Wrap returns a MarkdownTable rendering object for the given tabular.Table.
func Wrap(t tabular.Table) *MarkdownTable {
	var ws widthSetter
	t.RegisterPropertyCallback(t, tabular.CB_AT_RENDER, tabular.CB_ON_CELL, ws)
	return &MarkdownTable{
		Table: t,
	}
}

// New returns a MarkdownTable with a new Table inside it, access via .Table
// or just use the interface methods on the MarkdownTable.
func New() *MarkdownTable {
	return Wrap(tabular.New())
}

// Render takes a tabular.Table and creates a default options MarkdownTable object
// and then calls the Render method upon it.
func Render(t tabular.Table) (string, error) {
	return Wrap(t).Render()
}

// RenderTo takes a tabular.Table and creates a default options MarkdownTable object
// and calls the RenderTo method upon it.
func RenderTo(t tabular.Table, w io.Writer) error {
	return Wrap(t).RenderTo(w)
}

// Render takes a tabular Table and returns a string representing the fully
// rendered table or an error.
func (mt *MarkdownTable) Render() (string, error) {
	b := &bytes.Buffer{}
	err := mt.RenderTo(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderTo writes the table to the provided writer, stopping if it encounters
// an error.
func (mt *MarkdownTable) RenderTo(w io.Writer) error {
	mt.InvokeRenderCallbacks()
	var err error
	columnCount := mt.NColumns()
	if columnCount < 1 {
		return fmt.Errorf("markdown:RenderTo: can't emit a table with %d columns", columnCount)
	}

	// We have headers, then a control line which affects alignment, then the body.
	// No titles, that's not currently in the core tabular model.
	// We do want to try to align columns, for prettyness in markdown to be edited
	// by humans, but are willing to give up on alignment in degenerate cases (embedded newlines).

	headers := mt.Headers()
	if headers == nil {
		return fmt.Errorf("markdown:RenderTo: can't emit a table without headers")
	}

	widths := make([]int, columnCount)
	if len(headers) > columnCount {
		return fmt.Errorf("structural bug, columnCount %d but %d headers", columnCount, len(headers))
	}
	for i := range headers {
		widths[i] = CellPropertyExtractWidth(&headers[i])
	}
	for n, r := range mt.AllRows() {
		if r.IsSeparator() {
			continue
		}
		cells := r.Cells()
		if len(cells) > columnCount {
			return fmt.Errorf("structural bug, columnCount %d but %d cells in row %d", columnCount, len(cells), n)
		}
		for i := range cells {
			w := CellPropertyExtractWidth(&cells[i])
			if w > widths[i] {
				widths[i] = w
			}
		}
	}

	controlRowCells := make([]tabular.Cell, 0, columnCount)
	alignments := make([]align.Alignment, columnCount)
	for i := 0; i < columnCount; i++ {
		width := widths[i]
		// spec mandates at least three dashes
		if width < 3 {
			width = 3
		}
		var al align.Alignment
		alRaw := mt.Column(i + 1).GetProperty(align.PropertyType)
		if alRaw != nil {
			al = alRaw.(align.Alignment)
			alignments[i] = al
		}
		var content string
		switch al {
		case nil, align.Left:
			content = " " + strings.Repeat("-", width) + " "
		case align.Right:
			content = " " + strings.Repeat("-", width) + ":"
		case align.Center:
			content = ":" + strings.Repeat("-", width) + ":"
		default:
			content = " " + strings.Repeat("-", width) + " "
		}

		controlRowCells = append(controlRowCells, tabular.NewCell(content))
	}

	if err = mt.emitRow(w, columnCount, headers, widths, alignments, true); err != nil {
		return err
	}

	if err = mt.emitRow(w, columnCount, controlRowCells, widths, alignments, false); err != nil {
		return err
	}

	for _, r := range mt.AllRows() {
		if r.IsSeparator() {
			continue
		}
		if err = mt.emitRow(w, columnCount, r.Cells(), widths, alignments, true); err != nil {
			return err
		}
	}
	return nil
}

// emitRow handles just one row, whether from headers or body.
// It needs to know how many columns should be in the row, so that it can add extras,
// or error out, as needed.
func (mt *MarkdownTable) emitRow(
	w io.Writer,
	columnCount int,
	cells []tabular.Cell,
	widths []int,
	alignments []align.Alignment,
	addPads bool,
) error {
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
	//
	// For alignment, note that escaping will completely throw things off.  That's okay.
	// We don't align right for escaped content.  We're after "close enough to not be jarring".
	barLeft := "| "
	barRight := " |"
	barCenter := " | "
	if !addPads {
		barLeft = "|"
		barRight = "|"
		barCenter = "|"
	}
	io.WriteString(w, barLeft)
	for i = 0; i < max-1; i++ {
		if _, err := fmt.Fprint(w, mt.mdPaddedCellEscape(cells, widths, alignments, i), barCenter); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, mt.mdPaddedCellEscape(cells, widths, alignments, i), barRight); err != nil {
		return err
	}
	i++
	for ; i < columnCount; i++ {
		// these are the extra columns, always have one whitespace before bar
		if _, err := io.WriteString(w, " |"); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}
	return nil
}

// mdPaddedCellEscape is a wrapper around mdCellEscape which takes care of padding, widths, etc
// (and so needs more parameters)
func (mt *MarkdownTable) mdPaddedCellEscape(
	cells []tabular.Cell,
	widths []int,
	alignments []align.Alignment,
	i int,
) string {
	baseline := mt.mdCellEscape(cells[i].String())
	wantWidth := widths[i]
	haveWidth := length.StringCells(baseline)
	if haveWidth >= wantWidth {
		return baseline
	}
	pad := wantWidth - haveWidth
	switch alignments[i] {
	case nil, align.Left:
		return baseline + strings.Repeat(" ", pad)
	case align.Right:
		return strings.Repeat(" ", pad) + baseline
	case align.Center:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + baseline + strings.Repeat(" ", right)
	default:
		return baseline + strings.Repeat(" ", pad)
	}
}

// mdCellEscape handles producing the output for one field, with surrounding
// quotes.  Note that there are two separate issues here:
//  1. Use of the pipe character as a column separator
//  2. Use of HTML!  Our security model does not trust the content within the
//     table, so we should escape everything.  If there's a use-case for more
//     trusted content, that should be a non-default option which we can add
//     later.
//
// What about use of other markdown?  For now, we're going with "anything
// actively dangerous should require dropping to HTML to accomplish".  If I'm
// wrong about this, please file a bug report!
//
// For multi-line content, our approach is to replace newlines with HTML escape
// sequences.  Markdown tables do not support multiple "physical" lines in one
// cell, so this seems the only way.  We ignore CR, only handling LF.
// Really, at this point, we are past the limits of the spec.
func (mt *MarkdownTable) mdCellEscape(in string) string {
	return strings.Replace(strings.Replace(html.EscapeString(in), "|", "&#x7c;", -1), "\n", "&#x0a;", -1)
}
