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

// Pair is a type with two values. In other lisps it is often called "cons",
// "cons-cell", or "cell".
type Pair struct {
	first, second Value
}

// NewPair creates a new pair from two values.
func NewPair(first, second Value) *Pair {
	return &Pair{first, second}
}

// GetFirst returns the first value of a pair
func (p *Pair) GetFirst() Value {
	if p != nil {
		return p.first
	}
	return nil
}

// GetSecond returns the second value of a pair
func (p *Pair) GetSecond() Value {
	if p != nil {
		return p.second
	}
	return nil
}

// GetSlice returns the pair list elements as a slice of Values.
func (p *Pair) GetSlice() []Value {
	if p == nil {
		return nil
	}
	var result []Value
	cp := p
	for {
		result = append(result, cp.first)
		second := cp.second
		np, ok := second.(*Pair)
		if !ok {
			result = append(result, second)
			return result
		}
		if np == nil {
			return result
		}
		cp = np
	}
}

// Nil() returns the empty pair.
func Nil() *Pair { return nilPair }

var nilPair *Pair

func (p *Pair) Equal(other Value) bool {
	if p == nil || other == nil {
		return p == other
	}
	if o, ok := other.(*Pair); ok {
		return p.first.Equal(o.first) && p.second.Equal(o.second)
	}
	return false
}

func (p *Pair) String() string {
	if p == nil {
		return "()"
	}
	var buf bytes.Buffer
	buf.WriteByte('(')
	for cp := p; ; {
		if cp != p {
			buf.WriteByte(' ')
		}
		if first := cp.first; first != nil {
			if s := first.String(); s != "" {
				buf.WriteString(s)
			} else {
				buf.WriteString("()")
			}
		} else {
			buf.WriteString("()")
		}
		sval := cp.second
		if sval == nil {
			break
		}
		if np, ok := sval.(*Pair); ok {
			if np == nil {
				break
			}
			cp = np
			continue
		}
		if s := sval.String(); s != "" {
			buf.WriteString(" . ")
			buf.WriteString(s)
		}
		break
	}
	buf.WriteByte(')')
	return buf.String()
}
