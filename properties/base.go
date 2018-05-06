// Copyright © 2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package properties // import "go.pennock.tech/tabular/properties"

type propertyKey struct {
	name string
}

func (p *propertyKey) String() string { return "miscellaneous property keyid " + p.name }

var (
	Skipable = &propertyKey{"skipable"}
)
