// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package markdown_test // import "go.pennock.tech/tabular/markdown"

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"unicode"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/markdown"
)

func testViaCreatorFunc(t *testing.T, creator func() tabular.Table) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := creator()
	T.NotEqual(tb, nil, "have a table")

	have, err := markdown.Render(tb)
	T.Equal(tb.Errors(), nil, "no errors stored when rendering empty table via pkg func")
	T.ExpectError(err, "should have failed to render the table while empty")

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	should := `
| foo    | loquacious | x    |
| ------ | ---------- | ---- |
| 42     | .          | fred |
| snerty | word       | r    |
|        | true       |      |
`
	should = strings.TrimLeftFunc(should, unicode.IsSpace)

	// When creator is creating a markdown.MarkdownTable, this will double-wrap, but
	// should still work (same as we should be able to wrap tabular.Table
	// inside csv.CSVTable inside markdown.MarkdownTable)
	have, err = markdown.Render(tb)
	T.ExpectSuccess(err, "no errors returned when rendering basic table via pkg func")
	T.Equal(tb.Errors(), nil, "no errors stored when rendering basic table via pkg func")
	T.Equal(have, should, "basic output emits cleanly via pkg func")

	tb2 := creator()
	T.NotEqual(tb2, nil, "have a table")
	tb2.AddRowItems("42", "fred")
	_, err = markdown.Render(tb2)
	T.ExpectError(err, "headerless table should have failed to render")

	tb.AddRowItems("a\nb\nc", `d"e`, 1)
	should += strings.TrimLeftFunc(`
| a&#x0a;b&#x0a;c | d&#34;e    | 1    |
`, unicode.IsSpace)

	have, err = markdown.Render(tb)
	T.ExpectSuccess(err, "no errors returned when rendering multiline quote-containing table")
	T.Equal(tb.Errors(), nil, "no errors stored when rendering multiline quote-containing table")
	T.Equal(have, should, "multiline quote-containing output emits 'cleanly'")
}

func TestTableMarkdownOfTabular(t *testing.T) {
	testViaCreatorFunc(t, func() tabular.Table { return tabular.New() })
}

func TestTableMarkdownDirectly(t *testing.T) {
	testViaCreatorFunc(t, func() tabular.Table { return markdown.New() })
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
	broken := markdown.Wrap(&brokeTB)

	_, err := broken.Render()
	T.ExpectError(err, "table should fail to render if too few columns for rows")

	tb.AddHeaders("1", "2", "3", "4")
	brokeTB.overrideColumns = 4
	_, err = broken.Render()
	T.ExpectSuccess(err, "table renders something when enough columns")
	brokeTB.overrideColumns = 3
	_, err = broken.Render()
	T.ExpectError(err, "table should fail to render if too few columns for headers")
}

func TestFileWritingMarkdown(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := tabular.New()
	T.NotEqual(tb, nil, "have a table")

	should := []byte(`
| foo    | loquacious | x    |
| ------ | ---------- | ---- |
| 42     | .          | fred |
| snerty | word       | r    |
|        | true       |      |
`)
	should = bytes.TrimLeftFunc(should, unicode.IsSpace)

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	temp := T.TempFile()
	err := markdown.RenderTo(tb, temp)
	T.ExpectSuccess(err, "table renders to temporary file")

	newOff, err := temp.Seek(0, 0)
	T.ExpectSuccess(err, "seek-to-start of temp file succeeds")
	T.Equal(newOff, int64(0), "offset after seeking to start is 0")

	contents, err := ioutil.ReadAll(temp)
	T.ExpectSuccess(err, "reading contents from temp file")
	T.Equal(contents, should, "content in filesystem is as expected")
}

// degenerate case, no affixed separators
func TestSingleColumnMarkdown(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := markdown.New()
	T.NotEqual(tb, nil, "have a table")
	should := `
| header |
| ------ |
| line 1 |
| two    |
`
	should = strings.TrimLeftFunc(should, unicode.IsSpace)

	tb.AddHeaders("header")
	tb.AddRowItems("line 1")
	tb.AddRowItems("two")
	T.Equal(tb.Errors(), nil, "no errors just adding items")
	have, err := tb.Render()
	T.ExpectSuccess(err, "single-column table renders without errors")
	T.Equal(have, should, "got correct single-column output")
}
