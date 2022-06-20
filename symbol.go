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

import "strings"

// Symbol is a value that identifies something.
type Symbol struct {
	val string
}

// GetValue returns the string value of the symbol.
func (sym *Symbol) GetValue() string { return sym.val }

// Equal retruns true if the other value is equal to this one.
func (sym *Symbol) Equal(other Value) bool {
	if sym == nil || other == nil {
		return sym == other
	}
	if o, ok := other.(*Symbol); ok {
		return strings.EqualFold(sym.val, o.val)
	}
	return false
}

func (sym *Symbol) String() string { return sym.val }
func (sym *Symbol) Value() string  { return sym.val }

// SymbolTable allows to create unique symbols.
type SymbolTable struct {
	m map[string]*Symbol
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{map[string]*Symbol{}}
}

func (st *SymbolTable) MakeSymbol(s string) *Symbol {
	if s == "" {
		return nil
	}
	sym, found := st.m[s]
	if !found {
		sym = &Symbol{strings.ToUpper(s)}
		st.m[s] = sym
	}
	return sym
}
