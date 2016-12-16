// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular_test // import "go.pennock.tech/tabular"

import (
	"strings"
	"testing"

	"github.com/liquidgecka/testlib"

	"go.pennock.tech/tabular"
)

func TestErrorNoSuchCell(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	e := tabular.NoSuchCellError{Location: tabular.CellLocation{Row: 3, Column: 2}}
	T.ExpectError(e, "no such cell is an error")
	needle := "table does not contain cell at coordinates"
	T.Equalf(strings.Contains(e.Error(), needle), true, "got an error message containing needle %q", needle)
}

func TestErrorContainer(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	ec := tabular.NewErrorContainer()
	T.Equal(ec.Errors(), nil, "empty error-container has no errors, is nil")

	var noErr error
	oneErr := tabular.NoSuchCellError{Location: tabular.CellLocation{Row: 1, Column: 1}}

	ec.AddError(tabular.NoSuchCellError{Location: tabular.CellLocation{Row: 3, Column: 2}})
	ec.AddError(noErr)
	ec.AddError(tabular.NoSuchCellError{Location: tabular.CellLocation{Row: 5, Column: 7}})

	el := ec.Errors()
	T.Equal(len(el), 2, "should have 2 errors")
	ec.AddErrorList(el)
	T.Equal(len(ec.Errors()), 4, "should have 4 errors")

	ec.AddErrorList(nil)
	T.Equal(len(ec.Errors()), 4, "should have ignored nil errorlist")

	ec.AddErrorList([]error{})
	T.Equal(len(ec.Errors()), 4, "should have handled empty list of errors")
	ec.AddErrorList([]error{oneErr, oneErr})
	T.Equal(len(ec.Errors()), 6, "should have handled simple list of errors")
	ec.AddErrorList([]error{oneErr, noErr, oneErr})
	T.Equal(len(ec.Errors()), 8, "should have handled list of errors with nil error in it")

}

func TestNilContainer(t *testing.T) {
	T := testlib.NewT(t)
	defer T.Finish()

	e := tabular.NoSuchCellError{Location: tabular.CellLocation{Row: 3, Column: 2}}

	var ec *tabular.ErrorContainer

	ec.AddError(e)
	T.Equal(ec.Errors(), nil, "no errors from nil container")

	ec = &tabular.ErrorContainer{}
	T.Equal(ec.Errors(), nil, "no errors from badly created container")

	ec.AddError(e)
	T.Equal(len(ec.Errors()), 1, "error stored in badly created container")
	T.Equal(ec.Errors()[0], e, "round-trip error integrity checks out")

	ec2 := &tabular.ErrorContainer{}
	ec2.AddErrorList(ec.Errors())
	T.Equal(len(ec.Errors()), 1, "merged list of errors into missing list")
}
