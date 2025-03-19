// Copyright © 2016,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package auto_test // import "go.pennock.tech/tabular/auto"

import (
	"bytes"
	"reflect"
	"regexp"
	"sort"
	"strings"
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

func TestColor(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	tb := auto.New("texttable.utf8-light-curved.red.blue")
	populate(T, tb)
	buf := &bytes.Buffer{}
	err := tb.RenderTo(buf)
	T.ExpectSuccess(err, "rendered to buffer")

	const fg = "\x1b[38;2;255;0;0m"
	const bg = "\x1b[48;2;0;0;255m"
	const rs = "\x1b[m"
	const START = fg + bg
	const STOP = rs
	const BGONLY = rs + bg

	shouldHaveClean := "" +
		"╭────────┬────────────┬──────╮\n" +
		"│ foo    │ loquacious │ x    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│ 42     │ .          │ fred │\n" +
		"│ snerty │ word       │ r    │\n" +
		"├────────┼────────────┼──────┤\n" +
		"│        │ true       │      │\n" +
		"╰────────┴────────────┴──────╯\n"

	expected := "" +
		"START╭────────┬────────────┬──────╮STOP\n" +
		"START│STOP foo    START│STOP loquacious START│STOP x    START│STOP\n" +
		"START├────────┼────────────┼──────┤STOP\n" +
		"START│STOP 42     START│STOP .          START│STOP fred START│STOP\n" +
		"START│STOP snerty START│STOP word       START│STOP r    START│STOP\n" +
		"START├────────┼────────────┼──────┤STOP\n" +
		"START│STOP        START│STOP true       START│STOP      START│STOP\n" +
		"START╰────────┴────────────┴──────╯STOP\n"
	expected = strings.Replace(expected, "START", START, -1)
	expected = strings.Replace(expected, "STOP", STOP, -1)

	have := buf.String()
	T.Equal(have, expected, "got colored table text")

	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	have = re.ReplaceAllString(have, "")

	T.Equal(have, shouldHaveClean, "colored table strips down to expected clean state")

	buf.Reset()
	tt := tb.(*texttable.TextTable)
	tt.SetBGSolid(true)
	err = tb.RenderTo(buf)
	T.ExpectSuccess(err, "rendered to buffer")

	expectedSolid := "" +
		"START╭────────┬────────────┬──────╮STOP\n" +
		"START│BGONLY foo    START│BGONLY loquacious START│BGONLY x    START│STOP\n" +
		"START├────────┼────────────┼──────┤STOP\n" +
		"START│BGONLY 42     START│BGONLY .          START│BGONLY fred START│STOP\n" +
		"START│BGONLY snerty START│BGONLY word       START│BGONLY r    START│STOP\n" +
		"START├────────┼────────────┼──────┤STOP\n" +
		"START│BGONLY        START│BGONLY true       START│BGONLY      START│STOP\n" +
		"START╰────────┴────────────┴──────╯STOP\n"
	expectedSolid = strings.Replace(expectedSolid, "START", START, -1)
	expectedSolid = strings.Replace(expectedSolid, "STOP", STOP, -1)
	expectedSolid = strings.Replace(expectedSolid, "BGONLY", BGONLY, -1)

	have = buf.String()
	T.Equal(have, expectedSolid, "got solid colored table text")
}
