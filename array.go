//-----------------------------------------------------------------------------
// Copyright (c) 2022 Detlef Stern
//
// This file is part of sxpf.
//
// sxpf is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//-----------------------------------------------------------------------------

package sxpf

import "bytes"

// Array is a sequence of values, including sub-arrays.
type Array struct {
	val    []Value
	frozen bool
}

// Empty is the defined value for an empty array.
func Empty() *Array { return &myNIL }

var myNIL = Array{val: []Value{}, frozen: true}

// NewArray creates a new array with the given values.
func NewArray(lstVal ...Value) *Array {
	if len(lstVal) == 0 {
		return Empty()
	}
	for _, v := range lstVal {
		if v == nil {
			return Empty()
		}
	}
	return &Array{lstVal, false}
}

// Append some more value to an array.
func (lst *Array) Append(lstVal ...Value) {
	if lst.frozen {
		return
	}
	for _, v := range lstVal {
		if v == nil {
			return
		}
	}
	lst.val = append(lst.val, lstVal...)
}

// Extend the array by another
func (lst *Array) Extend(o *Array) {
	if lst.frozen {
		return
	}
	if o != nil {
		for _, v := range o.val {
			if v == nil {
				return
			}
		}
		lst.val = append(lst.val, o.val...)
	}
}

// GetValue returns the array value as a slice of Values.
func (lst *Array) GetValue() []Value {
	if lst == nil {
		return nil
	}
	return lst.val
}

func (lst *Array) GetSlice() []Value { return lst.GetValue() }

// Equal retruns true if the other value is equal to this one.
func (lst *Array) Equal(other Value) bool {
	if lst == nil || other == nil {
		return lst == other
	}
	o, ok := other.(*Array)
	if !ok || len(lst.val) != len(o.val) {
		return false
	}
	for i, val := range lst.val {
		if !val.Equal(o.val[i]) {
			return false
		}
	}
	return true
}

var (
	space    = []byte{' '}
	lBracket = []byte{'['}
	rBracket = []byte{']'}
)

func (lst *Array) String() string {
	var buf bytes.Buffer
	buf.Write(lBracket)
	for i, val := range lst.val {
		if i > 0 {
			buf.Write(space)
		}
		buf.WriteString(val.String())
	}
	buf.Write(rBracket)
	return buf.String()
}
