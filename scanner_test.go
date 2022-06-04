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
	"strings"
	"testing"

	"github.com/t73fde/sxpf"
)

func TestValidScans(t *testing.T) {
	t.Parallel()
	s := sxpf.NewScanner(strings.NewReader(" "))
	tok := s.Next()
	if tok.Typ != sxpf.TokEOF {
		t.Error(tok)
	}
}

func TestScanner(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		src string
		exp string
	}{
		{"", ""},
		{"a", "a"},
		{"a ; comment", "a"},
		{"(", "("},
		{"(.)[]", "(.)[]"},
		{"a.", "a."},
		{`""`, ``},
		{`"a"`, `a`},
		{`"\""`, `"`},
		{`"\\"`, "\\"},
		{`"\t"`, "\t"},
		{`"\r"`, "\r"},
		{`"\n"`, "\n"},
		{`"\x"`, `\x`}, {`"\x4"`, `\x4`}, {`"\x41"`, `A`}, {`"\x4g"`, `\x4g`},
		{`"\u"`, `\u`}, {`"\u0"`, `\u0`}, {`"\u00"`, `\u00`}, {`"\u004"`, `\u004`}, {`"\u0042"`, `B`},
		{`"\U"`, `\U`}, {`"\U0"`, `\U0`}, {`"\U00"`, `\U00`}, {`"\U000"`, `\U000`}, {`"\U0000"`, `\U0000`},
		{`"\U00004"`, `\U00004`}, {`"\U000043"`, `C`},
	}
	for i, tc := range testcases {
		var buf bytes.Buffer
		s := sxpf.NewScanner(strings.NewReader(tc.src))
		for {
			tok := s.Next()
			if tok.Typ == sxpf.TokEOF {
				break
			}
			if tok.Typ == sxpf.TokErr {
				buf.WriteByte('E')
				break
			}
			buf.WriteString(tok.Val)
		}
		got := buf.String()
		if got != tc.exp {
			t.Errorf("%d: %q -> %q, but got %q", i, tc.src, tc.exp, got)
		}
	}
}
