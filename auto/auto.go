// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package auto // import "go.pennock.tech/tabular/auto"

import (
	"io"
	"sort"
	"strings"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/csv"
	"go.pennock.tech/tabular/html"
	"go.pennock.tech/tabular/json"
	"go.pennock.tech/tabular/markdown"
	"go.pennock.tech/tabular/texttable"
	"go.pennock.tech/tabular/texttable/decoration"
)

// Any table sub-package which can be rendered should meet this
// interface.
type RenderTable interface {
	tabular.Table
	Render() (string, error)
	RenderTo(io.Writer) error
}

// Wrap takes a tabular.Table and a style string.  While Wrap returns
// the RenderTable interface, the type underlying that interface is
// dependent upon the contents of the style string.
//
// The style is a dot (period) separated sequence of fields, but will usually
// just be a single field with no dot.  The first section is either the name of
// a sub-package, or a known registered decoration of texttable.
// Wrap(t, "texttable.utf8-light") is the same as Wrap(t, "utf8-light").
//
// The style sections after the first are interpreted dependent upon the first
// section and not yet locked down by API.
func Wrap(t tabular.Table, style string) RenderTable {
	sections := strings.Split(style, ".")
	switch strings.ToLower(sections[0]) {
	case "csv":
		return csv.Wrap(t)
	case "html":
		// TODO: do we want to take `html.caption="foo.bar".class="a b c".whatever ?
		return html.Wrap(t)
	case "markdown":
		return markdown.Wrap(t)
	case "json":
		return json.Wrap(t)
	case "texttable":
		tt := texttable.Wrap(t)
		if len(sections) > 1 {
			tt.SetDecorationNamed(sections[1])
		}
		return tt
	default:
		tt := texttable.Wrap(t)
		tt.SetDecorationNamed(sections[0])
		return tt
	}
}

// New creates a new tabular.Table and Wrap()s it.
func New(style string) RenderTable {
	return Wrap(tabular.New(), style)
}

// ListStyles returns a sorted list of strings, where each string is a valid
// input to the New function.  The list is not exhaustive, sub-styles may be
// either included or omitted; the results should be sufficiently
// representative to be used by help systems for listing available styles,
// without having to go into every possible mutation.
func ListStyles() []string {
	l := decoration.RegisteredDecorationNames()
	l = append(l, "csv", "html", "json", "markdown")
	sort.Strings(l)
	return l
}

// Render takes a tabular.Table and a style and creates a default-options
// object of the type indicated by the style, with any tuning options from that
// style applied, and then calls the Render method upon it.
func Render(t tabular.Table, style string) (string, error) {
	return Wrap(t, style).Render()
}

// RenderTo takes a tabular.Table and a style and creates a default-options
// object of the type indicated by the style, with any tuning options from that
// style applied, and then calls the RenderTo method upon it.  API considers
// the style to be a decoration/refinement and places that parameter last.
func RenderTo(t tabular.Table, w io.Writer, style string) error {
	return Wrap(t, style).RenderTo(w)
}
