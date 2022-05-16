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
	"unicode"
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
	ch, err := skipSpace(r)
	if err != nil {
		return nil, err
	}
	return parseValue(smk, r, ch)
}

func skipSpace(r Reader) (ch rune, err error) {
	for {
		ch, _, err = r.ReadRune()
		if err != nil {
			return 0, err
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if ch != ';' {
			return ch, nil
		}
		for {
			ch, _, err = r.ReadRune()
			if err != nil {
				return 0, err
			}
			if ch == '\n' || ch == '\r' {
				break
			}
		}
	}
}

func parseValue(smk SymbolMaker, r Reader, ch rune) (Value, error) {
	switch ch {
	case '(':
		return parseList(smk, r)
	case '[':
		return parseArray(smk, r)
	case '"':
		return parseString(r)
	case ')', '.':
		return nil, ErrMissingOpenParenthesis
	case ']':
		return nil, ErrMissingOpenBracket
	default: // Must be symbol
		return parseSymbol(smk, r, ch)
	}
}

func parseSymbol(smk SymbolMaker, r Reader, ch rune) (res Value, err error) {
	var buf bytes.Buffer
	_, err = buf.WriteRune(ch)
	if err != nil {
		return nil, err
	}
	for {
		ch, _, err = r.ReadRune()
		if err == io.EOF {
			return smk.MakeSymbol(buf.String()), nil
		}
		if err != nil {
			return nil, err
		}
		switch ch {
		case '(', ')', '[', ']', '"', ';', '.':
			err = r.UnreadRune()
			return smk.MakeSymbol(buf.String()), err
		}
		if unicode.In(ch, unicode.Space, unicode.C) {
			return smk.MakeSymbol(buf.String()), nil
		}
		_, err = buf.WriteRune(ch)
		if err != nil {
			return nil, err
		}
	}
}

func parseString(r Reader) (Value, error) {
	var buf bytes.Buffer
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, ErrMissingQuote
			}
			return nil, err
		}
		switch ch {
		case '"':
			return NewString(buf.String()), nil
		case '\\':
			ch, _, err = r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, ErrMissingQuote
				}
				return nil, err
			}
			switch ch {
			case 't':
				err = buf.WriteByte('\t')
			case 'r':
				err = buf.WriteByte('\r')
			case 'n':
				err = buf.WriteByte('\n')
			case 'x':
				err = parseRune(r, &buf, ch, 2)
			case 'u':
				err = parseRune(r, &buf, ch, 4)
			case 'U':
				err = parseRune(r, &buf, ch, 6)
			default:
				if unicode.In(ch, unicode.C) {
					return nil, ErrMissingQuote
				}
				_, err = buf.WriteRune(ch)
			}
		default:
			if unicode.In(ch, unicode.C) {
				return nil, ErrMissingQuote
			}
			_, err = buf.WriteRune(ch)
		}
		if err != nil {
			return nil, err
		}
	}
}

var hexMap = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

func parseRune(r Reader, buf *bytes.Buffer, curCh rune, numDigits int) error {
	var arr [8]rune
	arr[0] = curCh
	result := rune(0)
	for i := 0; i < numDigits; i++ {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return ErrMissingQuote
			}
			return err
		}
		if ch == '"' {
			err = flushRunes(buf, &arr, i)
			if err != nil {
				return err
			}
			return r.UnreadRune()
		}
		arr[i+1] = ch
		if hexVal, found := hexMap[ch]; found {
			result = (result << 4) + rune(hexVal)
			continue
		}
		return flushRunes(buf, &arr, i+1)
	}
	_, err := buf.WriteRune(result)
	return err
}

func flushRunes(buf *bytes.Buffer, arr *[8]rune, i int) error {
	for j := 0; j <= i; j++ {
		_, err := buf.WriteRune(arr[j])
		if err != nil {
			return err
		}
	}
	return nil
}

func parseArray(smk SymbolMaker, r Reader) (Value, error) {
	elems := []Value{}
	for {
		ch, err := skipSpace(r)
		if err != nil {
			if err == io.EOF {
				return nil, ErrMissingCloseBracket
			}
			return nil, err
		}
		if ch == ']' {
			return NewArray(elems...), nil
		}
		val, err := parseValue(smk, r, ch)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}
}

func parseList(smk SymbolMaker, r Reader) (Value, error) {
	elems := []Value{}
	for {
		ch, err := skipSpace(r)
		if err != nil {
			if err == io.EOF {
				return nil, ErrMissingCloseParenthesis
			}
			return nil, err
		}
		if ch == ')' {
			p := Nil()
			for i := len(elems) - 1; i >= 0; i-- {
				p = NewPair(elems[i], p)
			}
			return p, nil
		}
		if ch == '.' {
			if len(elems) == 0 {
				return nil, ErrMissingCloseParenthesis
			}
			break
		}
		val, err := parseValue(smk, r, ch)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}

	ch, err := skipSpace(r)
	if err != nil {
		if err == io.EOF {
			return nil, ErrMissingCloseParenthesis
		}
		return nil, err
	}
	val, err := parseValue(smk, r, ch)
	if err != nil {
		return nil, err
	}
	ch, err = skipSpace(r)
	if err != nil {
		if err == io.EOF {
			return nil, ErrMissingCloseParenthesis
		}
		return nil, err
	}
	if ch != ')' {
		return nil, ErrMissingCloseParenthesis
	}
	p := NewPair(elems[len(elems)-1], val)
	for i := len(elems) - 2; i >= 0; i-- {
		p = NewPair(elems[i], p)
	}
	return p, nil
}
