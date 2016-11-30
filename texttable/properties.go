// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package texttable

import (
	"errors"

	"github.com/PennockTech/tabular"
	"github.com/PennockTech/tabular/length"
	"github.com/PennockTech/tabular/texttable/decoration"
)

type propertyKey struct {
	name string
}

func (p *propertyKey) String() string { return "texttable property keyid " + p.name }

var (
	propDimensions  = &propertyKey{"dimensions"}
	propLinesWidths = &propertyKey{"lineswidths"}
)

type dimensions struct {
	cellWidth int
	height    int
}

type dimensionSetter struct{}

var ErrNotCellProperties = errors.New("texttable: dimensions-set: not given a cell")

func (_ dimensionSetter) UpdateProperties(po tabular.PropertyOwner) error {
	cell, ok := po.(*tabular.Cell)
	if !ok {
		return ErrNotCellProperties
	}

	// TODO: move these calculations into separate common sub-package?
	dims := dimensions{
		cellWidth: cell.TerminalCellWidth(),
		height:    cell.Height(),
	}

	linesWidths := make([]decoration.WidthString, dims.height)
	for i, l := range cell.Lines() {
		linesWidths[i] = decoration.WidthString{
			S: l,
			W: length.StringCells(l),
		}
	}

	po.SetProperty(propDimensions, dims)
	po.SetProperty(propLinesWidths, linesWidths)
	return nil
}

func CellPropertyExtractDimensions(cell *tabular.Cell) dimensions {
	dimsI := cell.GetProperty(propDimensions)
	if dimsI == nil {
		return dimensions{0, 0}
	}
	dims, ok := dimsI.(dimensions)
	if ok {
		return dims
	}
	return dimensions{0, 0}
}

func CellPropertyExtractLinesWidths(cell *tabular.Cell) []decoration.WidthString {
	linesI := cell.GetProperty(propLinesWidths)
	if linesI == nil {
		return nil
	}
	lines, ok := linesI.([]decoration.WidthString)
	if ok {
		return lines
	}
	return nil
}
