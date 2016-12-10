// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable_test // import "go.pennock.tech/tabular/texttable"

import (
	"io/ioutil"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular/texttable"

	// for getting the CellLocation type
	"go.pennock.tech/tabular"
	// for testing via the named constants & functions
	"go.pennock.tech/tabular/texttable/decoration"
)

type str struct {
	S string
}

func (s str) String() string { return s.S }

type goStr struct {
	S string
}

func (s goStr) GoString() string { return s.S }

type errStr struct {
	S string
}

func (s errStr) Error() string { return s.S }

func TestTableCreation(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := texttable.New()
	T.NotEqual(tb, nil, "have a table")

	const emptyTableRenderBoxHeavy = "" +
		"┏┓\n" +
		"┗┛\n" +
		""

	rendered, err := tb.Render()
	T.ExpectSuccess(err, "empty table rendered")
	T.Equal(rendered, emptyTableRenderBoxHeavy, "empty table rendered correctly")
	T.Equal(tb.Errors(), nil, "no errors in empty table")

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	tb.AddRowItems(str{"alpha"}, goStr{"beta"}, errStr{"gamma"})
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	T.Equal(tb.NColumns(), 3, "table should have 3 columns")
	T.Equal(tb.NRows(), 5, "table should have 5 rows in the body")

	for _, status := range []struct {
		good bool
		r, c int
		want string
	}{
		{true, 1, 1, "42"},
		{true, 1, 3, "fred"},
		{false, 0, 3, ""},
		{false, 3, 0, ""},
		{false, -1, 3, ""},
		{false, 1, 4, ""},
		{false, 3, 1, ""},
		{true, 4, 1, " "},
		{true, 5, 1, "alpha"},
		{true, 5, 2, "beta"},
		{true, 5, 3, "gamma"},
		{false, 6, 1, ""},
	} {
		c, err := tb.CellAt(tabular.CellLocation{Row: status.r, Column: status.c})
		if status.good {
			T.ExpectSuccessf(err, "cell [%d,%d] should have been available", status.r, status.c)
			if c != nil {
				T.Equalf(c.String(), status.want,
					"cell [%d,%d] contents did not render to expected %q",
					status.r, status.c, status.want)
			}
		} else {
			T.ExpectErrorf(err, "cell [%d,%d] should not have been available", status.r, status.c)
		}
	}
}

func createStdTableContents(T *testlib.T) *texttable.TextTable {
	tb := texttable.New()
	T.NotEqual(tb, nil, "have a table")

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	return tb
}

func TestTableRenderingDefault(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	// Happens to be the heavy box style, is this subject to change?
	should := "" +
		"┏━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━┓\n" +
		"┃ foo    ┃ loquacious ┃ x    ┃\n" +
		"┣━━━━━━━━╇━━━━━━━━━━━━╇━━━━━━┫\n" +
		"┃ 42     │ .          │ fred ┃\n" +
		"┃ snerty │ word       │ r    ┃\n" +
		"┠────────┼────────────┼──────┨\n" +
		"┃        │ true       │      ┃\n" +
		"┗━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━┛\n" +
		""
	rendered, err := tb.Render()
	T.ExpectSuccess(err, "simple table rendered (default style)")
	T.Equal(tb.Errors(), nil, "no errors rendering table (default style)")
	T.Equal(rendered, should, "simple table rendered correctly (default style)")

	// confirm that the table isn't damaged by rendering and emits fine _repeatedly_
	for i := 2; i <= 5; i++ {
		r2, err := tb.Render()
		T.ExpectSuccessf(err, "simple table rendered [pass %d]", i)
		T.Equalf(tb.Errors(), nil, "no errors rendering table [pass %d]", i)
		T.Equalf(r2, should, "simple table rendered correctly [pass %d]", i)
	}

	rawTb := tb.Table
	temp := T.TempFile()
	texttable.RenderTo(rawTb, temp)
	temp.Close()
	tempContents, err := ioutil.ReadFile(temp.Name())
	T.ExpectSuccess(err, "unable to re-open tempfile")
	T.Equal(tempContents, []byte(should), "simple table wrote to file-system fine (via pkg function)")
	T.Equal(tb.Errors(), nil, "no errors rendering table to file (via pkg function)")
}

func TestTableRenderingAlign(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	tabular.SetAlignmentStatic(tb.Column(1), tabular.ALIGN_RIGHT)

	should := "" +
		"┏━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━┓\n" +
		"┃    foo ┃ loquacious ┃ x    ┃\n" +
		"┣━━━━━━━━╇━━━━━━━━━━━━╇━━━━━━┫\n" +
		"┃     42 │ .          │ fred ┃\n" +
		"┃ snerty │ word       │ r    ┃\n" +
		"┠────────┼────────────┼──────┨\n" +
		"┃        │ true       │      ┃\n" +
		"┗━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━┛\n" +
		""

	rendered, err := tb.Render()
	T.ExpectSuccess(err, "table rendered, align: col1-right")
	T.Equal(tb.Errors(), nil, "no errors rendering table (align: col1-right)")
	T.Equal(rendered, should, "table rendered correctly (align: col1-right)")
}

func TestTableRenderingLightByConstNamed(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	// UTF8BoxLight
	should := "" +
		"┌────────┬────────────┬──────┐\n" +
		"│ foo    │ loquacious │ x    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│ 42     │ .          │ fred │\n" +
		"│ snerty │ word       │ r    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│        │ true       │      │\n" +
		"└────────┴────────────┴──────┘\n" +
		""
	_, err := tb.SetDecorationNamed(decoration.D_UTF8_LIGHT)
	T.ExpectSuccess(err, "got a decoration with a const imported name")
	rendered, err := tb.Render()
	T.ExpectSuccess(err, "simple table rendered (light style)")
	T.Equal(tb.Errors(), nil, "no errors rendering table (light style)")
	T.Equal(rendered, should, "simple table rendered correctly (light style)")
}

func TestTableRenderingLightCurvedByName(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	// we unpick this to have a test path covering the pkg functions
	rawTb := createStdTableContents(T).Table

	// UTF8BoxLightCurved
	should := "" +
		"╭────────┬────────────┬──────╮\n" +
		"│ foo    │ loquacious │ x    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│ 42     │ .          │ fred │\n" +
		"│ snerty │ word       │ r    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│        │ true       │      │\n" +
		"╰────────┴────────────┴──────╯\n" +
		""
	tb, err := texttable.Wrap(rawTb).SetDecorationNamed("utf8-light-curved")
	T.ExpectSuccess(err, "got a decoration with a standard name")
	rendered, err := tb.Render()
	T.ExpectSuccess(err, "simple table rendered (light-curved style)")
	T.Equal(tb.Errors(), nil, "no errors rendering table (light-curved style)")
	T.Equal(rendered, should, "simple table rendered correctly (light-curved style)")

	temp := T.TempFile()
	tb.RenderTo(temp)
	temp.Close()
	tempContents, err := ioutil.ReadFile(temp.Name())
	T.ExpectSuccess(err, "unable to re-open tempfile")
	T.Equal(tempContents, []byte(should), "simple table wrote to file-system fine (via method)")
	T.Equal(tb.Errors(), nil, "no errors rendering table to file (via method)")

	// pkg function is not in this test function, but in the function which
	// works with the default decoration.  The reasons should be blindingly
	// obvious, even if only in retrospect after head-scratching and then
	// feeling very stupid.
}

func TestTableRenderingDoubleByDecor(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	// UTF8BoxDouble
	should := "" +
		"╔════════╦════════════╦══════╗\n" +
		"║ foo    ║ loquacious ║ x    ║\n" +
		"╠════════╪════════════╪══════╣\n" +
		"║ 42     │ .          │ fred ║\n" +
		"║ snerty │ word       │ r    ║\n" +
		"╟────────┼────────────┼──────╢\n" +
		"║        │ true       │      ║\n" +
		"╚════════╧════════════╧══════╝\n" +
		""
	tb2 := tb.SetDecoration(decoration.UTF8BoxDouble())
	T.Equal(tb, tb2, "SetDecoration should have mutated and returned self-same object")
	rendered, err := tb.Render()
	T.ExpectSuccess(err, "simple table rendered (double style, func-called style set)")
	T.Equal(tb.Errors(), nil, "no errors rendering table (double style)")
	T.Equal(rendered, should, "simple table rendered correctly (double style)")
}

func TestTableUnknownDecoration(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	tb2, err := tb.SetDecorationNamed("unknown, do not make me exist, pretty please")
	T.ExpectError(err, "setting decoration to a bad one should have succeeded")
	T.NotEqual(tb2, nil, "failed decoration setting should still have chained back the core object")
	T.Equal(tb2, tb, "core object and returned object same when set bad decoration")
	rendered, err := tb.Render()
	T.ExpectError(err, "should have failed to render with unknown decoration")
	T.Equal(rendered, "", "should have gotten empty render contents")
}

func TestTableBoxless(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()
	tb := createStdTableContents(T)

	// NoBox
	should := "" +
		"foo    loquacious x   \n" +
		"42     .          fred\n" +
		"snerty word       r   \n" +
		"       true           \n" +
		""
	_, err := tb.SetDecorationNamed(decoration.D_NONE)
	T.ExpectSuccess(err, "set a const decoration")
	rendered, err := tb.Render()
	T.ExpectSuccess(err, "simple table rendered (boxless)")
	T.Equal(tb.Errors(), nil, "no errors rendering table (boxless)")
	T.Equal(rendered, should, "simple table rendered correctly (boxless)")
}
