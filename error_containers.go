// Copyright © 2016 Pennock Tech, LLC.
// All rights reserved, except as granted under license.
// Licensed per file LICENSE.txt

package tabular

// ErrorContainer holds a list of errors.
// The tabular package uses a batched error model, where errors accumulate
// in the top-level linked object to which an item is being added.
type ErrorContainer struct {
	errors_ []error
}

// An ErrorReceiver can accumulate errors for later retrieval.
type ErrorReceiver interface {
	AddError(error)
}

// An ErrorSource can be interrogated for errors.
type ErrorSource interface {
	Errors() []error
}

// NewErrorContainer creates an object which satisfies both ErrorReceiver
// and ErrorSource.
func NewErrorContainer() *ErrorContainer {
	return &ErrorContainer{errors_: make([]error, 0, 10)}
}

// AddError puts a Golang error object into the container.  If the error is
// nil, this is a no-op.
func (ec *ErrorContainer) AddError(err error) {
	if ec == nil {
		return
	}
	if err == nil {
		return
	}
	if ec.errors_ == nil {
		ec.errors_ = make([]error, 0, 10)
	}
	ec.errors_ = append(ec.errors_, err)
}

// AddErrorList takes a list of errors and adds them to the container.  Any errors
// which are nil will be dropped.
func (ec *ErrorContainer) AddErrorList(el []error) {
	if ec.errors_ == nil {
		ec.errors_ = el
		return
	}
	for i := range el {
		if el[i] != nil {
			continue
		}
		// drat, long way around
		for j := range el {
			if el[j] != nil {
				ec.errors_ = append(ec.errors_, el[j])
			}
		}
		return
	}
	ec.errors_ = append(ec.errors_, el...)
}

// Errors returns either a non-empty list of errors, or nil.
func (ec *ErrorContainer) Errors() []error {
	// We deliberately handle a nil ErrorContainer, so that other objects delegating to us don't have to worry about it.
	if ec == nil {
		return nil
	}
	if len(ec.errors_) == 0 {
		return nil
	}
	return ec.errors_
}
