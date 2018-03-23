// Copyright © 2016,2018 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"errors"
	"fmt"
)

// NoSuchCellError is returned when a table is asked for a cell at some
// coordinates which do not exist within the table.
type NoSuchCellError struct {
	_        struct{}
	Location CellLocation
}

func (e NoSuchCellError) Error() string {
	return fmt.Sprintf("tabular: table does not contain cell at coordinates [row %d, col %d]", e.Location.Row, e.Location.Column)
}

// ErrMissingPropertyHolder is returned if you call SetProperty upon a nil object.
var ErrMissingPropertyHolder = errors.New("tabular: missing item upon which to set a property")
