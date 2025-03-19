// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable // import "go.pennock.tech/tabular/texttable"

import (
	"go.pennock.tech/tabular"
	"go.pennock.tech/tabular/color"
	"go.pennock.tech/tabular/texttable/decoration"
)

type colorFlags uint8

const (
	colorBGSolid colorFlags = 1 << iota
)

// A TextTable wraps a tabular.Table to act as the render control for
// tabular output to a fixed-cell grid system such as a terminal emulator
// (in the Unix TTY model).
type TextTable struct {
	tabular.Table

	decor   decoration.Decoration
	fgcolor *color.Color
	bgcolor *color.Color
	bgflags colorFlags
}

// Wrap returns a TextTable rendering object for the given tabular.Table
func Wrap(t tabular.Table) *TextTable {
	var ds dimensionSetter
	t.RegisterPropertyCallback(t, tabular.CB_AT_RENDER, tabular.CB_ON_CELL, ds)
	return &TextTable{
		Table: t,
		decor: decoration.UTF8BoxHeavy(),
	}
}

// New returns a TextTable with a new tabular.Table inside it.
func New() *TextTable {
	return Wrap(tabular.New())
}
