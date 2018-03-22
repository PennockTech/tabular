// Copyright © 2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

/*
The decoration package provides controls for how text tables are decorated for
rendering.  Unfortunately, this is a mix of functions which need to be
accessible to package API callers and functions which should be internal to
tabular.  Tabular pre-dates the `internal/` system of Go and needs some
refactoring.

The symbols from registry.go and styles.go are stable and part of the public
API.  The rest should probably move to become internal.  If you find yourself
using one of those, and know about it and have read this documentation, please
consider opening an issue so that a future change for v2 will better support you.

Sorry about this mess.
*/
package decoration // import "go.pennock.tech/tabular/texttable/decoration"
