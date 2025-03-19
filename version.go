// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

// We are a library, not a top-level binary, so we can't depend upon any
// particular top-level linker action specifying versions.  That said, we
// can provide hooks for applications to cooperate, if they so choose.
//
// *Clients* please consider invoking the `.version` shell-script inside
// this repo and passing it to the Go linker.
// See github.com/philpennock/character for an example.
var LinkerSpecifiedVersion string

const packageVersionName = "tabular"
const APIVersion string = "1.3"

func Versions() []string {
	vl := make([]string, 0, 2)
	if LinkerSpecifiedVersion != "" {
		vl = append(vl, packageVersionName+": build-time specified: "+LinkerSpecifiedVersion)
	}
	vl = append(vl, packageVersionName+": API version: "+APIVersion)
	return vl
}
