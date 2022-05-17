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
	"testing"

	"github.com/t73fde/sxpf"
)

func TestEvaluate(t *testing.T) {
	testcases := []struct {
		src string
		exp string
	}{
		{"a", "A"},
		{`"a"`, `"a"`},
		{"(CAT a b)", `"AB"`},
		{"(QUOTE [(A b) c])", "[(A B) C]"},
		{"[CAT a b]", `"AB"`},
		{"[QUOTE [[A b] c]]", "[[A B] C]"},
	}
	env := newTestEnv()
	for i, tc := range testcases {
		expr, err := sxpf.ReadString(env, tc.src)
		if err != nil {
			t.Error(err)
			continue
		}
		val, err := sxpf.Evaluate(env, expr)
		if err != nil {
			t.Error(err)
			continue
		}
		got := val.String()
		if got != tc.exp {
			t.Errorf("%d: %v should evaluate to %v, but got: %v", i, tc.src, tc.exp, got)
		}
	}
}

type testEnv struct {
	symbols sxpf.SymbolTable
	symMap  *sxpf.SymbolMap
}

func newTestEnv() *testEnv {
	env := testEnv{symbols: sxpf.NewSymbolTable()}
	symMap := sxpf.NewSymbolMap(nil)
	for _, form := range testForms {
		symMap.Set(env.MakeSymbol(form.Name()), form)
	}
	env.symMap = symMap
	return &env
}

var testForms = []*sxpf.Builtin{
	sxpf.NewBuiltin(
		"CAT",
		false, 0, -1,
		func(env sxpf.Environment, args []sxpf.Value) (sxpf.Value, error) {
			var buf bytes.Buffer
			for _, arg := range args {
				buf.WriteString(arg.String())
			}
			return sxpf.NewString(buf.String()), nil
		},
	),
	sxpf.NewBuiltin(
		"QUOTE",
		true, 1, 1,
		func(env sxpf.Environment, args []sxpf.Value) (sxpf.Value, error) {
			return args[0], nil
		},
	),
}

func (te *testEnv) MakeSymbol(s string) *sxpf.Symbol                 { return te.symbols.MakeSymbol(s) }
func (te *testEnv) LookupForm(sym *sxpf.Symbol) (sxpf.Form, error)   { return te.symMap.LookupForm(sym) }
func (*testEnv) EvaluateSymbol(sym *sxpf.Symbol) (sxpf.Value, error) { return sym, nil }
func (*testEnv) EvaluateString(str *sxpf.String) (sxpf.Value, error) { return str, nil }
func (e *testEnv) EvaluateList(p *sxpf.Pair) (sxpf.Value, error)     { return e.evalAsCall(p.GetSlice()) }
func (e *testEnv) EvaluateArray(lst *sxpf.Array) (sxpf.Value, error) {
	return e.evalAsCall(lst.GetValue())
}

func (e *testEnv) evalAsCall(vals []sxpf.Value) (sxpf.Value, error) {
	res, err, done := sxpf.EvaluateCall(e, vals)
	if done {
		return res, err
	}
	result, err := sxpf.EvaluateSlice(e, vals)
	if err != nil {
		return nil, err
	}
	return sxpf.NewArray(result...), nil
}
