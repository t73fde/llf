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
	for i, tc := range testcases {
		s := sxpf.NewSymbol(tc.val)
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

		s2 := sxpf.NewSymbol(tc.val)
		if s2 != s {
			t.Errorf("%d: NewSymbol(%q) produces different values if called multiple times", i, tc.val)
		}
	}
}

func FuzzSymbol(f *testing.F) {
	f.Fuzz(func(t *testing.T, in string) {
		t.Parallel()
		s := sxpf.NewSymbol(in)
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
