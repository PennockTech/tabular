// Copyright © 2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package align // import "go.pennock.tech/tabular/properties/align"

type propertyKey struct {
	name string
}

func (p *propertyKey) String() string { return "alignment property keyid " + p.name }

var (
	PropertyType = &propertyKey{"type"}
)

type Alignment interface {
	isAlignment() struct{}
}

type alignSimple struct {
	which int
}

func (l alignSimple) isAlignment() struct{} { return struct{}{} }

var (
	Left   = alignSimple{1}
	Right  = alignSimple{2}
	Center = alignSimple{3}
)

// TestingInvalidAlignment returns an alignment which should not be handled by
// code and may be used to exercise default handling, such as panics.
func TestingInvalidAlignment() Alignment {
	return alignSimple{99999}
}
