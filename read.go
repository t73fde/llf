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
	"io"
	"strings"
	"unicode"
)

func ReadString(env Environment, src string) (Value, error) {
	return ReadValue(env, strings.NewReader(src))
}

func ReadBytes(env Environment, src []byte) (Value, error) {
	return ReadValue(env, bytes.NewBuffer(src))
}

type Reader interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

func ReadValue(env Environment, r Reader) (Value, error) {
	ch, err := skipSpace(r)
	if err != nil {
		return nil, err
	}
	return parseValue(env, r, ch)
}

func skipSpace(r Reader) (rune, error) {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return 0, err
		}
		if unicode.IsSpace(ch) {
			continue
		}
		return ch, nil
	}
}

func parseValue(env Environment, r Reader, ch rune) (Value, error) {
	switch ch {
	case '(':
		return parseList(env, r)
	case '"':
		return parseString(r)
	default: // Must be symbol
		return parseSymbol(env, r, ch)
	}
}

func parseSymbol(env Environment, r Reader, ch rune) (res Value, err error) {
	var buf bytes.Buffer
	buf.WriteRune(ch)
	for {
		ch, _, err = r.ReadRune()
		if err == io.EOF {
			return env.MakeSymbol(buf.String()), nil
		}
		if err != nil {
			return nil, err
		}
		switch ch {
		case ')':
			err = r.UnreadRune()
			fallthrough
		case '(', '"':
			return env.MakeSymbol(buf.String()), err
		}
		if unicode.In(ch, unicode.Space, unicode.C) {
			return env.MakeSymbol(buf.String()), nil
		}
		buf.WriteRune(ch)
	}
}

func parseString(r Reader) (Value, error) {
	var buf bytes.Buffer
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}
		switch ch {
		case '"':
			return NewString(buf.String()), nil
		case '\\':
			ch, _, err = r.ReadRune()
			if err != nil {
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
				_, err = buf.WriteRune(ch)
			}
			if err != nil {
				return nil, err
			}
		default:
			buf.WriteRune(ch)
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

func parseList(env Environment, r Reader) (Value, error) {
	elems := []Value{}
	for {
		ch, err := skipSpace(r)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		if ch == ')' {
			return NewList(elems...), nil
		}
		val, err := parseValue(env, r, ch)
		if err != nil {
			return nil, err
		}
		elems = append(elems, val)
	}
}
