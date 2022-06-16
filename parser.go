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

// ErrMissingOpenCurly is raised if there is one additional closing curly bracket.
var ErrMissingOpenCurly = errors.New("missing opening opening curly brackets")

// ErrMissingCloseCurly raised if there is one additonal opening curly bracket.
var ErrMissingCloseCurly = errors.New("missing closing curly brackets")

// ErrNestedTooDeeply is raised if brackets / parentheses were nested too deeply.
var ErrNestedTooDeeply = errors.New("exceeded max depth")

// ErrMissingQuote is raised if there is no closing quote character.
var ErrMissingQuote = errors.New("missing quote character")

// ErrMissing EOF is raised if there is additional input after an expression.
var ErrMissingEOF = errors.New("missing end of input")

func ParseString(smk SymbolMaker, src string) (Value, error) {
	return consumeReader(smk, strings.NewReader(src))
}

func ParseBytes(smk SymbolMaker, src []byte) (Value, error) {
	return consumeReader(smk, bytes.NewBuffer(src))
}

func consumeReader(smk SymbolMaker, r RuneReader) (Value, error) {
	val, err := ParseValue(smk, r)
	if err != nil {
		return val, err
	}
	_, err = ParseValue(smk, r)
	if err == io.EOF {
		return val, nil
	}
	return val, ErrMissingEOF
}

type RuneReader interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

func ParseValue(smk SymbolMaker, rr RuneReader) (Value, error) {
	pa := NewParser(smk, rr)
	return pa.Parse()
}

type Parser struct {
	smk        SymbolMaker
	sc         *Scanner
	tbuf       []*Token
	maxNesting uint
}

func NewParser(smk SymbolMaker, rr RuneReader) *Parser {
	return &Parser{
		smk:        smk,
		sc:         NewScanner(rr),
		tbuf:       nil,
		maxNesting: 10000,
	}
}

func (pa *Parser) SetMaxNesting(n uint) uint {
	prevN := pa.maxNesting
	pa.maxNesting = n
	return prevN
}

func (pa *Parser) Parse() (Value, error) {
	return pa.parseValue(pa.next())
}

func (pa *Parser) next() Token {
	if tb := pa.tbuf; len(tb) > 0 {
		result := tb[0]
		tb[0] = nil
		if len(tb) > 1 {
			pa.tbuf = tb[1:]
		} else {
			pa.tbuf = nil
		}
		return *result
	}

	tok := pa.sc.Next()
	switch tok.Typ {
	case TokLeftBrack:
		// Fill buffer until right bracket
		return pa.fillBuffer(&tok, TokRightBrack, ErrMissingCloseBracket)
	case TokLeftParen:
		// Fill buffer until right parenthesis
		return pa.fillBuffer(&tok, TokRightParen, ErrMissingCloseParenthesis)
	case TokLeftCurly:
		// Fill buffer until right curly bracket
		return pa.fillBuffer(&tok, TokRightCurly, ErrMissingCloseCurly)
	}
	return tok
}

func (pa *Parser) fillBuffer(token *Token, etyp TokenType, errEOF error) Token {
	typ := token.Typ
	nesting := uint(0)
	for {
		tok := pa.sc.Next()
		switch tok.Typ {
		case TokEOF:
			pa.sc.err = errEOF
			return Token{Typ: TokErr}
		case TokErr:
			return tok
		case TokLeftBrack, TokLeftParen, TokLeftCurly:
			pa.tbuf = append(pa.tbuf, &tok)
			nesting++
			if nesting >= pa.maxNesting {
				pa.sc.err = ErrNestedTooDeeply
				return Token{Typ: TokErr}
			}
		case TokRightBrack, TokRightParen, TokRightCurly:
			pa.tbuf = append(pa.tbuf, &tok)
			if nesting == 0 {
				if tok.Typ == etyp {
					return *token
				}
				switch typ {
				case TokRightBrack:
					pa.sc.err = ErrMissingCloseBracket
				case TokRightParen:
					pa.sc.err = ErrMissingCloseParenthesis
				case TokRightCurly:
					pa.sc.err = ErrMissingCloseCurly
				}
				return Token{Typ: TokErr}
			}
			nesting--
		default:
			pa.tbuf = append(pa.tbuf, &tok)
		}
	}
}
func (pa *Parser) err() error { return pa.sc.Err() }

func (pa *Parser) parseValue(tok Token) (Value, error) {
	switch tok.Typ {
	case TokEOF:
		return nil, io.EOF
	case TokErr:
		return nil, pa.err()
	case TokLeftParen:
		return pa.parseList()
	case TokLeftBrack:
		return pa.parseVector()
	case TokLeftCurly:
		return pa.parseTable()
	case TokString:
		return NewString(tok.Val), nil
	case TokRightParen, TokPeriod:
		return nil, ErrMissingOpenParenthesis
	case TokRightBrack:
		return nil, ErrMissingOpenBracket
	case TokRightCurly:
		return nil, ErrMissingOpenCurly
	case TokSymbol:
		return pa.smk.MakeSymbol(tok.Val), nil
	default:
		panic(tok)
	}
}

func (pa *Parser) parseVector() (Value, error) {
	elems := []Value{}
	for {
		tok := pa.next()
		switch tok.Typ {
		case TokEOF:
			return nil, ErrMissingCloseBracket
		case TokErr:
			return nil, pa.err()
		case TokRightBrack:
			return NewVector(elems...), nil
		}
		val, err := pa.parseValue(tok)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}
}

func (pa *Parser) parseList() (Value, error) {
	elems := []Value{}
loop:
	for {
		tok := pa.next()
		switch tok.Typ {
		case TokEOF:
			return nil, ErrMissingCloseParenthesis
		case TokErr:
			return nil, pa.err()
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
		val, err := pa.parseValue(tok)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}

	tok := pa.next()
	switch tok.Typ {
	case TokEOF:
		return nil, ErrMissingCloseParenthesis
	case TokErr:
		return nil, pa.err()
	}
	val, err := pa.parseValue(tok)
	if err != nil {
		return nil, err
	}
	tok = pa.next()
	switch tok.Typ {
	case TokErr:
		return nil, pa.err()
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

func (pa *Parser) parseTable() (Value, error) {
	elems := []Value{}
	for {
		tok := pa.next()
		switch tok.Typ {
		case TokEOF:
			return nil, ErrMissingCloseCurly
		case TokErr:
			return nil, pa.err()
		case TokRightCurly:
			return NewTable(elems...), nil
		}
		val, err := pa.parseValue(tok)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}
}
