//-----------------------------------------------------------------------------
// Copyright (c) 2022 Detlef Stern
//
// This file is part of sxpf.
//
// sxpf is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//-----------------------------------------------------------------------------

// Package sxpf allows to work with symbolic expressions, s-expressions.
package sxpf

import "fmt"

// Value is a generic value, the set of all possible values of a s-expression.
type Value interface {
	Equal(Value) bool
	String() string
}

// Sequence is a generic value that is a sequence of values.
type Sequence interface {
	Value
	GetSlice() []Value
}

// NewSequence stores the slice of Values in a pair list or in an array, depending
// on its length.
func NewSequence(values ...Value) Sequence {
	// A pair list has an overhead of four words per list element (both first
	// and second value is a reference to an interface, two words each).
	// An array has a constant overhead of 4 words (the Array struct) +
	// 3 words (for the slice) and a linear overhead of two word per
	// list element (for the interface reference to the value).
	// Let x be the number of list elements. Then we need a x where
	// 4*x > 2*x + 7 <=> x+x+x+x > x+x+7 <=> x+x > 7 <=> x > 3.

	if len(values) > 3 {
		return NewArray(values...)
	}
	return NewPairFromSlice(values)
}

// GetSymbol returns the idx value of args as a Symbol.
func GetSymbol(args []Value, idx int) (*Symbol, error) {
	if idx < 0 || len(args) <= idx {
		return nil, makeErrIndexOutOfBounds(args, idx)
	}
	if val, ok := args[idx].(*Symbol); ok {
		return val, nil
	}
	return nil, fmt.Errorf("%v / %d is not a symbol", args[idx], idx)
}

// GetString returns the idx value of args as a String.
func GetString(args []Value, idx int) (string, error) {
	if idx < 0 || len(args) <= idx {
		return "", makeErrIndexOutOfBounds(args, idx)
	}
	if val, ok := args[idx].(*String); ok {
		return val.GetValue(), nil
	}
	if val, ok := args[idx].(*Symbol); ok {
		return val.GetValue(), nil
	}
	return "", fmt.Errorf("%v / %d is not a string", args[idx], idx)
}

// GetArray returns the idx value of args as an array.
func GetArray(args []Value, idx int) (*Array, error) {
	if idx < 0 || len(args) <= idx {
		return nil, makeErrIndexOutOfBounds(args, idx)
	}
	if val, ok := args[idx].(*Array); ok {
		return val, nil
	}
	return nil, fmt.Errorf("%v / %d is not an array", args[idx], idx)
}

func makeErrIndexOutOfBounds(args []Value, idx int) error {
	return fmt.Errorf("index %d out of bounds: %v", idx, args)
}
