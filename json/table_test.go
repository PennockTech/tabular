// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package json_test // import "go.pennock.tech/tabular/json"

import (
	"io/ioutil"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
	tab_json "go.pennock.tech/tabular/json"
)

func testViaCreatorFunc(t *testing.T, creator func() tabular.Table) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := creator()
	T.NotEqual(tb, nil, "have a table")

	have, err := tab_json.Render(tb)
	T.Equal(tb.Errors(), nil, "no errors stored when rendering empty table via pkg func")
	T.ExpectError(err, "should have failed to render the table while empty")

	tb.AddHeaders("foo", "loquacious", "x")
	have, err = tab_json.Render(tb)
	T.Equal(tb.Errors(), nil, "no errors stored when rendering headers-only table via pkg func")
	T.Equal(have, "[\n\n]\n", "empty table renders to empty-but-present array")

	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	should := `[
{"foo": 42, "loquacious": ".", "x": "fred"},
{"foo": "snerty", "loquacious": "word", "x": "r"},

{"foo": " ", "loquacious": true, "x": null}
]
`

	have, err = tab_json.Render(tb)
	T.ExpectSuccess(err, "no errors returned when rendering basic table via pkg func")
	T.Equal(tb.Errors(), nil, "no errors stored when rendering basic table via pkg func")
	T.Equal(have, should, "basic output emits cleanly via pkg func")

	tb = creator()
	T.NotEqual(tb, nil, "have a table")
	tb.AddRowItems("42", "fred")
	have, err = tab_json.Render(tb)
	T.ExpectErrorf(err, "shoulf have failed to render headerless table. instead got: %v", have)
}

func TestTableJSONOfTabular(t *testing.T) {
	testViaCreatorFunc(t, func() tabular.Table { return tabular.New() })
}

func TestTableJSONDirectly(t *testing.T) {
	testViaCreatorFunc(t, func() tabular.Table { return tab_json.New() })
}

type brokenTable struct {
	*tabular.ATable
	overrideColumns int
}

func (b brokenTable) NColumns() int { return b.overrideColumns }

func TestBrokenTables(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := tabular.New()
	T.NotEqual(tb, nil, "have a table")

	tb.AddRowItems("alpha", "beta", "gamma")
	brokeTB := brokenTable{ATable: tb}
	brokeTB.overrideColumns = 2
	broken := tab_json.Wrap(&brokeTB)

	have, err := broken.Render()
	T.ExpectErrorf(err, "table should fail to render if too few columns for rows, got: %v", have)

	tb.AddHeaders("1", "2", "3", "4")
	brokeTB.overrideColumns = 4
	_, err = broken.Render()
	T.ExpectSuccess(err, "table renders something when enough columns")
	brokeTB.overrideColumns = 3
	have, err = broken.Render()
	T.ExpectSuccess(err, "table renders something when have header surplus")
	brokeTB.overrideColumns = 3
	tb.AddHeaders("x", "y")
	have, err = broken.Render()
	T.ExpectErrorf(err, "table should fail to render if too few headers for columns, got: %v", have)
}

func TestFileWritingJSON(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := tabular.New()
	T.NotEqual(tb, nil, "have a table")

	should := []byte(`[
{"foo": 42, "loquacious": ".", "x": "fred"},
{"foo": "snerty", "loquacious": "word", "x": "r"},

{"foo": " ", "loquacious": true, "x": null}
]
`)

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	temp := T.TempFile()
	err := tab_json.RenderTo(tb, temp)
	T.ExpectSuccess(err, "table renders to temporary file")

	newOff, err := temp.Seek(0, 0)
	T.ExpectSuccess(err, "seek-to-start of temp file succeeds")
	T.Equal(newOff, int64(0), "offset after seeking to start is 0")

	contents, err := ioutil.ReadAll(temp)
	T.ExpectSuccess(err, "reading contents from temp file")
	T.Equal(string(contents), string(should), "content in filesystem is as expected")
}

// degenerate case, no internal commas in rows
func TestSingleColumnJSON(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := tab_json.New()
	T.NotEqual(tb, nil, "have a table")
	should := `[
{"header": "line 1"},
{"header": "two"}
]
`
	tb.AddHeaders("header")
	tb.AddRowItems("line 1")
	tb.AddRowItems("two")
	T.Equal(tb.Errors(), nil, "no errors just adding items")
	have, err := tb.Render()
	T.ExpectSuccess(err, "single-column table renders without errors")
	T.Equal(have, should, "got correct single-column output")
}