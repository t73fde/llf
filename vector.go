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

// Vector is a sequence of values, including sub-vectors.
type Vector struct {
	val    []Value
	frozen bool
}

// Empty is the defined value for an empty vector.
func Empty() *Vector { return &myNIL }

var myNIL = Vector{val: []Value{}, frozen: true}

// NewVector creates a new vector with the given values.
func NewVector(vals ...Value) *Vector {
	if len(vals) == 0 {
		return Empty()
	}
	for _, v := range vals {
		if v == nil {
			return Empty()
		}
	}
	return &Vector{vals, false}
}

// Append some more value to a vector.
func (v *Vector) Append(lstVal ...Value) {
	if v.frozen {
		return
	}
	for _, v := range lstVal {
		if v == nil {
			return
		}
	}
	v.val = append(v.val, lstVal...)
}

// Extend the vector by another
func (v *Vector) Extend(o *Vector) {
	if v.frozen {
		return
	}
	if o != nil {
		for _, v := range o.val {
			if v == nil {
				return
			}
		}
		v.val = append(v.val, o.val...)
	}
}

func (v *Vector) GetSlice() []Value {
	if v == nil {
		return nil
	}
	return v.val
}

func (v *Vector) Equal(other Value) bool {
	if v == nil || other == nil {
		return v == other
	}
	o, ok := other.(*Vector)
	if !ok || len(v.val) != len(o.val) {
		return false
	}
	for i, val := range v.val {
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

func (v *Vector) String() string {
	var buf bytes.Buffer
	buf.Write(lBracket)
	for i, val := range v.val {
		if i > 0 {
			buf.Write(space)
		}
		buf.WriteString(val.String())
	}
	buf.Write(rBracket)
	return buf.String()
}
