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

import "io"

// Form is a wrapper for a primitive or a user defined form / function.
// Currently, only primitive functions are allowed.
type Form struct {
	name      string
	primitive PrimForm
	special   bool
}

// PrimForm is a primitve form that is implemented in Go.
type PrimForm func(Environment, []Value) (Value, error)

// NewPrimForm returns a new primitive form.
func NewPrimForm(name string, special bool, f PrimForm) *Form {
	return &Form{name, f, special}
}

func (f *Form) Equal(other Value) bool {
	if f == nil || other == nil {
		return Value(f) == other
	}
	if o, ok := other.(*Form); ok {
		return f.name == o.name
	}
	return false
}

func (f *Form) Encode(w io.Writer) (int, error) { return io.WriteString(w, f.String()) }

func (f *Form) String() string { return "#" + f.name }

func (f *Form) IsSpecial() bool { return f != nil && f.special }
func (f *Form) Name() string {
	if f == nil {
		return ""
	}
	return f.name
}

func (f *Form) Call(env Environment, args []Value) (Value, error) {
	return f.primitive(env, args)
}
