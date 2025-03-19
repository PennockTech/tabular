// Copyright © 2018,2025 Pennock Tech, LLC.
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
	tab_csv "go.pennock.tech/tabular/csv"
	tab_html "go.pennock.tech/tabular/html"
	tab_json "go.pennock.tech/tabular/json"
	tab_md "go.pennock.tech/tabular/markdown"
	"go.pennock.tech/tabular/properties"
	"go.pennock.tech/tabular/properties/align"
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
	loc := cell.Location()
	fmt.Printf("At row %d column %d the value is %q, self-described as row %d column %d\n",
		row, column, cell, loc.Row, loc.Column)

	rowKiwi := tabular.NewRowWithCapacity(2)
	rowKiwi.Add(tabular.NewCell("kiwis"))
	rowKiwi.Add(tabular.NewCell("2.50"))
	t.AddRow(rowKiwi)
	kiwiLocation := rowKiwi.Location()
	fmt.Printf("Kiwi row is row %d\n", kiwiLocation.Row)

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
	// At row 3 column 2 the value is "2.22", self-described as row 3 column 2
	// Kiwi row is row 5
	//
	// Manual style at render time:
	// ╔══════════════╦════════════╗
	// ║ Fruit        ║ Price $/lb ║
	// ╠══════════════╪════════════╣
	// ║ bananas      │ 0.56       ║
	// ║ apples       │ 1.31       ║
	// ║ strawberries │ 2.22       ║
	// ║ pears        │ 1.90       ║
	// ║ kiwis        │ 2.50       ║
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

func Example_alignment() {
	t := auto_table.New(decoration.D_UTF8_LIGHT_CURVED)
	t = populateTable(t)
	t.AddRowItems("Sekai-ichi apple", "10.50")

	t.Column(2).SetProperty(align.PropertyType, align.Right)

	if err := t.RenderTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
	}
	// Output:
	// ╭──────────────────┬────────────╮
	// │ Fruit            │ Price $/lb │
	// ├──────────────────┼────────────┤
	// │ bananas          │       0.56 │
	// │ apples           │       1.31 │
	// │ strawberries     │       2.22 │
	// │ pears            │       1.90 │
	// │ Sekai-ichi apple │      10.50 │
	// ╰──────────────────┴────────────╯
}

func Example_alignment2() {
	t := auto_table.New(decoration.D_UTF8_LIGHT_CURVED)
	t = populateTable(t)
	t.AddRowItems("Sekai-ichi apple", "10.50")

	t.Column(0).SetProperty(align.PropertyType, align.Right)
	t.Column(1).SetProperty(align.PropertyType, align.Center)

	if err := t.RenderTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
	}
	// Output:
	// ╭──────────────────┬────────────╮
	// │      Fruit       │ Price $/lb │
	// ├──────────────────┼────────────┤
	// │     bananas      │       0.56 │
	// │      apples      │       1.31 │
	// │   strawberries   │       2.22 │
	// │      pears       │       1.90 │
	// │ Sekai-ichi apple │      10.50 │
	// ╰──────────────────┴────────────╯
}

func populateFourTable(t tabular.Table) {
	t.AddHeaders("Person", "Age", "Score", "Color")
	t.AddRowItems("Fred", 34, 0.6, "red")
	t.AddRowItems("Gladys", 32, 0.8, "yellow")
	t.AddRowItems("Bert", 57, 0.7, "green")
	t.AddRowItems("Belinda", 58, 0.8, "blue")

	for _, err := range t.Errors() {
		fmt.Fprintf(os.Stderr, "table error: %s\n", err)
	}

	t.Column(2).SetProperty(align.PropertyType, align.Right)
	// FIXME: add Column(3) align-float dot alignment
}

func Example_omit_text() {
	t := auto_table.New(decoration.D_UTF8_LIGHT_CURVED)
	populateFourTable(t)

	show := func(tb auto_table.RenderTable) {
		if err := tb.RenderTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	show(t)
	t.Column(3).SetProperty(properties.Omit, true)
	show(t)
	t.Column(3).SetProperty(properties.Omit, false) // explicitly false
	show(t)
	t.Column(3).SetProperty(properties.Omit, nil) // nil removes the property
	show(t)

	t.Column(3).SetProperty(properties.Omit, true)
	t.Column(4).SetProperty(properties.Omit, true)
	show(t)
	t.Column(2).SetProperty(properties.Omit, true)
	show(t)
	t.Column(1).SetProperty(properties.Omit, true)
	if err := t.RenderTo(os.Stdout); err == nil || err != tabular.ErrNoColumnsToDisplay {
		if err == nil {
			fmt.Fprintf(os.Stderr, "table did not fail to render with no columns to display\n")
		} else {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	// Delete all the properties, then set a default column property and override for the columns we _should_ display
	for i := range 4 {
		t.Column(i).SetProperty(properties.Omit, nil)
	}
	t.Column(0).SetProperty(properties.Omit, true)
	for _, cname := range []string{"Person", "Color"} {
		c, err := t.ColumnNamed(cname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "table lookup of column %q failed: %v\n", cname, err)
			continue
		}
		c.SetProperty(properties.Omit, false)
	}
	show(t)

	// Output:
	// ╭─────────┬─────┬───────┬────────╮
	// │ Person  │ Age │ Score │ Color  │
	// ├─────────┼─────┼───────┼────────┤
	// │ Fred    │  34 │ 0.6   │ red    │
	// │ Gladys  │  32 │ 0.8   │ yellow │
	// │ Bert    │  57 │ 0.7   │ green  │
	// │ Belinda │  58 │ 0.8   │ blue   │
	// ╰─────────┴─────┴───────┴────────╯
	// ╭─────────┬─────┬────────╮
	// │ Person  │ Age │ Color  │
	// ├─────────┼─────┼────────┤
	// │ Fred    │  34 │ red    │
	// │ Gladys  │  32 │ yellow │
	// │ Bert    │  57 │ green  │
	// │ Belinda │  58 │ blue   │
	// ╰─────────┴─────┴────────╯
	// ╭─────────┬─────┬───────┬────────╮
	// │ Person  │ Age │ Score │ Color  │
	// ├─────────┼─────┼───────┼────────┤
	// │ Fred    │  34 │ 0.6   │ red    │
	// │ Gladys  │  32 │ 0.8   │ yellow │
	// │ Bert    │  57 │ 0.7   │ green  │
	// │ Belinda │  58 │ 0.8   │ blue   │
	// ╰─────────┴─────┴───────┴────────╯
	// ╭─────────┬─────┬───────┬────────╮
	// │ Person  │ Age │ Score │ Color  │
	// ├─────────┼─────┼───────┼────────┤
	// │ Fred    │  34 │ 0.6   │ red    │
	// │ Gladys  │  32 │ 0.8   │ yellow │
	// │ Bert    │  57 │ 0.7   │ green  │
	// │ Belinda │  58 │ 0.8   │ blue   │
	// ╰─────────┴─────┴───────┴────────╯
	// ╭─────────┬─────╮
	// │ Person  │ Age │
	// ├─────────┼─────┤
	// │ Fred    │  34 │
	// │ Gladys  │  32 │
	// │ Bert    │  57 │
	// │ Belinda │  58 │
	// ╰─────────┴─────╯
	// ╭─────────╮
	// │ Person  │
	// ├─────────┤
	// │ Fred    │
	// │ Gladys  │
	// │ Bert    │
	// │ Belinda │
	// ╰─────────╯
	// ╭─────────┬────────╮
	// │ Person  │ Color  │
	// ├─────────┼────────┤
	// │ Fred    │ red    │
	// │ Gladys  │ yellow │
	// │ Bert    │ green  │
	// │ Belinda │ blue   │
	// ╰─────────┴────────╯
}

func Example_omit_json() {
	t := tabular.New()
	populateFourTable(t)

	show := func(tb tabular.Table) {
		if err := tab_json.Wrap(tb).RenderTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	show(t)
	t.Column(3).SetProperty(properties.Omit, true)
	fmt.Println("---")
	show(t)

	// Output:
	// [
	// {"Person": "Fred", "Age": 34, "Score": 0.6, "Color": "red"},
	// {"Person": "Gladys", "Age": 32, "Score": 0.8, "Color": "yellow"},
	// {"Person": "Bert", "Age": 57, "Score": 0.7, "Color": "green"},
	// {"Person": "Belinda", "Age": 58, "Score": 0.8, "Color": "blue"}
	// ]
	// ---
	// [
	// {"Person": "Fred", "Age": 34, "Color": "red"},
	// {"Person": "Gladys", "Age": 32, "Color": "yellow"},
	// {"Person": "Bert", "Age": 57, "Color": "green"},
	// {"Person": "Belinda", "Age": 58, "Color": "blue"}
	// ]
}

func Example_omit_csv() {
	t := tabular.New()
	populateFourTable(t)

	show := func(tb tabular.Table) {
		if err := tab_csv.Wrap(tb).RenderTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	show(t)
	t.Column(3).SetProperty(properties.Omit, true)
	fmt.Println("---")
	show(t)

	// Output:
	// "Person","Age","Score","Color"
	// "Fred","34","0.6","red"
	// "Gladys","32","0.8","yellow"
	// "Bert","57","0.7","green"
	// "Belinda","58","0.8","blue"
	// ---
	// "Person","Age","Color"
	// "Fred","34","red"
	// "Gladys","32","yellow"
	// "Bert","57","green"
	// "Belinda","58","blue"
}

func Example_omit_markdown() {
	t := tabular.New()
	populateFourTable(t)

	show := func(tb tabular.Table) {
		if err := tab_md.Wrap(tb).RenderTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	show(t)
	t.Column(3).SetProperty(properties.Omit, true)
	fmt.Println("---")
	show(t)

	// Output:
	// | Person  | Age | Score | Color  |
	// | ------- | ---:| ----- | ------ |
	// | Fred    |  34 | 0.6   | red    |
	// | Gladys  |  32 | 0.8   | yellow |
	// | Bert    |  57 | 0.7   | green  |
	// | Belinda |  58 | 0.8   | blue   |
	// ---
	// | Person  | Age | Color  |
	// | ------- | ---:| ------ |
	// | Fred    |  34 | red    |
	// | Gladys  |  32 | yellow |
	// | Bert    |  57 | green  |
	// | Belinda |  58 | blue   |
}

func Example_omit_html() {
	t := tabular.New()
	populateFourTable(t)

	show := func(tb tabular.Table) {
		if err := tab_html.Wrap(tb).RenderTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
		}
	}

	show(t)
	t.Column(3).SetProperty(properties.Omit, true)
	fmt.Println("---")
	show(t)

	// Output:
	// <table>
	//   <thead>
	//     <tr><th>Person</th><th>Age</th><th>Score</th><th>Color</th></tr>
	//   </thead>
	//   <tbody>
	//     <tr><td>Fred</td><td>34</td><td>0.6</td><td>red</td></tr>
	//     <tr><td>Gladys</td><td>32</td><td>0.8</td><td>yellow</td></tr>
	//     <tr><td>Bert</td><td>57</td><td>0.7</td><td>green</td></tr>
	//     <tr><td>Belinda</td><td>58</td><td>0.8</td><td>blue</td></tr>
	//   </tbody>
	// </table>
	// ---
	// <table>
	//   <thead>
	//     <tr><th>Person</th><th>Age</th><th>Color</th></tr>
	//   </thead>
	//   <tbody>
	//     <tr><td>Fred</td><td>34</td><td>red</td></tr>
	//     <tr><td>Gladys</td><td>32</td><td>yellow</td></tr>
	//     <tr><td>Bert</td><td>57</td><td>green</td></tr>
	//     <tr><td>Belinda</td><td>58</td><td>blue</td></tr>
	//   </tbody>
	// </table>

}
