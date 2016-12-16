// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

// The const-style names provide for language error-checking for those who
// explicitly import this package.  Both this and just passing a string
// are supported.  Package clients are allowed to register their own names,
// for which no const will be available here.  Which to use is a trade-off for
// the client to choose.
const (
	D_ASCII_SIMPLE      = "ascii-simple"
	D_NONE              = "none"
	D_UTF8_LIGHT        = "utf8-light"
	D_UTF8_LIGHT_CURVED = "utf8-light-curved"
	D_UTF8_HEAVY        = "utf8-heavy"
	D_UTF8_DOUBLE       = "utf8-double"
)

// ASCIIBoxSimple provides traditionalist boxes drawn with ASCII characters.
func ASCIIBoxSimple() Decoration {
	d := Decoration{
		Horizontal: "-",
		Vertical:   "|",
		CrossPiece: "+",
	}
	d.Populate()
	return d
}

func init() { RegisterDecorationName(D_ASCII_SIMPLE, ASCIIBoxSimple()) }

// NoBox ("none") does not draw lines.
func NoBox() Decoration {
	d := Decoration{
		// isBoxless distinguishes us from the empty decoration error type
		isBoxless: true,
	}
	// Do not populate
	return d
}

func init() { RegisterDecorationName(D_NONE, NoBox()) }

// UTF8BoxLight provides boxes drawn from the BOX DRAWINGS LIGHT characters.
func UTF8BoxLight() Decoration {
	d := Decoration{
		Horizontal:    "─",
		Vertical:      "│",
		CrossPiece:    "┼",
		TopDown:       "┬",
		BBottomUp:     "┴",
		TopLeft:       "┌",
		BottomLeft:    "└",
		TopRight:      "┐",
		BottomRight:   "┘",
		LeftBodyRule:  "├",
		RightBodyRule: "┤",
	}
	d.Populate()
	return d
}

func init() { RegisterDecorationName(D_UTF8_LIGHT, UTF8BoxLight()) }

// UTF8BoxLightCurved is very similar to UTF8BoxLight but rounds the corners.
func UTF8BoxLightCurved() Decoration {
	d := UTF8BoxLight()
	d.TopLeft = "╭"
	d.BottomLeft = "╰"
	d.TopRight = "╮"
	d.BottomRight = "╯"
	return d
}

func init() { RegisterDecorationName(D_UTF8_LIGHT_CURVED, UTF8BoxLightCurved()) }

// UTF8BoxHeavy provides boxes with a combination of heavy and light lines around borders and rules.
func UTF8BoxHeavy() Decoration {
	// ┏━━━┳━━━━━━━┳━━━┓
	// ┃ C ┃ Name  ┃ N ┃
	// ┣━━━╇━━━━━━━╇━━━┫
	// ┃ a │ Funky │ 1 ┃
	// ┃ b │ Hello │ 2 ┃
	// ┠───┼───────┼───┨
	// ┃ c │ Final │ 3 ┃
	// ┗━━━┷━━━━━━━┷━━━┛
	d := Decoration{
		Horizontal:    "─",
		Vertical:      "│",
		CrossPiece:    "┼",
		VBorder:       "┃",
		HOuter:        "━",
		HTopDown:      "┳",
		BTopDown:      "┯",
		BBottomUp:     "┷",
		TopLeft:       "┏",
		BottomLeft:    "┗",
		TopRight:      "┓",
		BottomRight:   "┛",
		LeftBodyRule:  "┠",
		RightBodyRule: "┨",
		HBLeft:        "┣",
		HBCross:       "╇",
		HBRight:       "┫",
	}
	d.Populate()
	return d
}

func init() { RegisterDecorationName(D_UTF8_HEAVY, UTF8BoxHeavy()) }

// There are no "BOX DRAWINGS HEAVY ARC DOWN AND RIGHT" etc characters
// So no UTF8BoxHeavyCurved

// UTF8BoxDouble is similar to UTF8BoxHeavy but uses doubled lines instead of heavy lines.
// Note though that because of a missing character in Unicode, this should probably not be used
// for tables with headers (the interior cross-piece along the header/body boundary will not match up).
func UTF8BoxDouble() Decoration {
	d := Decoration{
		Horizontal:    "─",
		Vertical:      "│",
		CrossPiece:    "┼",
		VBorder:       "║",
		HOuter:        "═",
		HTopDown:      "╦",
		BTopDown:      "╤",
		BBottomUp:     "╧",
		TopLeft:       "╔",
		BottomLeft:    "╚",
		TopRight:      "╗",
		BottomRight:   "╝",
		LeftBodyRule:  "╟",
		RightBodyRule: "╢",
		HBLeft:        "╠",
		HBCross:       "╪", // BROKEN: Unicode missing "BOX DRAWINGS DOWN SINGLE AND UP HORIZONTAL DOUBLE"?
		HBRight:       "╣",
		// Do the doubling lines really not have an analogy to "BOX DRAWINGS DOWN LIGHT AND UP HORIZONTAL HEAVY"
	}
	d.Populate()
	return d
}

func init() { RegisterDecorationName(D_UTF8_DOUBLE, UTF8BoxDouble()) }
