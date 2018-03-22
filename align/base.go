// Copyright © 2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package align // import "go.pennock.tech/tabular/align"

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

type alignLeft struct{}

func (l alignLeft) isAlignment() struct{} { return struct{}{} }

type alignRight struct{}

func (l alignRight) isAlignment() struct{} { return struct{}{} }

var (
	Left  = alignLeft{}
	Right = alignRight{}
)
