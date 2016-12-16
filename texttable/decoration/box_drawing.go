// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import (
	"reflect"
)

// A Decoration holds the rules for drawing a table's lines.
// Some fields can be inferred from others if left empty.
// We currently use "string" rather than "rune"; if non-empty,
// the contents _must_ render as one terminal cell width.
// Multiple runes are allowed (combining chars, etc).
// TODO (low priority): measure width, handle non-one-cell-width characters.
type Decoration struct {
	/*
		Our model tables:   Minimal table:
		┏━━━┳━━━━━━━┳━━━┓   +---+-------+---+  TopLeft HOuter HTopDown HOuter HTopDown HOuter TopRight
		┃ C ┃ Name  ┃ N ┃   |   |       |   |  VHeader
		┣━━━╇━━━━━━━╇━━━┫   +---+-------+---+  HBLeft HOuter HBCross HOuter HBCross HOuter HBRight
		┃ a │ Funky │ 1 ┃   |   |       |   |  VBodyBorder VBodyInner VBodyInner VBodyBorder
		┃ b │ Hello │ 2 ┃   |   |       |   |
		┠───┼───────┼───┨   +---+-------+---+  LeftBodyRule HRule CrossPiece HRule CrossPiece HRule RightBodyRule
		┃ c │ Final │ 3 ┃   |   |       |   |
		┗━━━┷━━━━━━━┷━━━┛   +---+-------+---+  BottomLeft HOuter BBottomUp HOuter BBottomUp HOuter BottomRight

		┏━━━┯━━━━━━━┯━━━┓   +---+-------+---+  TopLeft HOuter BTopDown HOuter BTopDown HOuter TopRight
		┃ o │ Only  │ 0 ┃   |   |       |   |
		┗━━━┷━━━━━━━┷━━━┛   +---+-------+---+

		No header-only tables; degenerate to:
		┏━━━┳━━━━━━━┳━━━┓
		┃ C ┃ Name  ┃ N ┃
		┣━━━╇━━━━━━━╇━━━┫
		┗━━━┷━━━━━━━┷━━━┛
	*/
	Horizontal string // unused-for-render
	Vertical   string // unused-for-render
	CrossPiece string

	TopDown string // unused-for-render
	VBorder string // unused-for-render

	HOuter        string
	HRule         string
	VHeader       string
	VBodyBorder   string
	VBodyInner    string
	TopLeft       string
	TopRight      string
	BottomLeft    string
	BottomRight   string
	LeftBodyRule  string
	RightBodyRule string
	HTopDown      string
	BTopDown      string
	BBottomUp     string
	HBCross       string
	HBLeft        string
	HBRight       string

	// might change this to non-bool, if we want to control options such as blank line
	// between headers and content, etc.
	isBoxless bool
}

// Populate fills in a table given some key points.
func (d *Decoration) Populate() {
	if d == nil {
		return
	}
	// These defaults deliberately not drawn nice, but make it easier to debug what's missing.
	if d.Horizontal == "" {
		d.Horizontal = "H"
	}
	if d.Vertical == "" {
		d.Vertical = "V"
	}
	if d.CrossPiece == "" {
		d.CrossPiece = "X"
	}
	// and the remaining unused-for-render templates:
	decorateDefaultTo(d, "TopDown", "CrossPiece")
	decorateDefaultTo(d, "VBorder", "Vertical")

	decorateDefaultTo(d, "HOuter", "Horizontal")
	decorateDefaultTo(d, "HRule", "Horizontal")
	decorateDefaultTo(d, "VHeader", "VBorder")
	decorateDefaultTo(d, "VBodyBorder", "VBorder")
	decorateDefaultTo(d, "VBodyInner", "Vertical")

	decorateDefaultTo(d, "TopLeft", "CrossPiece")
	decorateDefaultTo(d, "TopRight", "CrossPiece")
	decorateDefaultTo(d, "BottomLeft", "CrossPiece")
	decorateDefaultTo(d, "BottomRight", "CrossPiece")
	decorateDefaultTo(d, "LeftBodyRule", "CrossPiece")
	decorateDefaultTo(d, "RightBodyRule", "CrossPiece")
	decorateDefaultTo(d, "HTopDown", "TopDown")
	decorateDefaultTo(d, "BTopDown", "TopDown")
	decorateDefaultTo(d, "BBottomUp", "CrossPiece")
	decorateDefaultTo(d, "HBCross", "CrossPiece")
	decorateDefaultTo(d, "HBLeft", "LeftBodyRule")
	decorateDefaultTo(d, "HBRight", "RightBodyRule")
}

func decorateDefaultTo(d *Decoration, toFill, src string) {
	dv := reflect.ValueOf(d)
	if dv.Kind() != reflect.Ptr {
		panic("decoration ptr not a ptr??")
	}
	dev := dv.Elem()
	if dev.Kind() != reflect.Struct {
		panic("decoration ptr doesn't point to struct??")
	}
	target := dev.FieldByName(toFill)
	if target.Kind() == 0 {
		panic("decoration struct doesn't hold target field: " + toFill)
	}
	if target.Len() > 0 {
		return
	}
	origin := dev.FieldByName(src)
	if origin.Kind() == 0 {
		panic("decoration struct doesn't hold origin field: " + src)
	}
	target.Set(origin)
}
