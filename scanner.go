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

import (
	"fmt"
	"io"
)

// S-Expression scanner state machine.
// It is loosely modeled after the JSON value parser state machine of the Go
// standard library.
//
// To be more exact, "scanner" is both a lexical analyzer *and* some kind of
// parser. An ordinary scanner is effectively equivalent to a final state machine,
// which does not allow recursion. Since an s-expression is defined recursively,
// the "scanner" defined here maintains a stack where it tracks nested lists and
// similar s-expression elements.

// A SyntaxError is a description of a s-expression syntax error.
type SyntaxError struct {
	text string
}

func (e *SyntaxError) Error() string { return e.text }

// stepResult is the return value of the transition function of a scanner ("step").
type stepResult int

// Constant values for stepResult.
const (
	scanContinue stepResult = iota // Uninteresting rune

	// Values that signals an end state.
	scanEnd
	scanError // Some error detected, error value is in scanner.err
)

// parseState is an indicator about the current state when the scanner enters a
// recursive element, like a left parenthesis, which indicated a nested list.
type parseState int

// Constand values for parseState.
const (
	parseListValue parseState = iota
)

// A scanner is the state machine to scan s-expressions.
type scanner struct {
	err        error
	step       func(rune) stepResult
	parseStack []parseState
}

func newScanner() *scanner {
	s := &scanner{}
	s.reset()
	return s
}

// reset prepared the scanner to be used from a safe / initial state.
func (s *scanner) reset() {
	s.err = nil
	s.step = s.stepError
	s.parseStack = s.parseStack[0:0]
}

func (s *scanner) checkEOI() stepResult {
	if s.err != nil {
		return scanError
	}
	return scanEnd
}

func (s *scanner) checkValid(r Reader) error {
	s.reset()
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if s.step(ch) == scanError {
			return s.err
		}
	}
	if s.checkEOI() == scanError {
		return s.err
	}
	return nil
}

// stepError is the error state transition function.
func (*scanner) stepError(rune) stepResult { return scanError }

func (s *scanner) error(ch rune, msg string) stepResult {
	rs := string(ch)
	s.step = s.stepError
	s.err = &SyntaxError{fmt.Sprintf("invalid character %q %s", rs, msg)}
	return scanError
}
