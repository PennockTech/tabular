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

// GetAligmentOffset returns an AlignmentOffset for a given cell, which is a
// terminal display cell offset from the left.  This is suitable for outputs
// where we control alignments via spaces etc, but not suitable for use with
// HTML which has its own capable alignment model.
//
// This function iterates through a cell, its row, its column and its table
// looking for alignment controls.
func GetAligmentOffset(cell *Cell) AlignmentOffset {
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

// GetStaticAlignment gets the static alignment for one specific property owner.
// It does not iterate through sources, it does not handle AlignmentFunc.
// It is suitable for use with rich display alignment controls such as those in HTML.
func GetStaticAlignment(po PropertyOwner) (StaticAlignment, bool) {
	sI := po.GetProperty(propAlignmentStatic)
	if sI == nil {
		return ALIGN_LEFT, false
	}
	s, ok := sI.(StaticAlignment)
	if !ok {
		return ALIGN_LEFT, false
	}
	return s, true
}

// SetAlignmentStatic can be called with any PropertyOwner (cell, row, etc)
// to specify an explicit StaticAlignment property to record.
func SetAlignmentStatic(po PropertyOwner, a StaticAlignment) {
	po.SetProperty(propAlignmentStatic, a)
}

// SetAlignmentFunc can be called with any PropertyOwner (cell, row, etc)
// to specify an AlignmentFunc to record as a property; that function is called
// upon individual cells.  It is less supported than static alignments, but
// potentially more powerful.
func SetAlignmentFunc(po PropertyOwner, f AlignmentFunc) {
	po.SetProperty(propAlignmentFunc, f)
}
