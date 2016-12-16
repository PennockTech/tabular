// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package auto_test // import "go.pennock.tech/tabular/auto"

import (
	"bytes"
	"reflect"
	"sort"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/auto"
	"go.pennock.tech/tabular/csv"
	"go.pennock.tech/tabular/html"
	"go.pennock.tech/tabular/markdown"
	"go.pennock.tech/tabular/texttable"
)

func populate(T *testlib.T, tb tabular.Table) {
	tb.AddHeaders("foo", "loquacious", "x")
	tb.AddRowItems(42, ".", "fred")
	tb.AddRowItems("snerty", "word", "r")
	tb.AddSeparator()
	tb.AddRowItems(" ", true, nil)
	T.Equal(tb.Errors(), nil, "no errors just adding items")
}

func TestNewBasic(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	bare := tabular.New()

	tac := auto.New("csv")
	tdc := csv.New()

	tah := auto.New("html")
	tdh := html.New()

	tam := auto.New("markdown")
	tdm := markdown.New()

	tat := auto.New("texttable")
	tdt := texttable.New()

	taa1 := auto.New("ascii-simple")
	taa2 := auto.New("texttable.ascii-simple")
	tda := texttable.New()
	tda.SetDecorationNamed("ascii-simple")

	for _, tb := range []tabular.Table{
		bare, tac, tdc, tah, tdh, tam, tdm, tat, tdt, taa1, taa2, tda,
	} {
		populate(T, tb)
	}

	// This also serves to confirm that the tables created directly satisfy our interface.
	for i, tuple := range []struct {
		a, b  auto.RenderTable
		style string
	}{
		{tac, tdc, "csv"},
		{tah, tdh, "html"},
		{tam, tdm, "markdown"},
		{tat, tdt, "texttable"},
		{taa1, tda, "ascii-simple"},
		{taa2, tda, "texttable.ascii-simple"},
		{taa1, taa2, "texttable.ascii-simple"},
	} {
		typeA := reflect.TypeOf(tuple.a).String()
		typeB := reflect.TypeOf(tuple.b).String()
		T.Equalf(typeA, typeB, "auto [%d] tuple's pairs of identical type", i)

		ra, err := tuple.a.Render()
		T.ExpectSuccessf(err, "auto [%d] tuple.a rendered without error", i)

		bufB := &bytes.Buffer{}
		err = tuple.b.RenderTo(bufB)
		T.ExpectSuccessf(err, "auto [%d] tuple.b rendered without error", i)

		T.Equalf(ra, bufB.String(), "auto [%d] tuple renders identically", i)

		if _, ok := tuple.b.(auto.RenderTable); !ok {
			T.Fatalf("auto [%d] tuple.b does not satisfy RenderTable interface", i)
		}

		wrapped := auto.Wrap(bare, tuple.style)
		typeW := reflect.TypeOf(wrapped).String()
		T.Equalf(typeW, typeB, "auto [%d] wrapping of bare to %q of identical type", i, tuple.style)
		bufW := &bytes.Buffer{}
		err = wrapped.RenderTo(bufW)
		T.ExpectSuccessf(err, "auto [%d] bare-wrapped %q rendered without error", i, tuple.style)

		T.Equalf(ra, bufW.String(), "auto [%d] bare-wrapped %q rendered identically", i, tuple.style)
	}
}

func TestListStyles(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	available := auto.ListStyles()
	if !sort.StringsAreSorted(available) {
		T.Error("ListStyles() strings not sorted")
	}

	have := make(map[string]struct{}, len(available))
	for _, h := range available {
		if _, already := have[h]; already {
			T.Errorf("ListStyles() duplicated result: %q", h)
		}
		have[h] = struct{}{}
	}

	for _, expect := range []string{"csv", "html", "markdown", "ascii-simple", "utf8-light"} {
		if _, ok := have[expect]; !ok {
			T.Errorf("ListStyles() missing expected value %q", expect)
		}
	}
}
