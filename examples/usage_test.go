package main

import (
	"fmt"
	"os"

	auto_table "go.pennock.tech/tabular/auto"
)

func Example_simple() {
	t := auto_table.New("utf8-heavy")
	t.AddHeaders("Fruit", "Price $/lb")
	t.AddRowItems("bananas", "0.56")
	t.AddRowItems("apples", "1.31")
	t.AddRowItems("strawberries", "2.22")
	t.AddRowItems("pears", "1.90")

	if errs := t.Errors(); errs != nil {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "table error: %s\n", err)
		}
	}

	// This would get you the table as a string, instead of using an io.Writer:
	// tString, err := t.Render()

	err := t.RenderTo(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
	}
	// Output:
	// ┏━━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
	// ┃ Fruit        ┃ Price $/lb ┃
	// ┣━━━━━━━━━━━━━━╇━━━━━━━━━━━━┫
	// ┃ bananas      │ 0.56       ┃
	// ┃ apples       │ 1.31       ┃
	// ┃ strawberries │ 2.22       ┃
	// ┃ pears        │ 1.90       ┃
	// ┗━━━━━━━━━━━━━━┷━━━━━━━━━━━━┛
}

func populateTable(t auto_table.RenderTable) auto_table.RenderTable {
	t.AddHeaders("Fruit", "Price $/lb")
	t.AddRowItems("bananas", "0.56")
	t.AddRowItems("apples", "1.31")
	t.AddRowItems("strawberries", "2.22")
	t.AddRowItems("pears", "1.90")

	// Instead of checking for nil, as in Example_simple, which lets you decide
	// whether or not to error out etc, if you're just going to iterate then
	// you can just rely upon range over nil never invoking the body.
	for _, err := range t.Errors() {
		fmt.Fprintf(os.Stderr, "table error: %s\n", err)
	}

	return t
}

func Example_ascii() {
	t := auto_table.New("ascii-simple")
	t = populateTable(t)

	err := t.RenderTo(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "table rendering error: %s\n", err)
	}
	fmt.Printf("Table has %d rows and %d columns\n", t.NRows(), t.NColumns())

	// Output:
	// +--------------+------------+
	// | Fruit        | Price $/lb |
	// +--------------+------------+
	// | bananas      | 0.56       |
	// | apples       | 1.31       |
	// | strawberries | 2.22       |
	// | pears        | 1.90       |
	// +--------------+------------+
	// Table has 4 rows and 2 columns
}

func Example_html() {
	// We can use `"html"` as a style when creating the table too, but the
	// various styles let you embed other tables, and the auto package has
	// package-level functions to take any table and render in any style.
	t := auto_table.New("utf8-light-curved")
	t = populateTable(t)

	_ = t.RenderTo(os.Stdout)
	_ = auto_table.RenderTo(t, os.Stdout, "html")

	// Output:
	// ╭──────────────┬────────────╮
	// │ Fruit        │ Price $/lb │
	// ├──────────────┼────────────┤
	// │ bananas      │ 0.56       │
	// │ apples       │ 1.31       │
	// │ strawberries │ 2.22       │
	// │ pears        │ 1.90       │
	// ╰──────────────┴────────────╯
	// <table>
	//   <thead>
	//     <tr><th>Fruit</th><th>Price $/lb</th></tr>
	//   </thead>
	//   <tbody>
	//     <tr><td>bananas</td><td>0.56</td></tr>
	//     <tr><td>apples</td><td>1.31</td></tr>
	//     <tr><td>strawberries</td><td>2.22</td></tr>
	//     <tr><td>pears</td><td>1.90</td></tr>
	//   </tbody>
	// </table>
}
