// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular

import (
	"fmt"
	"strconv"
	"strings"
)

// StaticAligment indicates how cells should be aligned; it can be set as a
// property upon a cell, row, column or table and lookups search for an
// alignment in that order.
type StaticAlignment int

const (
	ALIGN_LEFT StaticAlignment = iota
	ALIGN_CENTER
	ALIGN_RIGHT
	ALIGN_PERIOD
	ALIGN_COMMA
)

// AlignmentOffset is a uint and says how many terminal display cells from the
// left the alignment offset is.
// 0 is the left-most character column (not cell column).
type AlignmentOffset uint

// An AlignmentFunc can be registered on a property owner to dynamically
// determine the AlignmentOffset; you can register your own functions, else
// any StaticAlignment found will be used, else no alignment (thus ALIGN_LEFT).
type AlignmentFunc func(*Cell) AlignmentOffset

type alignmentPropertyKey int

var (
	propAlignmentStatic = alignmentPropertyKey(1)
	propAlignmentFunc   = alignmentPropertyKey(2)
)

func (p alignmentPropertyKey) String() string {
	return "alignment property keyid " + strconv.Itoa(int(p))
}

func GetAligment(cell *Cell) AlignmentOffset {
	var errRecv ErrorReceiver
	if cell.inRow != nil {
		errRecv = cell.inRow
	}
	width := cell.TerminalCellWidth()
	if width < 0 {
		width = len(cell.String())
	}

	for _, src := range []struct {
		po    PropertyOwner
		label string
	}{
		{cell, fmt.Sprintf("cell %v", cell.Location())},
		{cell.inRow, fmt.Sprintf("row %v", cell.inRow.rowNum)},
		{cell.columnOfTable(), fmt.Sprintf("column %v", cell.columnNum)},
		{cell.tablePtr(), "table"},
	} {
		if src.po == nil {
			continue
		}
		fI := src.po.GetProperty(propAlignmentFunc)
		if fI != nil {
			f, ok := fI.(AlignmentFunc)
			if !ok {
				if errRecv != nil {
					errRecv.AddError(fmt.Errorf("%s has bad AlignmentFunc property", src.label))
				}
				return AlignmentOffset(0)
			}
			return f(cell)
		}
		sI := src.po.GetProperty(propAlignmentStatic)
		if sI != nil {
			s, ok := sI.(StaticAlignment)
			if !ok {
				if errRecv != nil {
					errRecv.AddError(fmt.Errorf("%s has bad AlignmentStatic property", src.label))
				}
				return AlignmentOffset(0)
			}
			switch s {
			case ALIGN_LEFT:
				return AlignmentOffset(0)
			case ALIGN_CENTER:
				return AlignmentOffset(width / 2)
			case ALIGN_RIGHT:
				return AlignmentOffset(width)

			case ALIGN_PERIOD:
				if cell.Height() > 1 {
					return AlignmentOffset(0)
				}
				offset := strings.LastIndex(cell.String(), ".")
				if offset == -1 {
					offset = 0
				}
				return AlignmentOffset(offset)

			case ALIGN_COMMA:
				if cell.Height() > 1 {
					return AlignmentOffset(0)
				}
				offset := strings.LastIndex(cell.String(), ",")
				if offset == -1 {
					offset = 0
				}
				return AlignmentOffset(offset)

			default:
				if errRecv != nil {
					errRecv.AddError(fmt.Errorf("%s has unknown AlignmentStatic value %v", src.label, s))
				}
				return AlignmentOffset(0)
			}
		}
	}
	return AlignmentOffset(0)
}
