// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// Open questions:
// * should we use interfaces to handle reflow and default to reflow

// Stringer to match fmt's.
// Embedded newlines will be handled to provide lines.
type Stringer interface {
	String() string
}

// GoStringer to match fmt's.
type GoStringer interface {
	GoString() string
}

// Fielder for ability to grab "a row" from a type.
type Fielder interface {
	Fields() []string
}

// AnonFielder to avoid caller having to iterate an []interface{}
// to construct strings; we then just do cell breakdown on each.
type AnonFielder interface {
	AnonFields() []interface{}
}

// TerminalCellWidther should be implemented to override tabular's conception
// of how "wide" a cell's contents are.
type TerminalCellWidther interface {
	TerminalCellWidth() int
}

// Heighter lets an object override how tall it is.
type Heighter interface {
	Height() int
}
