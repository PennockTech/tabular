// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration

import (
	"testing"

	"github.com/liquidgecka/testlib"
)

// This is a bit silly
func TestRegisteredNames(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	seen := make(map[string]struct{}, 10)

	for testI, tuple := range []struct {
		name     string
		callable func() Decoration
	}{
		{"ascii-simple", ASCIIBoxSimple},
		{"none", NoBox},
		{"utf8-light", UTF8BoxLight},
		{"utf8-light-curved", UTF8BoxLightCurved},
		{"utf8-heavy", UTF8BoxHeavy},
		{"utf8-double", UTF8BoxDouble},
	} {
		have := Named(tuple.name)
		want := tuple.callable()
		T.Equalf(have, want, "[%d] decoration equality", testI)
		seen[tuple.name] = struct{}{}
	}

	missing := make([]string, 0, 10)
	for _, n := range RegisteredDecorationNames() {
		if _, ok := seen[n]; !ok {
			missing = append(missing, n)
		}
	}
	if len(missing) > 0 {
		T.Fatalf("missing tests for registered decorations: %v", missing)
	}
}
