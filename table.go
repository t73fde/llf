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

// Table is a mapping between values.
type Table struct {
	m map[string]Value
	k map[string]Atom
}

// Rasa is the value of an empty table ("tabula rasa")
func Rasa() *Table { return &myTabulaRasa }

var myTabulaRasa = Table{make(map[string]Value, 0), make(map[string]Atom, 0)}

// NewTable creates a new table with the given values.
func NewTable(vals ...Value) *Table {
	if len(vals)%2 == 1 {
		vals = append(vals, Nil())
	}
	m := make(map[string]Value, len(vals)/2)
	k := make(map[string]Atom, len(vals)/2)
	t := Table{m, k}
	for i := 0; i < len(vals); i += 2 {
		valueKey, valueValue := vals[i], vals[i+1]
		if valueKey == nil || valueValue == nil {
			return Rasa()
		}
		key, ok := valueKey.(Atom)
		if !ok {
			return Rasa()
		}
		strKey := key.Value()
		m[strKey] = valueValue
		k[strKey] = key
	}
	return &t
}

func (tbl *Table) Equal(other Value) bool {
	if tbl == nil || other == nil {
		return tbl == other
	}
	o, ok := other.(*Table)
	if !ok || len(tbl.m) != len(o.m) {
		return false
	}
	for key, val := range tbl.m {
		if oval, found := o.m[key]; !found || !val.Equal(oval) {
			return false
		}
	}
	return true
}

var (
	lCurly = []byte{'{'}
	rCurly = []byte{'}'}
)

func (tbl *Table) String() string {
	var buf bytes.Buffer
	buf.Write(lCurly)
	first := true
	for key, val := range tbl.m {
		if first {
			first = false
		} else {
			buf.Write(space)
		}
		buf.WriteString(tbl.k[key].String())
		buf.Write(space)
		buf.WriteString(val.String())
	}
	buf.Write(rCurly)
	return buf.String()
}
