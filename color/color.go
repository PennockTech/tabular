// Copyright © 2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package color // import "go.pennock.tech/tabular/color"

import (
	"errors"
	"strconv"
	"strings"
)

var ErrUnknownColor = errors.New("tabular: unknown color")

// AnsiColor is a special case of Color, supporting the idea that we might allow
// something to emit "classic" ANSI escape sequences and it's still a color.
// (We don't, at time of writing.)
type AnsiColor interface {
	AnsiEscapeFG() string
}

// A Color represents an RGB color.
type Color struct {
	red, green, blue uint8
}

const nibbleToHex = "0123456789ABCDEF"
const ansiCSI = "\x1B["
const ResetColor = ansiCSI + "m"

func writeUint8Hex2(buf []byte, n uint8) {
	buf[0] = byte(nibbleToHex[n/16])
	buf[1] = byte(nibbleToHex[n%16])
}

// RGB24 makes a color from uint8 R,G,B values
func RGB24(red, green, blue uint8) Color {
	return Color{red, green, blue}
}

// ByHTMLNamedColor returns a color from the HTML color name, or returns ErrUnknownColor
func ByHTMLNamedColor(name string) (Color, error) {
	lname := strings.ToLower(name)
	if cVal, ok := htmlNamedColors[lname]; ok {
		return Color{
			red:   uint8((cVal / 0x10000) & 0xFF),
			green: uint8((cVal / 0x100) & 0xFF),
			blue:  uint8(cVal & 0xFF),
		}, nil
	}
	return Color{0, 0, 0}, ErrUnknownColor
}

// RGBHex returns the color as a 6 octet string RRGGBB
func (c Color) RGBHex() string {
	buf := make([]byte, 6)
	writeUint8Hex2(buf[0:2], c.red)
	writeUint8Hex2(buf[2:4], c.green)
	writeUint8Hex2(buf[4:6], c.blue)
	return string(buf)
}

// HTML returns the color as an HTML color sequence
func (c Color) HTML() string {
	buf := make([]byte, 7)
	buf[0] = '#'
	writeUint8Hex2(buf[1:3], c.red)
	writeUint8Hex2(buf[3:5], c.green)
	writeUint8Hex2(buf[5:7], c.blue)
	return string(buf)
}

// AnsiEscapeFG returns the color as an escape sequence setting the foreground
// Direct Color in RGB Space.
func (c Color) AnsiEscapeFG() string {
	// xterm-ctlseqs.ms gives a reference to ISO-8613-6 but it's a pay-to-play
	// "standard", which because of the pay status means it wasn't available to
	// many implementors and support for the "correct" form is not as
	// widespread as support for an "incorrect" form.
	// I used a search engine to try to find a reference to figure out which of these I should use:
	//   \e[38:2:<dummy>:<red>:<green>:<blue>m
	//   \e[38;2;<red>;<green>;<blue>m
	// and <https://chadaustin.me/2024/01/truecolor-terminal-emacs/> was very helpful.
	// The colon form is per-ISO but support is not as widespread as the
	// "incorrect" semi-colon form (which is the only form I'd seen ever
	// actually used).  I do understand the point about ambiguity, as being
	// something I'd wondered about myself when decoding this stuff in the
	// past.
	return ansiCSI + "38;2;" + strconv.Itoa(int(c.red)) + ";" + strconv.Itoa(int(c.green)) + ";" + strconv.Itoa(int(c.blue)) + "m"
}

// AnsiEscapeBG returns the color as an escape sequence setting the foreground
// Direct Color in RGB Space.
func (c Color) AnsiEscapeBG() string {
	return ansiCSI + "48;2;" + strconv.Itoa(int(c.red)) + ";" + strconv.Itoa(int(c.green)) + ";" + strconv.Itoa(int(c.blue)) + "m"
}
