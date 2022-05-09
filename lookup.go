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

// SymbolMap maps symbols to values.
type SymbolMap struct {
	parent *SymbolMap
	assoc  map[*Symbol]Value
}

func NewSymbolMap(parentMap *SymbolMap) *SymbolMap {
	return &SymbolMap{
		parent: parentMap,
		assoc:  map[*Symbol]Value{},
	}
}

// Set a symbol to its associated value.
func (sm *SymbolMap) Set(sym *Symbol, val Value) {
	sm.assoc[sym] = val
}

// Lookup the value assiated with a given symbol.
func (sm *SymbolMap) Lookup(sym *Symbol) (Value, bool) {
	for curSm := sm; curSm != nil; curSm = curSm.parent {
		if val, found := curSm.assoc[sym]; found {
			return val, true
		}
	}
	return nil, false
}

// LookupForm returns the value associated with the given symbol, if the value
// is a form.
func (sm *SymbolMap) LookupForm(sym *Symbol) (Form, error) {
	if val, found := sm.Lookup(sym); found {
		if form, ok := val.(Form); ok {
			return form, nil
		}
	}
	return nil, ErrNotFormBound(sym)
}

// AsList returns a list representation of the symbol map.
func (sm *SymbolMap) AsList() *List {
	if sm == nil {
		return Nil()
	}
	result := NewList(NewString("symbol"))
	parent := NewList(NewString("parent"))
	if sm.parent == nil {
		parent.Append(Nil())
	} else {
		parent.Append(sm.parent.AsList())
	}
	result.Append(parent)
	for sym, val := range sm.assoc {
		result.Append(NewList(sym, val))
	}
	return result
}

// Sexpr methods

func (sm *SymbolMap) Equal(other Value) bool {
	if sm == nil || other == nil {
		return sm == other
	}
	o, ok := other.(*SymbolMap)
	if !ok {
		return false
	}
	if sm == o {
		return true
	}
	if !sm.parent.Equal(o.parent) || len(sm.assoc) != len(o.assoc) {
		return false
	}
	for sym, val := range sm.assoc {
		if oval, found := o.assoc[sym]; !found || !val.Equal(oval) {
			return false
		}
	}
	return true
}

func (sm *SymbolMap) String() string {
	return sm.AsList().String()
}
