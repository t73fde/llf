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
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/t73fde/sxpf"
)

func TestReadString(t *testing.T) {
	pml, pmr := strings.Repeat("(", 5000), strings.Repeat(")", 5000)
	bml, bmr := strings.Repeat("[", 5000), strings.Repeat("]", 5000)
	testcases := []struct {
		src string
		exp string
	}{
		{"a", "A"},
		{"a ", "A"},
		{"a ; comment", "A"},
		{`""`, `""`},
		{`"a"`, `"a"`},
		{`"\""`, `"\""`},
		{`"\\"`, `"\\"`},
		{`"\t"`, `"\t"`},
		{`"\r"`, `"\r"`},
		{`"\n"`, `"\n"`},
		{`"\x"`, `"x"`}, {`"\x4"`, `"x4"`}, {`"\x41"`, `"A"`}, {`"\x4g"`, `"x4g"`},
		{`"\u"`, `"u"`}, {`"\u0"`, `"u0"`}, {`"\u00"`, `"u00"`}, {`"\u004"`, `"u004"`}, {`"\u0042"`, `"B"`},
		{`"\U"`, `"U"`}, {`"\U0"`, `"U0"`}, {`"\U00"`, `"U00"`}, {`"\U000"`, `"U000"`}, {`"\U0000"`, `"U0000"`},
		{`"\U00004"`, `"U00004"`}, {`"\U000043"`, `"C"`},

		{"()", "()"},
		{"(a)", "(A)"},
		{"((a))", "((A))"},
		{"(a b c)", "(A B C)"},
		{"(a b . c)", "(A B . C)"},
		{`("a" b "c")`, `("a" B "c")`},
		{`("a" "b" "c")`, `("a" "b" "c")`},
		{`("a""b""c")`, `("a" "b" "c")`},
		{"(A ((b c) d) (e f))", "(A ((B C) D) (E F))"},
		{"(A.B)", "(A . B)"},
		{`("A"."B")`, `("A" . "B")`},
		{`("A".b)`, `("A" . B)`},

		{"[]", "[]"},
		{"[a]", "[A]"},
		{"[[a]]", "[[A]]"},
		{"[a b c]", "[A B C]"},
		{`["a" b "c"]`, `["a" B "c"]`},
		{"[A [[b c] d] [e f]]", "[A [[B C] D] [E F]]"},

		{"A; bla", "A"},
		{"; bla\na", "A"},
		{"; bla\n\r\n\na", "A"},
		{"; bla\n; bla\na", "A"},
		{"; bla\n\n; bla\na", "A"},

		{pml + pml + pmr + pmr, pml + pml + pmr + pmr},
		{bml + bml + bmr + bmr, bml + bml + bmr + bmr},
		{pml + bml + bmr + pmr, pml + bml + bmr + pmr},
		{bml + pml + pmr + bmr, bml + pml + pmr + bmr},
	}
	for i, tc := range testcases {
		smk := sxpf.NewTrivialSymbolMaker()
		val, err := sxpf.ParseString(smk, tc.src)
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

func TestReadMultiple(t *testing.T) {
	testcases := []struct {
		src string
		exp string
	}{
		{"", ""},
		{"A B", "A B"},
		{" A  B ", "A B"},
		{"[A][B]", "[A] [B]"},
		{" [A] [B] ", "[A] [B]"},
		{"A;B\nC", "A C"},
	}
	for i, tc := range testcases {
		var buf bytes.Buffer
		smk := sxpf.NewTrivialSymbolMaker()
		reader := bytes.NewBufferString(tc.src)
		for cnt := 0; ; cnt++ {
			val, err := sxpf.ParseValue(smk, reader)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("%d: ReadString(%q) resulted in error: %v", i, tc.src, err)
				continue
			}
			if val == nil {
				t.Errorf("%d: ReadString(%q) resulted in nil value", i, tc.src)
				continue
			}
			if cnt > 0 {
				buf.WriteByte((' '))
			}
			buf.WriteString(val.String())
		}
		got := buf.String()
		if tc.exp != got {
			t.Errorf("%d: ReadString(%q) should return %q, but got %q", i, tc.src, tc.exp, got)
		}
	}
}

func TestReadBytesWithError(t *testing.T) {
	testcases := []struct {
		src string
		msg string
	}{
		{"A B", sxpf.ErrMissingEOF.Error()},

		{"(A", sxpf.ErrMissingCloseParenthesis.Error()},
		{"(", sxpf.ErrMissingCloseParenthesis.Error()},
		{")", sxpf.ErrMissingOpenParenthesis.Error()},
		{"())", sxpf.ErrMissingEOF.Error()}, // b/c "()" is already an expression
		{`("`, sxpf.ErrMissingQuote.Error()},
		{`(")`, sxpf.ErrMissingQuote.Error()},
		{`(")()`, sxpf.ErrMissingQuote.Error()},
		{"(.A)", sxpf.ErrMissingCloseParenthesis.Error()},

		{"[A", sxpf.ErrMissingCloseBracket.Error()},
		{"[", sxpf.ErrMissingCloseBracket.Error()},
		{"]", sxpf.ErrMissingOpenBracket.Error()},
		{"[]]", sxpf.ErrMissingEOF.Error()}, // b/c "[]" is already an expression
		{`["`, sxpf.ErrMissingQuote.Error()},
		{`["]`, sxpf.ErrMissingQuote.Error()},
		{`["][]`, sxpf.ErrMissingQuote.Error()},

		{`"`, sxpf.ErrMissingQuote.Error()},
		{`"a`, sxpf.ErrMissingQuote.Error()},
		{`"\`, sxpf.ErrMissingQuote.Error()},
		{`"\"`, sxpf.ErrMissingQuote.Error()},
		{`"\x`, sxpf.ErrMissingQuote.Error()},
		{`"\x1`, sxpf.ErrMissingQuote.Error()},
		{`"\x11`, sxpf.ErrMissingQuote.Error()},
	}
	for i, tc := range testcases {
		smk := sxpf.NewTrivialSymbolMaker()
		val, err := sxpf.ParseBytes(smk, []byte(tc.src))
		if err == nil {
			t.Errorf("%d: ReadString(%q) should result in error, but got value of type %T: %v", i, tc.src, val, val)
			continue
		}
		got := err.Error()
		if got != tc.msg {
			t.Errorf("%d: ReadString(%q) should result in error %q, but got %q", i, tc.src, tc.msg, got)
		}
	}
}

func FuzzReadBytes(f *testing.F) {
	smk := sxpf.NewTrivialSymbolMaker()
	f.Fuzz(func(t *testing.T, src []byte) {
		t.Parallel()
		sxpf.ParseBytes(smk, src)
	})
}
