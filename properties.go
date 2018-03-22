// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular // import "go.pennock.tech/tabular"

import (
	"fmt"
	"reflect"
)

// A PropertyCallback registration is used to update cell properties.
// Properties are akin to stdlib's "context", in that various users of cells
// can maintain their own namespaced properties, using a similar API.
// Examples might include terminal text properties, color, calculation attributes
// and more.  A callback can be registered to be applied to any PropertyOwner,
// whether cell, row, column or table.
type PropertyCallback interface {
	UpdateProperties(PropertyOwner) error
}

// callbackSet is a convenience, to easily embed the set of callbacks inside
// callback-holders.
type callbackSet struct {
	addTime            []PropertyCallback
	preCellRenderTime  []PropertyCallback
	renderTime         []PropertyCallback
	postCellRenderTime []PropertyCallback
}

type callbackTime int

// These constants are used for callback functions, to control when
// the function is invoked.
const (
	// CB_AT_ADD called when an item is added to a container
	CB_AT_ADD callbackTime = iota

	// TODO v2: is `CB_AT_UPDATE` a distinct time, for `c.SetItem()` ?

	// CB_AT_RENDER_PRECELL for callbacks registered on containers, to be
	// called before cell's own callbacks.  Use before dimensions/etc locked
	// down.
	CB_AT_RENDER_PRECELL

	// CB_AT_RENDER is called for two sets: table (not row/column) and cell
	// callbacks.  An assumption is that it might be used for deriving
	// properties such as dimensions, the effect of which cascade outwards.
	CB_AT_RENDER

	// CB_AT_RENDER_POSTCELL is called after the cell's own callbacks.
	CB_AT_RENDER_POSTCELL
)

func invokePropertyCallbacks(
	set callbackSet,
	t callbackTime,
	owner PropertyOwner,
	errTaker ErrorReceiver,
) {
	var cbList []PropertyCallback
	switch t {
	case CB_AT_ADD:
		cbList = set.addTime
	case CB_AT_RENDER:
		cbList = set.renderTime
	case CB_AT_RENDER_PRECELL:
		cbList = set.preCellRenderTime
	case CB_AT_RENDER_POSTCELL:
		cbList = set.postCellRenderTime
	default:
		// internal function, only invoked from within this package, panic is appropriate
		panic("unhandled callbackTime when invoking properties")
	}
	for i := range cbList {
		e := cbList[i].UpdateProperties(owner)
		if e != nil {
			errTaker.AddError(e)
		}
	}
}

type cbTarget int

// These constants are used when registering a callback to indicate what should
// be passed to the callback.  Eg, a row might have a set of callbacks to
// update the row itself, and a separate set of callbacks which are invoked
// upon each cell in the row.
const (
	CB_ON_ITSELF cbTarget = iota
	CB_ON_CELL
	CB_ON_ROW
)

// RegisterPropertyCallback is used to register your object, which has the
// UpdateProperties method available, upon something within the table, to be
// invoked against something (else, perhaps), at a specified time.
func (tb *ATable) RegisterPropertyCallback(
	owner PropertyOwner,
	when callbackTime,
	target cbTarget,
	theNewCallback PropertyCallback,
) error {
	var set *callbackSet

	// TODO: handle _time_ sanity checks, too; warn if would never be invoked.

	switch base := owner.(type) {
	case *ATable:
		switch target {
		case CB_ON_ITSELF:
			set = &base.tableItselfCallbacks
		case CB_ON_CELL:
			set = &base.tableCellCallbacks
		case CB_ON_ROW:
			set = &base.tableRowAdditionCallbacks
		}
	case *column:
		switch target {
		case CB_ON_ITSELF:
			set = &base.columnItselfCallbacks
		case CB_ON_CELL:
			set = &base.cellCallbacks
		}
	case *Row:
		switch target {
		case CB_ON_ITSELF, CB_ON_ROW:
			set = &base.rowItselfCallbacks
		case CB_ON_CELL:
			set = &base.rowCellCallbacks
		}
	case *Cell:
		switch target {
		case CB_ON_ITSELF, CB_ON_CELL:
			set = &base.callbacks
		}
	default:
		return fmt.Errorf("do not know how to register callbacks for type %T", owner)
	}
	if set == nil {
		return fmt.Errorf("unable to register a %v-targetted callback upon a %T", target, owner)
	}

	var cbListPtr *[]PropertyCallback
	switch when {
	case CB_AT_ADD:
		cbListPtr = &set.addTime
	case CB_AT_RENDER:
		cbListPtr = &set.renderTime
	case CB_AT_RENDER_PRECELL:
		cbListPtr = &set.preCellRenderTime
	case CB_AT_RENDER_POSTCELL:
		cbListPtr = &set.postCellRenderTime
	default:
		return fmt.Errorf("unhandled callbackTime when registering properties (%v)", when)
	}

	if cbListPtr == nil {
		*cbListPtr = make([]PropertyCallback, 0, 10)
	}

	*cbListPtr = append(*cbListPtr, theNewCallback)
	return nil
}

// PropertyOwner is the high-level interface satisfied by anything
// which holds metadata in the form of properties.
type PropertyOwner interface {
	SetProperty(interface{}, interface{}) error
	GetProperty(interface{}) interface{}
}

// propertySet is the type of anything which can "be" a property, which
// boils down to an empty property or a value property.
type propertySet interface {
	Value(key interface{}) interface{}
}

// propertyImpl is something which can be embedded in a struct to turn
// it into a PropertyOwner.  Embed as as non-pointer.
type propertyImpl struct {
	properties propertySet
}

func (pi *propertyImpl) GetProperty(key interface{}) interface{} {
	if pi.properties == nil {
		return nil
	}
	return pi.properties.Value(key)
}

func (pi *propertyImpl) SetProperty(key, value interface{}) error {
	if pi == nil {
		return ErrMissingPropertyHolder
	}
	_, remainder := stripReturnValue(pi.properties, key)
	pi.properties = withValue(remainder, key, value)
	return nil
}

type emptyProperty int

func (*emptyProperty) Value(key interface{}) interface{} {
	return nil
}

func (emptyProperty) GoString() string {
	return "NoValue()"
}

var noProperty = new(emptyProperty)

func NoProperty() propertySet {
	return noProperty
}

func withValue(parent propertySet, key, val interface{}) propertySet {
	if key == nil {
		panic("nil property key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueProperty{parent, key, val}
}

// A valueProperty carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded propertySet.
type valueProperty struct {
	propertySet
	key, val interface{}
}

func (c *valueProperty) GoString() string {
	if c.propertySet == nil || c.propertySet == noProperty {
		return fmt.Sprintf("Value(%#v, %#v)", c.key, c.val)
	}
	return fmt.Sprintf("%v.withValue(%#v, %#v)", c.propertySet, c.key, c.val)
}

func (c *valueProperty) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	if c.propertySet == nil || c.propertySet == noProperty {
		return nil
	}
	return c.propertySet.Value(key)
}

func stripReturnValue(ps propertySet, key interface{}) (interface{}, propertySet) {
	top, ok := ps.(*valueProperty)
	if !ok {
		return nil, ps
	}
	if top.propertySet == nil || top.propertySet == noProperty {
		return nil, ps
	}
	if top.key == key {
		return top.val, top.propertySet
	}
	return stripChainReturnValue(top, top, top.propertySet, key)
}

func stripChainReturnValue(top, parent *valueProperty, this_ propertySet, key interface{}) (interface{}, propertySet) {
	this, ok := this_.(*valueProperty)
	if !ok {
		// we break the chain if non-valueProperty are intermingled
		return nil, top
	}
	if this.propertySet == nil || this.propertySet == noProperty {
		return nil, top
	}
	if this.key == key {
		// caller ensures that this != top/parent
		parent.propertySet = this.propertySet
		this.propertySet = nil
		return this.val, top
	}
	return stripChainReturnValue(top, this, this.propertySet, key)
}
