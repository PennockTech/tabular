// Copyright © 2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular_test // import "go.pennock.tech/tabular"

import (
	"strings"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/texttable"
)

type greek struct {
	s string
	n int64
}

func (g greek) String() string   { return g.s }
func (g greek) SortInt64() int64 { return g.n }

var _ tabular.SortInter = greek{}

// alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu nu xi omicron pi rho sigma tau upsilon phi chi psi omega
// 1     2    3     4     5       6    7   8     9    10    11     12 13 14 15      16 17  18    19  20      21  22  23  24

func populate(T *testlib.T, tb tabular.Table) {
	tb.AddHeaders("Label", "Greek", "Price", "Quantity", "Rating")
	tb.AddRowItems("Mu", greek{"mu", 12}, 0.02, 500, 4.5)
	tb.AddRowItems("Zeta", greek{"zeta", 6}, 3, 1000, 3.0)
	tb.AddRowItems("Gamma", greek{"gamma", 3}, 400, 1, 4.2)
	tb.AddRowItems("Alpha", greek{"alpha", 1}, 50, 10, 4.0)
	tb.AddRowItems("Chi", greek{"chi", 22}, 42.34, 7, 0.0)
	T.Equal(tb.Errors(), nil, "no errors just adding items")
}

const insertionOrder = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Gamma | gamma | 400   | 1        | 4.2    |
| Alpha | alpha | 50    | 10       | 4      |
| Chi   | chi   | 42.34 | 7        | 0      |
+-------+-------+-------+----------+--------+
`

const sortedByLabelAsc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Alpha | alpha | 50    | 10       | 4      |
| Chi   | chi   | 42.34 | 7        | 0      |
| Gamma | gamma | 400   | 1        | 4.2    |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Zeta  | zeta  | 3     | 1000     | 3      |
+-------+-------+-------+----------+--------+
`

const sortedByLabelDesc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Zeta  | zeta  | 3     | 1000     | 3      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Gamma | gamma | 400   | 1        | 4.2    |
| Chi   | chi   | 42.34 | 7        | 0      |
| Alpha | alpha | 50    | 10       | 4      |
+-------+-------+-------+----------+--------+
`

const sortedByGreekAsc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Alpha | alpha | 50    | 10       | 4      |
| Gamma | gamma | 400   | 1        | 4.2    |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Chi   | chi   | 42.34 | 7        | 0      |
+-------+-------+-------+----------+--------+
`

const sortedByGreekDesc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Chi   | chi   | 42.34 | 7        | 0      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Gamma | gamma | 400   | 1        | 4.2    |
| Alpha | alpha | 50    | 10       | 4      |
+-------+-------+-------+----------+--------+
`

const sortedByQuantityAsc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Gamma | gamma | 400   | 1        | 4.2    |
| Chi   | chi   | 42.34 | 7        | 0      |
| Alpha | alpha | 50    | 10       | 4      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Zeta  | zeta  | 3     | 1000     | 3      |
+-------+-------+-------+----------+--------+
`

const sortedByQuantityDesc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Zeta  | zeta  | 3     | 1000     | 3      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Alpha | alpha | 50    | 10       | 4      |
| Chi   | chi   | 42.34 | 7        | 0      |
| Gamma | gamma | 400   | 1        | 4.2    |
+-------+-------+-------+----------+--------+
`

const sortedByPriceAsc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Chi   | chi   | 42.34 | 7        | 0      |
| Alpha | alpha | 50    | 10       | 4      |
| Gamma | gamma | 400   | 1        | 4.2    |
+-------+-------+-------+----------+--------+
`

const sortedByPriceDesc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Gamma | gamma | 400   | 1        | 4.2    |
| Alpha | alpha | 50    | 10       | 4      |
| Chi   | chi   | 42.34 | 7        | 0      |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Mu    | mu    | 0.02  | 500      | 4.5    |
+-------+-------+-------+----------+--------+
`

const sortedByRatingAsc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Chi   | chi   | 42.34 | 7        | 0      |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Alpha | alpha | 50    | 10       | 4      |
| Gamma | gamma | 400   | 1        | 4.2    |
| Mu    | mu    | 0.02  | 500      | 4.5    |
+-------+-------+-------+----------+--------+
`

const sortedByRatingDesc = `
+-------+-------+-------+----------+--------+
| Label | Greek | Price | Quantity | Rating |
+-------+-------+-------+----------+--------+
| Mu    | mu    | 0.02  | 500      | 4.5    |
| Gamma | gamma | 400   | 1        | 4.2    |
| Alpha | alpha | 50    | 10       | 4      |
| Zeta  | zeta  | 3     | 1000     | 3      |
| Chi   | chi   | 42.34 | 7        | 0      |
+-------+-------+-------+----------+--------+
`

func TestSortByColumn(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	var (
		rendered string
		err      error
	)

	tb := texttable.New()
	tb.SetDecorationNamed("ascii-simple")
	populate(T, tb)
	rendered, err = tb.Render()
	T.ExpectSuccess(err, "render before sorting")
	T.Equal(rendered, strings.TrimLeft(insertionOrder, "\n"), "unsorted matches what we inserted")
	err = tb.SortByNamedColumn("missing", tabular.SORT_ASC)
	T.ExpectError(err, "should have failed to sort by a missing column")

	checks := []struct {
		Column string
		Order  tabular.SortOrder
		Want   string
	}{
		{"Label", tabular.SORT_ASC, sortedByLabelAsc},
		{"Label", tabular.SORT_DESC, sortedByLabelDesc},
		{"Greek", tabular.SORT_ASC, sortedByGreekAsc},
		{"Greek", tabular.SORT_DESC, sortedByGreekDesc},
		{"Quantity", tabular.SORT_ASC, sortedByQuantityAsc},
		{"Quantity", tabular.SORT_DESC, sortedByQuantityDesc},
		{"Price", tabular.SORT_ASC, sortedByPriceAsc},
		{"Price", tabular.SORT_DESC, sortedByPriceDesc},
		{"Rating", tabular.SORT_ASC, sortedByRatingAsc},
		{"Rating", tabular.SORT_DESC, sortedByRatingDesc},
	}

	for i := range checks {
		err = tb.SortByNamedColumn(checks[i].Column, checks[i].Order)
		T.ExpectSuccess(err, "able to sort by "+checks[i].Column+" "+checks[i].Order.String())

		/*
			fmt.Printf("\nSORT BY %s %s\n", checks[i].Column, checks[i].Order)
			var col int
			switch checks[i].Column {
			case "Label":
				col = 1
			case "Greek":
				col = 2
			case "Price":
				col = 3
			case "Quantity":
				col = 4
			case "Rating":
				col = 5
			}
			for r := 0; r < tb.NRows(); r += 1 {
				c, err := tb.CellAt(tabular.CellLocation{Row: r + 1, Column: col})
				if err != nil {
					fmt.Printf("row %d cell failed: %v\n", r, err)
					continue
				}
				fmt.Printf("row %d cell: %s | %T | %v\n", r, c, c.Item(), c.Item())
			}
		*/

		rendered, err = tb.Render()
		T.ExpectSuccess(err, "rendered sorted by "+checks[i].Column+" "+checks[i].Order.String())
		T.Equal(rendered, strings.TrimLeft(checks[i].Want, "\n"), "sorted by "+checks[i].Column+" "+checks[i].Order.String())
		T.Log("match on sorting by " + checks[i].Column + " " + checks[i].Order.String())
	}
}
