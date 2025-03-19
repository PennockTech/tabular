// Copyright © 2018,2025 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package properties // import "go.pennock.tech/tabular/properties"

import (
	"fmt"
)

type propertyKey struct {
	name string
}

func (p *propertyKey) String() string { return "miscellaneous property keyid " + p.name }

var (
	Skipable = &propertyKey{"skipable"}
	Omit     = &propertyKey{"omit"}
)

type ErrPropertyNotBool struct {
	Property *propertyKey
}

func (e ErrPropertyNotBool) Error() string {
	return "tabular: property '" + e.Property.name + "' not bool value"
}

type ErrPropertyNotBoolForPosition struct {
	Label         string
	Property      *propertyKey
	PositionLabel string
	Position      any
	Item          any
}

func (e ErrPropertyNotBoolForPosition) Error() string {
	return fmt.Sprintf("tabular: %s: %s %v property %s not bool value but %T", e.Label, e.PositionLabel, e.Position, e.Property.name, e.Item)
}

type ErrPropertyNotSet struct {
	Property *propertyKey
}

func (e ErrPropertyNotSet) Error() string {
	return "tabular: property '" + e.Property.name + "' not set"
}

// ExpectBoolProperty converts a property value to a bool, or returns an error.
func ExpectBoolProperty(p *propertyKey, v any, label string, positionLabel string, position any) (bool, error) {
	if v == nil {
		return false, ErrPropertyNotSet{p}
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, ErrPropertyNotBoolForPosition{Label: label, Property: p, PositionLabel: positionLabel, Position: position, Item: v}
}

// ExpectBoolPropertyOrNil converts a property value to a bool, or returns an error.
// An unset (nil) property is treated as false
func ExpectBoolPropertyOrNil(p *propertyKey, v any, label string, positionLabel string, position any) (bool, error) {
	if v == nil {
		return false, nil
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, ErrPropertyNotBoolForPosition{Label: label, Property: p, PositionLabel: positionLabel, Position: position, Item: v}
}
