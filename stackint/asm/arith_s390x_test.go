// Copyright (c) 2009 The Go Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// +build s390x,!math_big_pure_go

package asm

import (
	"testing"
)

// Tests whether the non vector routines are working, even when the tests are run on a
// vector-capable machine

func TestFunVVnovec(t *testing.T) {
	if hasVX == true {
		for _, a := range sumVV {
			arg := a
			testFunVV(t, "AddVV_novec", AddVV_novec, arg)

			arg = argVV{a.z, a.y, a.x, a.c}
			testFunVV(t, "AddVV_novec symmetric", AddVV_novec, arg)

			arg = argVV{a.x, a.z, a.y, a.c}
			testFunVV(t, "SubVV_novec", SubVV_novec, arg)

			arg = argVV{a.y, a.z, a.x, a.c}
			testFunVV(t, "SubVV_novec symmetric", SubVV_novec, arg)
		}
	}
}

func TestFunVWnovec(t *testing.T) {
	if hasVX == true {
		for _, a := range sumVW {
			arg := a
			testFunVW(t, "AddVW_novec", AddVW_novec, arg)

			arg = argVW{a.x, a.z, a.y, a.c}
			testFunVW(t, "SubVW_novec", SubVW_novec, arg)
		}
	}
}
