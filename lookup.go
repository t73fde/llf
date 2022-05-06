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

// Add a symbol and its associated value.
// If symbol was already associated with a value, this association is overwritten.
func (sm *SymbolMap) Add(sym *Symbol, val Value) {
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
