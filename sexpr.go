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
