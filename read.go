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
	"bytes"
	"errors"
	"io"
	"strings"
)

// ErrMissingOpenBracket is raised if there is one additional closing bracket.
var ErrMissingOpenBracket = errors.New("missing opening bracket")

// ErrMissingCloseBracket raised if there is one additonal opening bracket.
var ErrMissingCloseBracket = errors.New("missing closing bracket")

// ErrMissingOpenParenthesis is raised if there is one additional closing parenthesis.
var ErrMissingOpenParenthesis = errors.New("missing opening parenthesis")

// ErrMissingCloseParenthesis raised if there is one additonal opening parenthesis.
var ErrMissingCloseParenthesis = errors.New("missing closing parenthesis")

// ErrMissingQuote is raised if there is no closing quote character.
var ErrMissingQuote = errors.New("missing quote character")

// ErrMissing EOF is raised if there is additional input after an expression.
var ErrMissingEOF = errors.New("missing end of input")

func ReadString(smk SymbolMaker, src string) (Value, error) {
	return consumeReader(smk, strings.NewReader(src))
}

func ReadBytes(smk SymbolMaker, src []byte) (Value, error) {
	return consumeReader(smk, bytes.NewBuffer(src))
}

func consumeReader(smk SymbolMaker, r Reader) (Value, error) {
	val, err := ReadValue(smk, r)
	if err != nil {
		return val, err
	}
	_, err = ReadValue(smk, r)
	if err == io.EOF {
		return val, nil
	}
	return val, ErrMissingEOF
}

type Reader interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

func ReadValue(smk SymbolMaker, r Reader) (Value, error) {
	s := NewScanner(r)
	return parseValue(smk, s, s.Next())
}

func parseValue(smk SymbolMaker, s *Scanner, tok Token) (Value, error) {
	switch tok.Typ {
	case TokEOF:
		return nil, io.EOF
	case TokErr:
		return nil, s.Error()
	case TokLeftParen:
		return parseList(smk, s)
	case TokLeftBrack:
		return parseArray(smk, s)
	case TokString:
		return NewString(tok.Val), nil
	case TokRightParen, TokPeriod:
		return nil, ErrMissingOpenParenthesis
	case TokRightBrack:
		return nil, ErrMissingOpenBracket
	case TokSymbol:
		return smk.MakeSymbol(tok.Val), nil
	default:
		panic(tok)
	}
}

func parseArray(smk SymbolMaker, s *Scanner) (Value, error) {
	elems := []Value{}
	for {
		tok := s.Next()
		switch tok.Typ {
		case TokEOF:
			return nil, ErrMissingCloseBracket
		case TokErr:
			return nil, s.Error()
		case TokRightBrack:
			return NewArray(elems...), nil
		}
		val, err := parseValue(smk, s, tok)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}
}

func parseList(smk SymbolMaker, s *Scanner) (Value, error) {
	elems := []Value{}
loop:
	for {
		tok := s.Next()
		switch tok.Typ {
		case TokEOF:
			return nil, ErrMissingCloseParenthesis
		case TokErr:
			return nil, s.Error()
		case TokRightParen:
			p := Nil()
			for i := len(elems) - 1; i >= 0; i-- {
				p = NewPair(elems[i], p)
			}
			return p, nil
		case TokPeriod:
			if len(elems) == 0 {
				return nil, ErrMissingCloseParenthesis
			}
			break loop
		}
		val, err := parseValue(smk, s, tok)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}

	tok := s.Next()
	switch tok.Typ {
	case TokEOF:
		return nil, ErrMissingCloseParenthesis
	case TokErr:
		return nil, s.Error()
	}
	val, err := parseValue(smk, s, tok)
	if err != nil {
		return nil, err
	}
	tok = s.Next()
	switch tok.Typ {
	case TokErr:
		return nil, s.Error()
	case TokRightParen:
	default:
		return nil, ErrMissingCloseParenthesis
	}
	p := NewPair(elems[len(elems)-1], val)
	for i := len(elems) - 2; i >= 0; i-- {
		p = NewPair(elems[i], p)
	}
	return p, nil
}
