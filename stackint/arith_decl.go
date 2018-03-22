// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !math_big_pure_go

package stackint

// implemented in arith_$GOARCH.s
func mulWW(x, y uint64) (z1, z0 uint64)
func divWW(x1, x0, y uint64) (q, r uint64)
func addVV(z, x, y []uint64) (c uint64)
func subVV(z, x, y []uint64) (c uint64)
func addVW(z, x []uint64, y uint64) (c uint64)
func subVW(z, x []uint64, y uint64) (c uint64)
func shlVU(z, x []uint64, s uint) (c uint64)
func shrVU(z, x []uint64, s uint) (c uint64)
func mulAddVWW(z, x []uint64, y, r uint64) (c uint64)
func addMulVVW(z, x []uint64, y uint64) (c uint64)
func divWVW(z []uint64, xn uint64, x []uint64, y uint64) (r uint64)
