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

// +build !math_big_pure_go

#include "textflag.h"

// This file provides fast assembly versions for the elementary
// arithmetic operations on vectors implemented in arith.go.

// func MulWW(x, y Word) (z1, z0 Word)
TEXT ·MulWW(SB),NOSPLIT,$0
	MOVD	x+0(FP), R0
	MOVD	y+8(FP), R1
	MUL	R0, R1, R2
	UMULH	R0, R1, R3
	MOVD	R3, z1+16(FP)
	MOVD	R2, z0+24(FP)
	RET


// func DivWW(x1, x0, y Word) (q, r Word)
TEXT ·DivWW(SB),NOSPLIT,$0
	B	·DivWW_g(SB) // ARM64 has no multiword division


// func AddVV(z, x, y []Word) (c Word)
TEXT ·AddVV(SB),NOSPLIT,$0
	MOVD	z+0(FP), R3
	MOVD	z_len+8(FP), R0
	MOVD	x+24(FP), R1
	MOVD	y+48(FP), R2
	ADDS	$0, R0 // clear carry flag
loop:
	CBZ	R0, done // careful not to touch the carry flag
	MOVD.P	8(R1), R4
	MOVD.P	8(R2), R5
	ADCS	R4, R5
	MOVD.P	R5, 8(R3)
	SUB	$1, R0
	B	loop
done:
	CSET	HS, R0 // extract carry flag
	MOVD	R0, c+72(FP)
	RET


// func SubVV(z, x, y []Word) (c Word)
TEXT ·SubVV(SB),NOSPLIT,$0
	MOVD	z+0(FP), R3
	MOVD	z_len+8(FP), R0
	MOVD	x+24(FP), R1
	MOVD	y+48(FP), R2
	CMP	R0, R0 // set carry flag
loop:
	CBZ	R0, done // careful not to touch the carry flag
	MOVD.P	8(R1), R4
	MOVD.P	8(R2), R5
	SBCS	R5, R4
	MOVD.P	R4, 8(R3)
	SUB	$1, R0
	B	loop
done:
	CSET	LO, R0 // extract carry flag
	MOVD	R0, c+72(FP)
	RET


// func AddVW(z, x []Word, y Word) (c Word)
TEXT ·AddVW(SB),NOSPLIT,$0
	MOVD	z+0(FP), R3
	MOVD	z_len+8(FP), R0
	MOVD	x+24(FP), R1
	MOVD	y+48(FP), R2
	CBZ	R0, return_y
	MOVD.P	8(R1), R4
	ADDS	R2, R4
	MOVD.P	R4, 8(R3)
	SUB	$1, R0
loop:
	CBZ	R0, done // careful not to touch the carry flag
	MOVD.P	8(R1), R4
	ADCS	$0, R4
	MOVD.P	R4, 8(R3)
	SUB	$1, R0
	B	loop
done:
	CSET	HS, R0 // extract carry flag
	MOVD	R0, c+56(FP)
	RET
return_y: // z is empty; copy y to c
	MOVD	R2, c+56(FP)
	RET


// func SubVW(z, x []Word, y Word) (c Word)
TEXT ·SubVW(SB),NOSPLIT,$0
	MOVD	z+0(FP), R3
	MOVD	z_len+8(FP), R0
	MOVD	x+24(FP), R1
	MOVD	y+48(FP), R2
	CBZ	R0, rety
	MOVD.P	8(R1), R4
	SUBS	R2, R4
	MOVD.P	R4, 8(R3)
	SUB	$1, R0
loop:
	CBZ	R0, done // careful not to touch the carry flag
	MOVD.P	8(R1), R4
	SBCS	$0, R4
	MOVD.P	R4, 8(R3)
	SUB	$1, R0
	B	loop
done:
	CSET	LO, R0 // extract carry flag
	MOVD	R0, c+56(FP)
	RET
rety: // z is empty; copy y to c
	MOVD	R2, c+56(FP)
	RET


// func ShlVU(z, x []Word, s uint) (c Word)
TEXT ·ShlVU(SB),NOSPLIT,$0
	B ·ShlVU_g(SB)


// func ShrVU(z, x []Word, s uint) (c Word)
TEXT ·ShrVU(SB),NOSPLIT,$0
	B ·ShrVU_g(SB)


// func MulAddVWW(z, x []Word, y, r Word) (c Word)
TEXT ·MulAddVWW(SB),NOSPLIT,$0
	MOVD	z+0(FP), R1
	MOVD	z_len+8(FP), R0
	MOVD	x+24(FP), R2
	MOVD	y+48(FP), R3
	MOVD	r+56(FP), R4
loop:
	CBZ	R0, done
	MOVD.P	8(R2), R5
	UMULH	R5, R3, R7
	MUL	R5, R3, R6
	ADDS	R4, R6
	ADC	$0, R7
	MOVD.P	R6, 8(R1)
	MOVD	R7, R4
	SUB	$1, R0
	B	loop
done:
	MOVD	R4, c+64(FP)
	RET


// func AddMulVVW(z, x []Word, y Word) (c Word)
TEXT ·AddMulVVW(SB),NOSPLIT,$0
	B ·AddMulVVW_g(SB)


// func DivWVW(z []Word, xn Word, x []Word, y Word) (r Word)
TEXT ·DivWVW(SB),NOSPLIT,$0
	B ·DivWVW_g(SB)

