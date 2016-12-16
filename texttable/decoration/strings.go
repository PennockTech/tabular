// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import "strings"

type WidthString struct {
	S string
	W int
}

func (ws WidthString) WithinWidth(available int) string {
	if ws.W < 0 {
		return strings.Repeat(" ", available)
	}
	pad := available - ws.W
	if pad < 0 {
		pad = 0
	}
	return ws.S + strings.Repeat(" ", pad)
}
