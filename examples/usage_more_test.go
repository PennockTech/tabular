// Copyright © 2018 Pennock Tech, LLC.
// These examples are deliberately for use in your code, and so while this file
// is still covered under LICENSE.txt, only the liability disclaimers apply.
// Copy/paste freely without the license applying directly to your code.  The
// license still applies to any bundled version of this library.

package examples

import (
	"fmt"
	"os"

	"go.pennock.tech/tabular"
	auto_table "go.pennock.tech/tabular/auto"
	"go.pennock.tech/tabular/texttable/decoration"
)

func Example_inspection() {
	t := auto_table.New("") // empty string is style-less, so a rendering error later
	t = populateTable(t)

	err := t.RenderTo(os.Stdout)
	if err != nil {
		// nb, we actually expect output, so using os.Stdout here this time
		fmt.Fprintf(os.Stdout, "table rendering error: %s\n", err)
	}
	fmt.Printf("Table has %d rows and %d columns\n", t.NRows(), t.NColumns())
	row := 3
	column := 2
	cell, err := t.CellAt(tabular.CellLocation{Row: row, Column: column})
	fmt.Printf("At row %d column %d the value is %q\n", row, column, cell)

	fmt.Printf("\nManual style at render time:\n")
	if err = auto_table.RenderTo(t, os.Stdout, "utf8-double"); err != nil {
		fmt.Fprintf(os.Stdout, "second render, table rendering error: %s\n", err)
	}
	// nb: this highlights a limitation of tabular: we can only use characters
	// which exist in Unicode, and the box-drawing double-lines are missing
	// double left/up/right and single down, so you'll notice a display glitch
	// below.

	// Output:
	// table rendering error: table has no decoration at all, can't render
	// Table has 4 rows and 2 columns
	// At row 3 column 2 the value is "2.22"
	//
	// Manual style at render time:
	// ╔══════════════╦════════════╗
	// ║ Fruit        ║ Price $/lb ║
	// ╠══════════════╪════════════╣
	// ║ bananas      │ 0.56       ║
	// ║ apples       │ 1.31       ║
	// ║ strawberries │ 2.22       ║
	// ║ pears        │ 1.90       ║
	// ╚══════════════╧════════════╝
}

func Example_rows() {
	// can use named constants so that typos can be caught at compile time, instead
	// of strings, _if_ you're willing to import the extra package; we don't force
	// callers to do this, letting them make the call on trade-offs.
	t := auto_table.New(decoration.D_UTF8_LIGHT_CURVED)
	t = populateTable(t)

	extra := t.NewRowSizedFor()
	// could call tabular.NewRow(), or tabular.NewRowWithCapacity(2), but the
	// table-method supplies the current table column count to the latter,
	// letting us just get the right size.  Can always add cells to a row, and
	// keep adding, which we'll show here, but this is getting the "right" size
	// with pre-allocation.
	c1 := tabular.NewCell("tomato")
	c2 := tabular.NewCell("2.24")
	extra = extra.Add(c1).Add(c2).Add(tabular.NewCell("but is it a fruit?"))

	t.AddSeparator()
	t.AddRow(extra)

	fmt.Printf("Now have %d columns\n", t.NColumns())
	if err := t.RenderTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
	}

	// Output:
	// Now have 3 columns
	// ╭──────────────┬────────────┬────────────────────╮
	// │ Fruit        │ Price $/lb │                    │
	// ├──────────────┼────────────┼────────────────────┤
	// │ bananas      │ 0.56       │                    │
	// │ apples       │ 1.31       │                    │
	// │ strawberries │ 2.22       │                    │
	// │ pears        │ 1.90       │                    │
	// ├──────────────┼────────────┼────────────────────┤
	// │ tomato       │ 2.24       │ but is it a fruit? │
	// ╰──────────────┴────────────┴────────────────────╯
}
