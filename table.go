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
	m map[Value]Value
}

// NewTable creates a new table with the given values.
func NewTable(vals ...Value) *Table {
	if len(vals)%2 == 1 {
		vals = append(vals, Nil())
	}
	m := make(map[Value]Value, len(vals)/2)
	t := Table{m}
	if len(vals) == 0 {
		return &t
	}
	for _, v := range vals {
		if v == nil {
			return &t
		}
	}
	for i := 0; i < len(vals); i += 2 {
		m[vals[i]] = vals[i+1]
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
		buf.WriteString(key.String())
		buf.Write(space)
		buf.WriteString(val.String())
	}
	buf.Write(rCurly)
	return buf.String()
}
