// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"fmt"
)

// NoSuchCellError is returned when a table is asked for a cell at some
// coordinates which do not exist within the table.
type NoSuchCellError struct {
	_        struct{}
	Location CellLocation
}

func (e NoSuchCellError) Error() string {
	return fmt.Sprintf("table does not contain cell at coordinates [row %d, col %d]", e.Location.Row, e.Location.Column)
}
