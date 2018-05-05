// Copyright © 2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular_test // import "go.pennock.tech/tabular"

import (
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
)

func TestVersions(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	linesPre := tabular.Versions()
	tabular.LinkerSpecifiedVersion = "testing-1.2.3"
	linesPost := tabular.Versions()
	for _, l := range linesPre {
		t.Log(l)
	}
	for _, l := range linesPost {
		t.Log(l)
	}
}
