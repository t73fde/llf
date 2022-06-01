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
	"strings"
	"testing"
)

func TestValidScans(t *testing.T) {
	t.Parallel()
	s := newScanner()
	err := s.checkValid(strings.NewReader(" "))
	if err != nil {
		t.Error(err)
	}
}
