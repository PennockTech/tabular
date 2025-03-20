// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package auto // import "go.pennock.tech/tabular/auto"

import (
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/color"
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
	var rt RenderTable
	sections := strings.Split(style, ".")
	switch strings.ToLower(sections[0]) {
	case "csv":
		rt = csv.Wrap(t)
	case "html":
		// TODO: do we want to take `html.caption="foo.bar".class="a b c".whatever ?
		rt = html.Wrap(t)
	case "markdown":
		rt = markdown.Wrap(t)
	case "json":
		rt = json.Wrap(t)
	case "texttable":
		tt := texttable.Wrap(t)
		setColorsOrDecorationsFromSections(tt, sections[1:])
		rt = tt
	default:
		tt := texttable.Wrap(t)
		setColorsOrDecorationsFromSections(tt, sections)
		rt = tt
	}
	return rt
}

var reColorHex *regexp.Regexp

func init() {
	reColorHex = regexp.MustCompile(`^#?([0-9A-Fa-f]{2})([0-9A-Fa-f]{2})([0-9A-Fa-f]{2})\z`)
}

func setColorsOrDecorationsFromSections(tt *texttable.TextTable, sections []string) {
	doneFG, doneBG, doneDecoration := false, false, false
	for _, section := range sections {
		var (
			c                color.Color
			err              error
			red, green, blue uint64
			haveColor        bool
		)
		haveColor = false
		if c, err = color.ByHTMLNamedColor(section); err == nil {
			haveColor = true
		} else if m := reColorHex.FindStringSubmatch(section); m != nil {
			red, err = strconv.ParseUint(m[1], 16, 8)
			if err == nil {
				green, err = strconv.ParseUint(m[2], 16, 8)
			}
			if err == nil {
				blue, err = strconv.ParseUint(m[3], 16, 8)
			}
			if err == nil {
				c = color.RGB24(uint8(red), uint8(green), uint8(blue))
				haveColor = true
			}
		} else if section == "solid" {
			tt.SetBGSolid(true)
		} else if !doneDecoration {
			tt.SetDecorationNamed(section)
			doneDecoration = true
		}
		if haveColor {
			if doneFG && doneBG {
				continue
			} else if doneFG {
				tt.SetBGColor(c)
				doneBG = true
			} else {
				tt.SetFGColor(c)
				doneFG = true
			}
		}
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
