// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable

import (
	"fmt"

	"github.com/PennockTech/tabular/texttable/decoration"
)

// SetDecoration sets a Decoration type for rendering a table.  The caller
// must provide the decoration object.
func (t *TextTable) SetDecoration(decor decoration.Decoration) *TextTable {
	t.decor = decor
	return t
}

// SetDecorationNamed selects a decoration by name, from the decoration package.
// It returns an error if the name is not known.
// If the name is not known, this is still registered in the TextTable, so that
// later attempts to render it will fail, instead of succeeding with unexpected
// output.
func (t *TextTable) SetDecorationNamed(n string) (*TextTable, error) {
	d := decoration.Named(n)
	t.decor = d
	if d == decoration.EmptyDecoration {
		return t, fmt.Errorf("unknown decoration name %q", n)
	}
	return t, nil
}
