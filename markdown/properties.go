// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package markdown

import (
	"errors"

	"github.com/PennockTech/tabular"
)

type propertyKey struct {
	name string
}

func (p *propertyKey) String() string { return "markdowntable property keyid " + p.name }

var (
	propWidth = &propertyKey{"width"}
)

type width struct {
	cellWidth int
}

type widthSetter struct{}

var ErrNotCellProperties = errors.New("markdowntable: width-set: not given a cell")

func (_ widthSetter) UpdateProperties(po tabular.PropertyOwner) error {
	cell, ok := po.(*tabular.Cell)
	if !ok {
		return ErrNotCellProperties
	}

	// TODO: move these calculations into separate common sub-package?
	dims := width{
		cellWidth: cell.TerminalCellWidth(),
	}

	po.SetProperty(propWidth, dims)
	return nil
}

func CellPropertyExtractWidth(cell *tabular.Cell) int {
	dimsI := cell.GetProperty(propWidth)
	if dimsI == nil {
		return 0
	}
	dims, ok := dimsI.(width)
	if ok {
		return dims.cellWidth
	}
	return 0
}
