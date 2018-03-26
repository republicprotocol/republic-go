// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !math_big_pure_go

package asm

var MulWW = mulWW
var AddVV = addVV_g
var AddVW = addVW_g
var AddWW = addWW_g
var SubVV = subVV_g
var SubVW = subVW_g
var MulAddVWW = mulAddVWW_g
var MulAddWWW = mulAddWWW_g
var DivWW = divWW_g
var ShlVU = shlVU_g
var ShrVU = shrVU_g
var Nlz = nlz

const M Word = _M
const S int = _S

// (x, y Word) (z1, z0 Word)
// func divWW(x1, x0, y Word) (q, r Word)
// func addVV(z, x, y []Word) (c Word)
// func subVV(z, x, y []Word) (c Word)
// func addVW(z, x []Word, y Word) (c Word)
// func subVW(z, x []Word, y Word) (c Word)
// func shlVU(z, x []Word, s uint) (c Word)
// func shrVU(z, x []Word, s uint) (c Word)
// func mulAddVWW(z, x []Word, y, r Word) (c Word)
// func addMulVVW(z, x []Word, y Word) (c Word)
// func divWVW(z []Word, xn Word, x []Word, y Word) (r Word)
