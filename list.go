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

// List is a sequence of values, including sub-lists.
type List struct {
	val    []Value
	frozen bool
}

// Nil is the defined value for an empty list.
func Nil() *List { return &myNIL }

var myNIL = List{val: []Value{}, frozen: true}

// NewList creates a new list with the given values.
func NewList(lstVal ...Value) *List {
	if len(lstVal) == 0 {
		return Nil()
	}
	for _, v := range lstVal {
		if v == nil {
			return Nil()
		}
	}
	return &List{lstVal, false}
}

// Append some more value to a list.
func (lst *List) Append(lstVal ...Value) {
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

// Extend the list by another
func (lst *List) Extend(o *List) {
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

// GetValue returns the list value.
func (lst *List) GetValue() []Value {
	if lst == nil {
		return nil
	}
	return lst.val
}

// Equal retruns true if the other value is equal to this one.
func (lst *List) Equal(other Value) bool {
	if lst == nil || other == nil {
		return lst == other
	}
	o, ok := other.(*List)
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
	space  = []byte{' '}
	lParen = []byte{'('}
	rParen = []byte{')'}
)

func (lst *List) String() string {
	var buf bytes.Buffer
	buf.Write(lParen)
	for i, val := range lst.val {
		if i > 0 {
			buf.Write(space)
		}
		buf.WriteString(val.String())
	}
	buf.Write(rParen)
	return buf.String()
}
