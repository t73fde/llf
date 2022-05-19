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

func TestSymbol(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		val string
		ok  bool
		exp string
	}{
		{"", false, ""},
		{"a", true, "A"},
	}
	smk := sxpf.NewTrivialSymbolMaker()
	for i, tc := range testcases {
		s := smk.MakeSymbol(tc.val)
		if (s != nil) != tc.ok {
			if s == nil {
				t.Errorf("%d: NewSymbol(%q) must not be nil, but is", i, tc.val)
			} else {
				t.Errorf("%d: NewSymbol(%q) must be nil, but is not: %q", i, tc.val, s.GetValue())
			}
			continue
		}
		if s == nil {
			continue
		}
		got := s.GetValue()
		if tc.exp != got {
			t.Errorf("%d: GetValue(%q) != %q, but got %q", i, tc.val, tc.exp, got)
		}
		if !s.Equal(s) {
			t.Errorf("%d: %q is not equal to itself", i, got)
		}

		s2 := smk.MakeSymbol(tc.val)
		if s2 != s {
			t.Errorf("%d: NewSymbol(%q) produces different values if called multiple times", i, tc.val)
		}
	}
}

func FuzzSymbol(f *testing.F) {
	smk := sxpf.NewTrivialSymbolMaker()
	f.Fuzz(func(t *testing.T, in string) {
		t.Parallel()
		s := smk.MakeSymbol(in)
		if !s.Equal(s) {
			if s == nil {
				t.Errorf("nil symbol is not equal to itself")
			} else {
				t.Errorf("%q is not equal to itself", s.GetValue())
			}
		}
	})
}

func TestStringString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		val string
		exp string
	}{
		{"", ""},
		{"a", "a"},
		{"\n", "\\n"},
	}
	for i, tc := range testcases {
		s := sxpf.NewString(tc.val)
		if s == nil {
			t.Errorf("%d: NewString(%q) == nil", i, tc.val)
			continue
		}
		sVal := s.GetValue()
		if sVal != tc.val {
			t.Errorf("%d: NewString(%q) changed value to %q", i, tc.val, sVal)
			continue
		}
		got := s.String()
		if length := len(got); length < 2 {
			t.Errorf("%d: len(String(%q)) < 2: %q (%d)", i, tc.val, got, length)
			continue
		}
		exp := "\"" + tc.exp + "\""
		if got != exp {
			t.Errorf("%d: String(%q) expected %q, but got %q", i, tc.val, exp, got)
		}
	}
}

func TestNewList(t *testing.T) {
	t.Parallel()
	st := sxpf.NewSymbolTable()
	symA, symB, symC, symD := st.MakeSymbol("a"), st.MakeSymbol("b"), st.MakeSymbol("c"), st.MakeSymbol("d")
	testcases := []struct {
		values []sxpf.Value
		exp    string
	}{
		{[]sxpf.Value{}, "()"},
		{[]sxpf.Value{symA}, "(A)"},
		{[]sxpf.Value{symA, symB}, "(A B)"},
		{[]sxpf.Value{symA, symB, symC}, "(A B C)"},
		{[]sxpf.Value{symA, symB, symC, symD}, "[A B C D]"},
	}
	for i, tc := range testcases {
		lst := sxpf.NewSequence(tc.values...)
		got := lst.String()
		if got != tc.exp {
			t.Errorf("%d: expected %q, but got %q", i, tc.exp, got)
		}
	}
}
