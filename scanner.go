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
	"unicode"
)

// TokenType enumerates the concrete type of token.
type TokenType int

// Constants for TokenType.
const (
	TokErr        TokenType = iota // Error
	TokEOF                         // End of Input
	TokLeftParen                   // (
	TokPeriod                      // .
	TokRightParen                  // )
	TokLeftBrack                   // [
	TokRightBrack                  // ]
	TokSymbol                      // a symbol
	TokString                      // "..."
)

// Token is the result of calling a scanner.
type Token struct {
	Typ TokenType
	Val string
}

// Scanner are returning Token from a Reader.
type Scanner struct {
	rd  Reader
	pos uint64 // current bye position in Reader
	err error
}

// NewScanner creates a new scanner.
func NewScanner(rd Reader) *Scanner {
	return &Scanner{rd, 0, nil}
}

func (s *Scanner) Error() error { return s.err }

const (
	chErr rune = -1
	chEOF rune = 0
)

func (s *Scanner) read() rune {
	if s.err != nil {
		return chErr
	}
	ch, width, err := s.rd.ReadRune()
	if err != nil {
		if err == io.EOF {
			return chEOF
		}
		s.err = err
		return chErr
	}
	s.pos += uint64(width)
	return ch
}

func (s *Scanner) Next() Token {
	ch := s.read()
	for unicode.IsSpace(ch) {
		ch = s.read()
	}
	switch ch {
	case chEOF:
		return Token{TokEOF, ""}
	case chErr:
		return Token{TokErr, s.err.Error()}
	case '(':
		return Token{TokLeftParen, "("}
	case '.':
		return Token{TokPeriod, "."}
	case ')':
		return Token{TokRightParen, ")"}
	case '[':
		return Token{TokLeftBrack, "["}
	case ']':
		return Token{TokRightBrack, "]"}
	case '"':
		return s.nextString()
	}
	if unicode.In(ch, unicode.C) {
		// TODO: invalid unicode char at position
		s.err = io.EOF
		return Token{TokErr, s.err.Error()}
	}
	return s.nextSymbol(ch)
}

func (s *Scanner) nextSymbol(ch rune) Token {
	var buf bytes.Buffer
	for {
		buf.WriteRune(ch)
		ch = s.read()
		switch ch {
		case chEOF, '(', '.', ')', '[', ']', '"':
			err := s.rd.UnreadRune()
			if err == nil {
				return Token{TokSymbol, buf.String()}
			}
			s.err = err
			fallthrough
		case chErr:
			return Token{TokErr, s.err.Error()}
		}
		if unicode.IsSpace(ch) {
			// No need to unread, since space will be skipped next time
			return Token{TokSymbol, buf.String()}
		}
		if unicode.In(ch, unicode.C) {
			// TODO: invalid unicode char at position
			s.err = io.EOF
			return Token{TokErr, s.err.Error()}
		}
	}
}

func (s *Scanner) nextString() Token {
	var buf bytes.Buffer
	for {
		ch := s.read()
		switch ch {
		case chEOF:
			s.err = ErrMissingQuote
			fallthrough
		case chErr:
			return Token{TokErr, s.err.Error()}
		case '"':
			return Token{TokString, buf.String()}
		case '\\':
			ch = s.read()
			switch ch {
			case chEOF:
				s.err = ErrMissingQuote
				fallthrough
			case chErr:
				return Token{TokErr, s.err.Error()}
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			case 'n':
				buf.WriteByte('\n')
			case 'x':
				s.parseRune(&buf, ch, 2)
			case 'u':
				s.parseRune(&buf, ch, 4)
			case 'U':
				s.parseRune(&buf, ch, 6)
			default:
				buf.WriteRune(ch)
			}
		default:
			buf.WriteRune(ch)
		}
	}
}

func (s *Scanner) parseRune(buf *bytes.Buffer, curCh rune, numDigits int) {
	var arr [8]rune
	result := rune(0)
	for i := 0; i < numDigits; i++ {
		ch := s.read()
		switch ch {
		case chEOF:
			s.err = ErrMissingQuote
			return
		case chErr:
			return
		case '0':
			arr[i] = '0'
			result <<= 4
			continue
		case '1':
			arr[i] = '1'
			result = (result << 4) + 1
			continue
		case '2':
			arr[i] = '2'
			result = (result << 4) + 2
			continue
		case '3':
			arr[i] = '3'
			result = (result << 4) + 3
			continue
		case '4':
			arr[i] = '4'
			result = (result << 4) + 4
			continue
		case '5':
			arr[i] = '5'
			result = (result << 4) + 5
			continue
		case '6':
			arr[i] = '6'
			result = (result << 4) + 6
			continue
		case '7':
			arr[i] = '7'
			result = (result << 4) + 7
			continue
		case '8':
			arr[i] = '8'
			result = (result << 4) + 8
			continue
		case '9':
			arr[i] = '9'
			result = (result << 4) + 9
			continue
		case 'A', 'a':
			arr[i] = ch
			result = (result << 4) + 10
			continue
		case 'B', 'b':
			arr[i] = ch
			result = (result << 4) + 11
			continue
		case 'C', 'c':
			arr[i] = ch
			result = (result << 4) + 12
			continue
		case 'D', 'd':
			arr[i] = ch
			result = (result << 4) + 13
			continue
		case 'E', 'e':
			arr[i] = ch
			result = (result << 4) + 14
			continue
		case 'F', 'f':
			arr[i] = ch
			result = (result << 4) + 15
			continue
		default:
			xflushRunes(buf, &arr, i, curCh)
			s.err = s.rd.UnreadRune()
			return
		}
	}
	buf.WriteRune(result)
}

func xflushRunes(buf *bytes.Buffer, arr *[8]rune, i int, ch rune) {
	buf.WriteByte('\\')
	buf.WriteRune(ch)
	for j := 0; j < i; j++ {
		buf.WriteRune(arr[j])
	}
}
