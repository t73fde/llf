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

func TestReadString(t *testing.T) {
	testcases := []struct {
		src string
		exp string
	}{
		{"a", "A"},
		{`""`, `""`},
		{`"a"`, `"a"`},
		{`"\""`, `"\""`},
		{`"\\"`, `"\\"`},
		{`"\t"`, `"\t"`},
		{`"\r"`, `"\r"`},
		{`"\n"`, `"\n"`},
		{`"\x"`, `"x"`}, {`"\x4"`, `"x4"`}, {`"\x41"`, `"A"`},
		{`"\u"`, `"u"`}, {`"\u0"`, `"u0"`}, {`"\u00"`, `"u00"`}, {`"\u004"`, `"u004"`}, {`"\u0042"`, `"B"`},
		{`"\U"`, `"U"`}, {`"\U0"`, `"U0"`}, {`"\U00"`, `"U00"`}, {`"\U000"`, `"U000"`}, {`"\U0000"`, `"U0000"`},
		{`"\U00004"`, `"U00004"`}, {`"\U000043"`, `"C"`},
		{"()", "()"},
		{"(a)", "(A)"},
		{"((a))", "((A))"},
		{"(a b c)", "(A B C)"},
		{`("a" b "c")`, `("a" B "c")`},
		{"(A ((b c) d) (e f))", "(A ((B C) D) (E F))"},
	}
	env := sxpf.NewTrivialEnvironment()
	for i, tc := range testcases {
		val, err := sxpf.ReadString(env, tc.src)
		if err != nil {
			t.Errorf("%d: ReadString(%q) resulted in error: %v", i, tc.src, err)
			continue
		}
		if val == nil {
			t.Errorf("%d: ReadString(%q) resulted in nil value", i, tc.src)
			continue
		}
		got := val.String()
		if tc.exp != got {
			t.Errorf("%d: ReadString(%q) should return %q, but got %q", i, tc.src, tc.exp, got)
		}
	}
}

func FuzzReadBytes(f *testing.F) {
	env := sxpf.NewTrivialEnvironment()
	f.Fuzz(func(t *testing.T, src []byte) {
		t.Parallel()
		sxpf.ReadBytes(env, src)
	})
}
