// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package decoration // import "go.pennock.tech/tabular/texttable/decoration"

import (
	"sort"
	"sync"
)

// EmptyDecoration is a non-decoration and is returned for named lookups for an
// unknown name.
var EmptyDecoration = Decoration{}

var registry struct {
	sync.Mutex
	table map[string]Decoration
}

// provide a couple of entries for someone to register more
const registryApproxSize = 8

func init() {
	registry.Lock()
	registry.table = make(map[string]Decoration, registryApproxSize)
	registry.Unlock()
}

// RegisterDecorationName declares a given name to provide a given style.
// Existing entries may be overwritten.
func RegisterDecorationName(name string, decor Decoration) {
	registry.Lock()
	registry.table[name] = decor
	registry.Unlock()
}

// RegisteredDecorationNames returns a sorted list of registered decoration
// names.
func RegisteredDecorationNames() []string {
	registry.Lock()
	defer registry.Unlock()
	a := make([]string, len(registry.table))
	i := 0
	for k := range registry.table {
		a[i] = k
		i++
	}
	sort.Strings(a)
	return a
}

// Named returns the named decoration, or EmptyDecoration if the name is not known.
// The name is a simple string instead of typed, so as part of the point of this API
// is that callers don't need to explicitly import this package.
func Named(n string) Decoration {
	registry.Lock()
	d, ok := registry.table[n]
	registry.Unlock()
	if ok {
		return d
	}
	return EmptyDecoration
}
