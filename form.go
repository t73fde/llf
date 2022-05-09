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

// Form is a value that can be called. Depending on IsSpecial, the arguments
// are evaluated or not before calling the form.
type Form interface {
	Value

	Call(Environment, []Value) (Value, error)
	IsSpecial() bool
}

// Builtin is a wrapper for a builtin function.
type Builtin struct {
	name    string
	fn      BuiltinFn
	special bool
}

// BuiltinFn is a builtin form that is implemented in Go.
type BuiltinFn func(Environment, []Value) (Value, error)

// NewPrimForm returns a new primitive form.
func NewPrimForm(name string, special bool, f BuiltinFn) *Builtin {
	return &Builtin{name, f, special}
}

func (b *Builtin) Equal(other Value) bool {
	if b == nil || other == nil {
		return Value(b) == other
	}
	if o, ok := other.(*Builtin); ok {
		return b.name == o.name
	}
	return false
}

func (b *Builtin) String() string { return "#" + b.name }

func (b *Builtin) IsSpecial() bool { return b != nil && b.special }
func (b *Builtin) Name() string {
	if b == nil {
		return ""
	}
	return b.name
}

func (b *Builtin) Call(env Environment, args []Value) (Value, error) {
	return b.fn(env, args)
}
