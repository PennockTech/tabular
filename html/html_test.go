// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package html_test // import "go.pennock.tech/tabular/html"

import (
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/html"
)

func TestHTMLTableRendering(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := tabular.New()
	T.NotEqual(tb, nil, "have a table")

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	const shouldTail = `  <thead>
    <tr><th>foo</th><th>loquacious</th><th>x</th></tr>
  </thead>
  <tbody>
    <tr><td>42</td><td>.</td><td>fred</td></tr>
    <tr><td>snerty</td><td>word</td><td>r</td></tr>
    <tr><td> </td><td>true</td><td></td></tr>
  </tbody>
</table>
`

	should := "<table>\n" + shouldTail
	rendered, err := html.Wrap(tb).Render()
	T.ExpectSuccess(err, "simple table rendered to HTML")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "simple table rendered to HTML correctly")

	ht := html.Wrap(tb)
	ht.Id = "foo"
	should = "<table id=\"foo\">\n" + shouldTail
	rendered, err = ht.Render()
	T.ExpectSuccess(err, "table rendered to HTML with id")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "table rendered to HTML with id correctly")

	ht.Class = "bar"
	should = "<table class=\"bar\" id=\"foo\">\n" + shouldTail
	rendered, err = ht.Render()
	T.ExpectSuccess(err, "table rendered to HTML with id and class")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "table rendered to HTML with id and class correctly")

	ht.Id = ""
	should = "<table class=\"bar\">\n" + shouldTail
	rendered, err = ht.Render()
	T.ExpectSuccess(err, "table rendered to HTML with class")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "table rendered to HTML with class correctly")

	ht.Caption = "A test table"
	should = "<table class=\"bar\">\n" + "  <caption>A test table</caption>\n" + shouldTail
	rendered, err = ht.Render()
	T.ExpectSuccess(err, "table rendered to HTML with class and caption")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "table rendered to HTML with class and caption correctly")

	// Warning: table contents mutated here, tests below should hand;e accordingly

	tb.AddHeaders("foo", "less verbose", "x")
	should = strings.Replace(should, "loquacious", "less verbose", -1)
	rendered, err = ht.Render()
	T.ExpectSuccess(err, "table rendered, after changing headers")
	T.Equal(tb.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "table rendered fine after changing headers")
}

type flipflop struct {
	isBlue bool
}

func blueGreenFlipflop(rowNum int, ffint interface{}) template.HTMLAttr {
	ff := ffint.(*flipflop)
	ff.isBlue = !ff.isBlue
	if ff.isBlue {
		return template.HTMLAttr(fmt.Sprintf("blue r%d", rowNum))
	} else {
		return template.HTMLAttr(fmt.Sprintf("green r%d", rowNum))
	}
}

func TestHTMLRowClassGenerator(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := html.New()
	T.NotEqual(tb, nil, "have a table")

	tb2 := html.New()
	T.NotEqual(tb2, nil, "have a second table")
	T.Equal(tb, tb2, "empty tables equal")
	tb2.SetRowClassGenerator(blueGreenFlipflop, &flipflop{})
	T.NotEqual(tb, tb2, "setting row class generator mutated tb2")

	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")

	tb.SetRowClassGenerator(blueGreenFlipflop, &flipflop{})

	const should = `<table>
  <thead>
    <tr class="blue r0"><th>foo</th><th>loquacious</th><th>x</th></tr>
  </thead>
  <tbody>
    <tr class="green r1"><td>42</td><td>.</td><td>fred</td></tr>
    <tr class="blue r2"><td>snerty</td><td>word</td><td>r</td></tr>
    <tr class="green r4"><td> </td><td>true</td><td></td></tr>
  </tbody>
</table>
`

	rendered, err := tb.Render()
	T.ExpectSuccess(err, "rcg table rendered to HTML")
	T.Equal(tb.Errors(), nil, "no errors accumulated in rcg table through rendering")
	T.Equal(rendered, should, "rcg table rendered to HTML correctly")
}

func TestHTMLTableIsTable(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	ht := html.New()
	T.NotEqual(ht, nil, "have a table")

	ht.AddHeaders("foo", "loquacious", "x")
	ht.AddRowItems(42, ".", "fred")
	ht.AddRowItems("snerty", "word", "r")
	ht.AddSeparator()
	ht.AddRowItems(" ", true, nil)
	T.Equal(ht.Errors(), nil, "no errors just adding items")

	const shouldTail = `  <thead>
    <tr><th>foo</th><th>loquacious</th><th>x</th></tr>
  </thead>
  <tbody>
    <tr><td>42</td><td>.</td><td>fred</td></tr>
    <tr><td>snerty</td><td>word</td><td>r</td></tr>
    <tr><td> </td><td>true</td><td></td></tr>
  </tbody>
</table>
`

	should := "<table>\n" + shouldTail
	rendered, err := ht.Render()
	T.ExpectSuccess(err, "simple table rendered to HTML")
	T.Equal(ht.Errors(), nil, "no errors accumulated in table through rendering")
	T.Equal(rendered, should, "simple table rendered to HTML correctly")
}
