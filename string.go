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

import "bytes"

// String is a string value without any restrictions.
type String struct {
	val string
}

// NewString creates a new string with the given value.
func NewString(strVal string) *String { return &String{strVal} }

// GetValue returns the string value.
func (str *String) GetValue() string { return str.val }

// Equal retruns true if the other value is equal to this one.
func (str *String) Equal(other Value) bool {
	if str == nil || other == nil {
		return str == other
	}
	if o, ok := other.(*String); ok {
		return str.val == o.val
	}
	return false
}

var (
	quote        = []byte{'"'}
	encBackslash = []byte{'\\', '\\'}
	encQuote     = []byte{'\\', '"'}
	encNewline   = []byte{'\\', 'n'}
	encTab       = []byte{'\\', 't'}
	encCr        = []byte{'\\', 'r'}
	encUnicode   = []byte{'\\', 'x', '0', '0'}
	encHex       = []byte("0123456789ABCDEF")
)

func (str *String) String() string {
	var buf bytes.Buffer
	buf.Write(quote)
	last := 0
	for i, ch := range str.val {
		var b []byte
		switch ch {
		case '\t':
			b = encTab
		case '\r':
			b = encCr
		case '\n':
			b = encNewline
		case '"':
			b = encQuote
		case '\\':
			b = encBackslash
		default:
			if ch >= ' ' {
				continue
			}
			b = encUnicode
			b[2] = encHex[ch>>4]
			b[3] = encHex[ch&0xF]
		}
		buf.WriteString(str.val[last:i])
		buf.Write(b)
		last = i + 1
	}
	buf.WriteString(str.val[last:])
	buf.Write(quote)
	return buf.String()
}
