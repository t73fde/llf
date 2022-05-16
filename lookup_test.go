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

func TestSymbolMapSexpr(t *testing.T) {
	smk := sxpf.NewTrivialSymbolMaker()
	m := sxpf.NewSymbolMap(nil)
	sym := smk.MakeSymbol("map")
	m.Set(sym, m) // A SymbolMap is itself a Sexpr
	m1, found := m.Lookup(sym)
	if !found {
		t.Errorf("Symbol %v not found, but should be there", sym)
	} else if !m.Equal(m1) {
		t.Errorf("Expected map %v, but got: %v", m, m1)
	}
}

func TestSymbolMapPrint(t *testing.T) {
	smk := sxpf.NewTrivialSymbolMaker()
	sm1 := sxpf.NewSymbolMap(nil)
	sm1.Set(smk.MakeSymbol("sym1"), sxpf.NewString("val1"))
	got := sm1.String()
	exp := `["symbol" ["parent" []] [SYM1 "val1"]]`
	if exp != got {
		t.Errorf("sm1:\nexpected: %v,\nbut got: %v", exp, got)
	}
	sm2 := sxpf.NewSymbolMap(sm1)
	sm2.Set(smk.MakeSymbol("sym2"), sxpf.NewString("val2"))
	got = sm2.String()
	exp = `["symbol" ["parent" ["symbol" ["parent" []] [SYM1 "val1"]]] [SYM2 "val2"]]`
	if exp != got {
		t.Errorf("sm2:\nexpected: %v,\n but got: %v", exp, got)
	}
}
