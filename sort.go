// Copyright © 2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"errors"
	"sort"
	"strconv"
)

type ErrorNoSuchColumn string

func (e ErrorNoSuchColumn) Error() string { return "no such column " + strconv.Quote(string(e)) }

type ErrorColumnOutOfRange int

func (e ErrorColumnOutOfRange) Error() string {
	return "column " + strconv.Itoa(int(e)) + " out of range"
}

var ErrNoColumnHeaders = errors.New("no headers have defined columns")

type SortOrder int

const (
	SORT_ASC SortOrder = iota + 1
	SORT_DESC
)

func (so SortOrder) String() string {
	switch so {
	case SORT_ASC:
		return "ascending"
	case SORT_DESC:
		return "descending"
	default:
		panic("unhandled sort order for String")
	}
}

// SortByNamedColumn performs an in-place row-sort.
func (t *ATable) SortByNamedColumn(name string, order SortOrder) error {
	if t.columnNames == nil {
		return ErrNoColumnHeaders
	}
	columnNumber, ok := t.columnNames[name]
	if !ok {
		return ErrorNoSuchColumn(name)
	}
	return t.SortByColumnNumber(columnNumber, order)
}

type tableSorter struct {
	tb     *ATable
	column int
	order  SortOrder
}

func (t tableSorter) Len() int { return t.tb.NRows() }
func (t tableSorter) Swap(i, j int) {
	t.tb.rows[i], t.tb.rows[j] = t.tb.rows[j], t.tb.rows[i]
}
func (t tableSorter) Less(i, j int) bool {
	if t.order == SORT_ASC {
		return t.tb.rows[i].cells[t.column].LessThan(&t.tb.rows[j].cells[t.column])
	} else if t.order == SORT_DESC {
		return t.tb.rows[j].cells[t.column].LessThan(&t.tb.rows[i].cells[t.column])
	} else {
		panic("unhandled order in tableSorter.Less")
	}
}

func (t *ATable) SortByColumnNumber(sortCol int, order SortOrder) error {
	if sortCol < 0 || sortCol >= t.nColumns {
		return ErrorColumnOutOfRange(sortCol)
	}
	ts := tableSorter{t, sortCol, order}
	sort.Sort(ts)
	return nil
}
