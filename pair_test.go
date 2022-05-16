//-----------------------------------------------------------------------------
// Copyright (c) 2022 Detlef Stern
//
// This file is part of sxpf.
//
// sxpf is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//-----------------------------------------------------------------------------

package sxpf_test

import (
	"testing"

	"github.com/t73fde/sxpf"
)

func TestPairString(t *testing.T) {
	n := sxpf.Nil()
	st := sxpf.NewSymbolTable()
	a, b, c := st.MakeSymbol("A"), st.MakeSymbol("B"), st.MakeSymbol("C")
	np := sxpf.NewPair
	testcases := []struct {
		pair *sxpf.Pair
		exp  string
	}{
		{nil, "()"}, {n, "()"},
		{np(nil, nil), "(())"}, {np(n, nil), "(())"}, {np(nil, n), "(())"}, {np(n, n), "(())"},
		{np(a, nil), "(A)"}, {np(a, n), "(A)"},
		{np(nil, a), "(() . A)"}, {np(n, a), "(() . A)"},
		{np(a, np(b, nil)), "(A B)"}, {np(a, np(b, n)), "(A B)"},
		{np(a, np(b, c)), "(A B . C)"},
		{np(np(a, b), c), "((A . B) . C)"},
		{np(np(a, b), np(c, nil)), "((A . B) C)"}, {np(np(a, b), np(c, n)), "((A . B) C)"},
	}
	for i, tc := range testcases {
		got := tc.pair.String()
		if tc.exp != got {
			t.Errorf("%d: expected %q, but got %q", i, tc.exp, got)
		}
	}
}

func TestPairEqual(t *testing.T) {
	if !sxpf.Nil().Equal(sxpf.Nil()) {
		t.Error("Nil() is not equal to itself")
	}
	s1 := sxpf.NewString("a")
	p1 := sxpf.NewPair(s1, sxpf.Nil())
	if !p1.Equal(p1) {
		t.Errorf("%v is not equal to itself", p1)
	}
	p2 := sxpf.NewPair(s1, sxpf.Nil())
	if !p1.Equal(p2) {
		t.Errorf("%v is not equal to %v", p1, p2)
	}
	if !p2.Equal(p1) {
		t.Errorf("%v is not equal to %v", p2, p1)
	}
	p3 := sxpf.NewPair(p1, p2)
	if !p3.Equal(p3) {
		t.Errorf("%v is not equal to itself", p3)
	}
	if p3.Equal(s1) {
		t.Errorf("%v is equal to %v", p3, s1)
	}
}
