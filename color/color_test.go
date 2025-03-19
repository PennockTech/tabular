// Copyright © 2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package color_test

import (
	"strings"
	"testing"

	"go.pennock.tech/tabular/color"

	"github.com/liquidgecka/testlib"
)

func TestColorRGBHex(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	type rgbs struct {
		red, green, blue uint8
		hex              string
		terminal         string
	}
	for n, item := range []rgbs{
		{0, 0, 0, "000000", "\x1B[38;2;0;0;0m"},
		{255, 255, 255, "FFFFFF", "\x1B[38;2;255;255;255m"},
		{255, 10, 16, "FF0A10", "\x1B[38;2;255;10;16m"},
		{10, 255, 16, "0AFF10", "\x1B[38;2;10;255;16m"},
		{10, 16, 255, "0A10FF", "\x1B[38;2;10;16;255m"},
		{42, 87, 200, "2A57C8", "\x1B[38;2;42;87;200m"},
	} {
		c := color.RGB24(item.red, item.green, item.blue)
		T.Equalf(c.RGBHex(), item.hex, "RGBHex row %d", n)
		T.Equalf(c.AnsiEscapeFG(), item.terminal, "AnsiEscapeFG row %d", n)
		T.Equalf(c.AnsiEscapeBG(), strings.Replace(item.terminal, "[38;", "[48;", 1), "AnsiEscapeBG row %d", n)
		T.Equalf(c.HTML(), "#"+item.hex, "HTML row %d", n)
	}
}

func TestColorHTML(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	for n, item := range []struct{ name, html string }{
		{"black", "#000000"},
		{"white", "#FFFFFF"},
		{"cyan", "#00FFFF"},
		{"royalblue", "#4169E1"},
	} {
		c, err := color.ByHTMLNamedColor(item.name)
		T.ExpectSuccessf(err, "row %d color %q", n, item.name)
		T.Equalf(c.HTML(), item.html, "row %d color %q", n, item.name)
	}
}
