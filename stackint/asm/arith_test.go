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

package asm

import (
	"fmt"
	// "internal/testenv"
	"math/rand"
	"testing"
)

var isRaceBuilder = false

type nat []Word

type funWW func(x, y, c Word) (z1, z0 Word)
type argWW struct {
	x, y, c, z1, z0 Word
}

var sumWW = []argWW{
	{0, 0, 0, 0, 0},
	{0, 1, 0, 0, 1},
	{0, 0, 1, 0, 1},
	{0, 1, 1, 0, 2},
	{12345, 67890, 0, 0, 80235},
	{12345, 67890, 1, 0, 80236},
	{M, 1, 0, 1, 0},
	{M, 0, 1, 1, 0},
	{M, 1, 1, 1, 1},
	{M, M, 0, 1, M - 1},
	{M, M, 1, 1, M},
}

func testFunWW(t *testing.T, msg string, f funWW, a argWW) {
	z1, z0 := f(a.x, a.y, a.c)
	if z1 != a.z1 || z0 != a.z0 {
		t.Errorf("%s%+v\n\tgot z1:z0 = %#x:%#x; want %#x:%#x", msg, a, z1, z0, a.z1, a.z0)
	}
}

func TestFunWW(t *testing.T) {
	for _, a := range sumWW {
		arg := a
		testFunWW(t, "AddWW_g", AddWW_g, arg)

		arg = argWW{a.y, a.x, a.c, a.z1, a.z0}
		testFunWW(t, "AddWW_g symmetric", AddWW_g, arg)

		arg = argWW{a.z0, a.x, a.c, a.z1, a.y}
		testFunWW(t, "subWW_g", subWW_g, arg)

		arg = argWW{a.z0, a.y, a.c, a.z1, a.x}
		testFunWW(t, "subWW_g symmetric", subWW_g, arg)
	}
}

type funVV func(z, x, y []Word) (c Word)
type argVV struct {
	z, x, y nat
	c       Word
}

var sumVV = []argVV{
	{},
	{nat{0}, nat{0}, nat{0}, 0},
	{nat{1}, nat{1}, nat{0}, 0},
	{nat{0}, nat{M}, nat{1}, 1},
	{nat{80235}, nat{12345}, nat{67890}, 0},
	{nat{M - 1}, nat{M}, nat{M}, 1},
	{nat{0, 0, 0, 0}, nat{M, M, M, M}, nat{1, 0, 0, 0}, 1},
	{nat{0, 0, 0, M}, nat{M, M, M, M - 1}, nat{1, 0, 0, 0}, 0},
	{nat{0, 0, 0, 0}, nat{M, 0, M, 0}, nat{1, M, 0, M}, 1},
}

func testFunVV(t *testing.T, msg string, f funVV, a argVV) {
	z := make(nat, len(a.z))
	c := f(z, a.x, a.y)
	for i, zi := range z {
		if zi != a.z[i] {
			t.Errorf("%s%+v\n\tgot z[%d] = %#x; want %#x", msg, a, i, zi, a.z[i])
			break
		}
	}
	if c != a.c {
		t.Errorf("%s%+v\n\tgot c = %#x; want %#x", msg, a, c, a.c)
	}
}

func TestFunVV(t *testing.T) {
	for _, a := range sumVV {
		arg := a
		testFunVV(t, "AddVV_g", AddVV_g, arg)
		testFunVV(t, "AddVV", AddVV, arg)

		arg = argVV{a.z, a.y, a.x, a.c}
		testFunVV(t, "AddVV_g symmetric", AddVV_g, arg)
		testFunVV(t, "AddVV symmetric", AddVV, arg)

		arg = argVV{a.x, a.z, a.y, a.c}
		testFunVV(t, "SubVV_g", SubVV_g, arg)
		testFunVV(t, "SubVV", SubVV, arg)

		arg = argVV{a.y, a.z, a.x, a.c}
		testFunVV(t, "SubVV_g symmetric", SubVV_g, arg)
		testFunVV(t, "SubVV symmetric", SubVV, arg)
	}
}

// Always the same seed for reproducible results.
var rnd = rand.New(rand.NewSource(0))

func rndW() Word {
	return Word(rnd.Int63()<<1 | rnd.Int63n(2))
}

func rndV(n int) []Word {
	v := make([]Word, n)
	for i := range v {
		v[i] = rndW()
	}
	return v
}

var benchSizes = []int{1, 2, 3, 4, 5, 1e1, 1e2, 1e3, 1e4, 1e5}

func BenchmarkAddVV(b *testing.B) {
	for _, n := range benchSizes {
		if isRaceBuilder && n > 1e3 {
			continue
		}
		x := rndV(n)
		y := rndV(n)
		z := make([]Word, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.SetBytes(int64(n * _W))
			for i := 0; i < b.N; i++ {
				AddVV(z, x, y)
			}
		})
	}
}

type funVW func(z, x []Word, y Word) (c Word)
type argVW struct {
	z, x nat
	y    Word
	c    Word
}

var sumVW = []argVW{
	{},
	{nil, nil, 2, 2},
	{nat{0}, nat{0}, 0, 0},
	{nat{1}, nat{0}, 1, 0},
	{nat{1}, nat{1}, 0, 0},
	{nat{0}, nat{M}, 1, 1},
	{nat{0, 0, 0, 0}, nat{M, M, M, M}, 1, 1},
	{nat{585}, nat{314}, 271, 0},
}

var lshVW = []argVW{
	{},
	{nat{0}, nat{0}, 0, 0},
	{nat{0}, nat{0}, 1, 0},
	{nat{0}, nat{0}, 20, 0},

	{nat{M}, nat{M}, 0, 0},
	{nat{M << 1 & M}, nat{M}, 1, 1},
	{nat{M << 20 & M}, nat{M}, 20, M >> (_W - 20)},

	{nat{M, M, M}, nat{M, M, M}, 0, 0},
	{nat{M << 1 & M, M, M}, nat{M, M, M}, 1, 1},
	{nat{M << 20 & M, M, M}, nat{M, M, M}, 20, M >> (_W - 20)},
}

var rshVW = []argVW{
	{},
	{nat{0}, nat{0}, 0, 0},
	{nat{0}, nat{0}, 1, 0},
	{nat{0}, nat{0}, 20, 0},

	{nat{M}, nat{M}, 0, 0},
	{nat{M >> 1}, nat{M}, 1, M << (_W - 1) & M},
	{nat{M >> 20}, nat{M}, 20, M << (_W - 20) & M},

	{nat{M, M, M}, nat{M, M, M}, 0, 0},
	{nat{M, M, M >> 1}, nat{M, M, M}, 1, M << (_W - 1) & M},
	{nat{M, M, M >> 20}, nat{M, M, M}, 20, M << (_W - 20) & M},
}

func testFunVW(t *testing.T, msg string, f funVW, a argVW) {
	z := make(nat, len(a.z))
	c := f(z, a.x, a.y)
	for i, zi := range z {
		if zi != a.z[i] {
			t.Errorf("%s%+v\n\tgot z[%d] = %#x; want %#x", msg, a, i, zi, a.z[i])
			break
		}
	}
	if c != a.c {
		t.Errorf("%s%+v\n\tgot c = %#x; want %#x", msg, a, c, a.c)
	}
}

func makeFunVW(f func(z, x []Word, s uint) (c Word)) funVW {
	return func(z, x []Word, s Word) (c Word) {
		return f(z, x, uint(s))
	}
}

func TestFunVW(t *testing.T) {
	for _, a := range sumVW {
		arg := a
		testFunVW(t, "AddVW_g", AddVW_g, arg)
		testFunVW(t, "AddVW", AddVW, arg)

		arg = argVW{a.x, a.z, a.y, a.c}
		testFunVW(t, "SubVW_g", SubVW_g, arg)
		testFunVW(t, "SubVW", SubVW, arg)
	}

	shlVW_g := makeFunVW(ShlVU_g)
	shlVW := makeFunVW(ShlVU)
	for _, a := range lshVW {
		arg := a
		testFunVW(t, "ShlVU_g", shlVW_g, arg)
		testFunVW(t, "ShlVU", shlVW, arg)
	}

	shrVW_g := makeFunVW(ShrVU_g)
	shrVW := makeFunVW(ShrVU)
	for _, a := range rshVW {
		arg := a
		testFunVW(t, "ShrVU_g", shrVW_g, arg)
		testFunVW(t, "ShrVU", shrVW, arg)
	}
}

func BenchmarkAddVW(b *testing.B) {
	for _, n := range benchSizes {
		if isRaceBuilder && n > 1e3 {
			continue
		}
		x := rndV(n)
		y := rndW()
		z := make([]Word, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.SetBytes(int64(n * S))
			for i := 0; i < b.N; i++ {
				AddVW(z, x, y)
			}
		})
	}
}

type funVWW func(z, x []Word, y, r Word) (c Word)
type argVWW struct {
	z, x nat
	y, r Word
	c    Word
}

var prodVWW = []argVWW{
	{},
	{nat{0}, nat{0}, 0, 0, 0},
	{nat{991}, nat{0}, 0, 991, 0},
	{nat{0}, nat{M}, 0, 0, 0},
	{nat{991}, nat{M}, 0, 991, 0},
	{nat{0}, nat{0}, M, 0, 0},
	{nat{991}, nat{0}, M, 991, 0},
	{nat{1}, nat{1}, 1, 0, 0},
	{nat{992}, nat{1}, 1, 991, 0},
	{nat{22793}, nat{991}, 23, 0, 0},
	{nat{22800}, nat{991}, 23, 7, 0},
	{nat{0, 0, 0, 22793}, nat{0, 0, 0, 991}, 23, 0, 0},
	{nat{7, 0, 0, 22793}, nat{0, 0, 0, 991}, 23, 7, 0},
	{nat{0, 0, 0, 0}, nat{7893475, 7395495, 798547395, 68943}, 0, 0, 0},
	{nat{991, 0, 0, 0}, nat{7893475, 7395495, 798547395, 68943}, 0, 991, 0},
	{nat{0, 0, 0, 0}, nat{0, 0, 0, 0}, 894375984, 0, 0},
	{nat{991, 0, 0, 0}, nat{0, 0, 0, 0}, 894375984, 991, 0},
	{nat{M << 1 & M}, nat{M}, 1 << 1, 0, M >> (_W - 1)},
	{nat{M<<1&M + 1}, nat{M}, 1 << 1, 1, M >> (_W - 1)},
	{nat{M << 7 & M}, nat{M}, 1 << 7, 0, M >> (_W - 7)},
	{nat{M<<7&M + 1<<6}, nat{M}, 1 << 7, 1 << 6, M >> (_W - 7)},
	{nat{M << 7 & M, M, M, M}, nat{M, M, M, M}, 1 << 7, 0, M >> (_W - 7)},
	{nat{M<<7&M + 1<<6, M, M, M}, nat{M, M, M, M}, 1 << 7, 1 << 6, M >> (_W - 7)},
}

func testFunVWW(t *testing.T, msg string, f funVWW, a argVWW) {
	z := make(nat, len(a.z))
	c := f(z, a.x, a.y, a.r)
	for i, zi := range z {
		if zi != a.z[i] {
			t.Errorf("%s%+v\n\tgot z[%d] = %#x; want %#x", msg, a, i, zi, a.z[i])
			break
		}
	}
	if c != a.c {
		t.Errorf("%s%+v\n\tgot c = %#x; want %#x", msg, a, c, a.c)
	}
}

// TODO(gri) MulAddVWW and DivWVW are symmetric operations but
//           their signature is not symmetric. Try to unify.

type funWVW func(z []Word, xn Word, x []Word, y Word) (r Word)
type argWVW struct {
	z  nat
	xn Word
	x  nat
	y  Word
	r  Word
}

func testFunWVW(t *testing.T, msg string, f funWVW, a argWVW) {
	z := make(nat, len(a.z))
	r := f(z, a.xn, a.x, a.y)
	for i, zi := range z {
		if zi != a.z[i] {
			t.Errorf("%s%+v\n\tgot z[%d] = %#x; want %#x", msg, a, i, zi, a.z[i])
			break
		}
	}
	if r != a.r {
		t.Errorf("%s%+v\n\tgot r = %#x; want %#x", msg, a, r, a.r)
	}
}

func TestFunVWW(t *testing.T) {
	for _, a := range prodVWW {
		arg := a
		testFunVWW(t, "MulAddVWW_g", MulAddVWW_g, arg)
		testFunVWW(t, "MulAddVWW", MulAddVWW, arg)

		if a.y != 0 && a.r < a.y {
			arg := argWVW{a.x, a.c, a.z, a.y, a.r}
			testFunWVW(t, "DivWVW_g", DivWVW_g, arg)
			testFunWVW(t, "DivWVW", DivWVW, arg)
		}
	}
}

var MulWWTests = []struct {
	x, y Word
	q, r Word
}{
	{M, M, M - 1, 1},
	// 32 bit only: {0xc47dfa8c, 50911, 0x98a4, 0x998587f4},
}

func TestMulWW(t *testing.T) {
	for i, test := range MulWWTests {
		q, r := MulWW_g(test.x, test.y)
		if q != test.q || r != test.r {
			t.Errorf("#%d got (%x, %x) want (%x, %x)", i, q, r, test.q, test.r)
		}
	}
}

var MulAddWWWTests = []struct {
	x, y, c Word
	q, r    Word
}{
	// TODO(agl): These will only work on 64-bit platforms.
	// {15064310297182388543, 0xe7df04d2d35d5d80, 13537600649892366549, 13644450054494335067, 10832252001440893781},
	// {15064310297182388543, 0xdab2f18048baa68d, 13644450054494335067, 12869334219691522700, 14233854684711418382},
	{M, M, 0, M - 1, 1},
	{M, M, M, M, 0},
}

func TestMulAddWWW(t *testing.T) {
	for i, test := range MulAddWWWTests {
		q, r := MulAddWWW_g(test.x, test.y, test.c)
		if q != test.q || r != test.r {
			t.Errorf("#%d got (%x, %x) want (%x, %x)", i, q, r, test.q, test.r)
		}
	}
}

func BenchmarkAddMulVVW(b *testing.B) {
	for _, n := range benchSizes {
		if isRaceBuilder && n > 1e3 {
			continue
		}
		x := rndV(n)
		y := rndW()
		z := make([]Word, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.SetBytes(int64(n * _W))
			for i := 0; i < b.N; i++ {
				AddMulVVW(z, x, y)
			}
		})
	}
}
