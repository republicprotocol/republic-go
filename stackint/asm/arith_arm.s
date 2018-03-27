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

// func AddVV(z, x, y []Word) (c Word)
TEXT ·AddVV(SB),NOSPLIT,$0
	ADD.S	$0, R0		// clear carry flag
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R4
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	ADD	R4<<2, R1, R4
	B E1
L1:
	MOVW.P	4(R2), R5
	MOVW.P	4(R3), R6
	ADC.S	R6, R5
	MOVW.P	R5, 4(R1)
E1:
	TEQ	R1, R4
	BNE L1

	MOVW	$0, R0
	MOVW.CS	$1, R0
	MOVW	R0, c+36(FP)
	RET


// func SubVV(z, x, y []Word) (c Word)
// (same as AddVV except for SBC instead of ADC and label names)
TEXT ·SubVV(SB),NOSPLIT,$0
	SUB.S	$0, R0		// clear borrow flag
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R4
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	ADD	R4<<2, R1, R4
	B E2
L2:
	MOVW.P	4(R2), R5
	MOVW.P	4(R3), R6
	SBC.S	R6, R5
	MOVW.P	R5, 4(R1)
E2:
	TEQ	R1, R4
	BNE L2

	MOVW	$0, R0
	MOVW.CC	$1, R0
	MOVW	R0, c+36(FP)
	RET


// func AddVW(z, x []Word, y Word) (c Word)
TEXT ·AddVW(SB),NOSPLIT,$0
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R4
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	ADD	R4<<2, R1, R4
	TEQ	R1, R4
	BNE L3a
	MOVW	R3, c+28(FP)
	RET
L3a:
	MOVW.P	4(R2), R5
	ADD.S	R3, R5
	MOVW.P	R5, 4(R1)
	B	E3
L3:
	MOVW.P	4(R2), R5
	ADC.S	$0, R5
	MOVW.P	R5, 4(R1)
E3:
	TEQ	R1, R4
	BNE	L3

	MOVW	$0, R0
	MOVW.CS	$1, R0
	MOVW	R0, c+28(FP)
	RET


// func SubVW(z, x []Word, y Word) (c Word)
TEXT ·SubVW(SB),NOSPLIT,$0
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R4
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	ADD	R4<<2, R1, R4
	TEQ	R1, R4
	BNE L4a
	MOVW	R3, c+28(FP)
	RET
L4a:
	MOVW.P	4(R2), R5
	SUB.S	R3, R5
	MOVW.P	R5, 4(R1)
	B	E4
L4:
	MOVW.P	4(R2), R5
	SBC.S	$0, R5
	MOVW.P	R5, 4(R1)
E4:
	TEQ	R1, R4
	BNE	L4

	MOVW	$0, R0
	MOVW.CC	$1, R0
	MOVW	R0, c+28(FP)
	RET


// func ShlVU(z, x []Word, s uint) (c Word)
TEXT ·ShlVU(SB),NOSPLIT,$0
	MOVW	z_len+4(FP), R5
	TEQ	$0, R5
	BEQ	X7
	
	MOVW	z+0(FP), R1
	MOVW	x+12(FP), R2
	ADD	R5<<2, R2, R2
	ADD	R5<<2, R1, R5
	MOVW	s+24(FP), R3
	TEQ	$0, R3	// shift 0 is special
	BEQ	Y7
	ADD	$4, R1	// stop one word early
	MOVW	$32, R4
	SUB	R3, R4
	MOVW	$0, R7
	
	MOVW.W	-4(R2), R6
	MOVW	R6<<R3, R7
	MOVW	R6>>R4, R6
	MOVW	R6, c+28(FP)
	B E7

L7:
	MOVW.W	-4(R2), R6
	ORR	R6>>R4, R7
	MOVW.W	R7, -4(R5)
	MOVW	R6<<R3, R7
E7:
	TEQ	R1, R5
	BNE	L7

	MOVW	R7, -4(R5)
	RET

Y7:	// copy loop, because shift 0 == shift 32
	MOVW.W	-4(R2), R6
	MOVW.W	R6, -4(R5)
	TEQ	R1, R5
	BNE Y7

X7:
	MOVW	$0, R1
	MOVW	R1, c+28(FP)
	RET


// func ShrVU(z, x []Word, s uint) (c Word)
TEXT ·ShrVU(SB),NOSPLIT,$0
	MOVW	z_len+4(FP), R5
	TEQ	$0, R5
	BEQ	X6

	MOVW	z+0(FP), R1
	MOVW	x+12(FP), R2
	ADD	R5<<2, R1, R5
	MOVW	s+24(FP), R3
	TEQ	$0, R3	// shift 0 is special
	BEQ Y6
	SUB	$4, R5	// stop one word early
	MOVW	$32, R4
	SUB	R3, R4
	MOVW	$0, R7

	// first word
	MOVW.P	4(R2), R6
	MOVW	R6>>R3, R7
	MOVW	R6<<R4, R6
	MOVW	R6, c+28(FP)
	B E6

	// word loop
L6:
	MOVW.P	4(R2), R6
	ORR	R6<<R4, R7
	MOVW.P	R7, 4(R1)
	MOVW	R6>>R3, R7
E6:
	TEQ	R1, R5
	BNE	L6

	MOVW	R7, 0(R1)
	RET

Y6:	// copy loop, because shift 0 == shift 32
	MOVW.P	4(R2), R6
	MOVW.P	R6, 4(R1)
	TEQ R1, R5
	BNE Y6

X6:
	MOVW	$0, R1
	MOVW	R1, c+28(FP)
	RET


// func MulAddVWW(z, x []Word, y, r Word) (c Word)
TEXT ·MulAddVWW(SB),NOSPLIT,$0
	MOVW	$0, R0
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R5
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	MOVW	r+28(FP), R4
	ADD	R5<<2, R1, R5
	B E8

	// word loop
L8:
	MOVW.P	4(R2), R6
	MULLU	R6, R3, (R7, R6)
	ADD.S	R4, R6
	ADC	R0, R7
	MOVW.P	R6, 4(R1)
	MOVW	R7, R4
E8:
	TEQ	R1, R5
	BNE	L8

	MOVW	R4, c+32(FP)
	RET


// func AddMulVVW(z, x []Word, y Word) (c Word)
TEXT ·AddMulVVW(SB),NOSPLIT,$0
	MOVW	$0, R0
	MOVW	z+0(FP), R1
	MOVW	z_len+4(FP), R5
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	ADD	R5<<2, R1, R5
	MOVW	$0, R4
	B E9

	// word loop
L9:
	MOVW.P	4(R2), R6
	MULLU	R6, R3, (R7, R6)
	ADD.S	R4, R6
	ADC	R0, R7
	MOVW	0(R1), R4
	ADD.S	R4, R6
	ADC	R0, R7
	MOVW.P	R6, 4(R1)
	MOVW	R7, R4
E9:
	TEQ	R1, R5
	BNE	L9

	MOVW	R4, c+28(FP)
	RET


// func DivWVW(z* Word, xn Word, x []Word, y Word) (r Word)
TEXT ·DivWVW(SB),NOSPLIT,$0
	// ARM has no multiword division, so use portable code.
	B ·DivWVW_g(SB)


// func DivWW(x1, x0, y Word) (q, r Word)
TEXT ·DivWW(SB),NOSPLIT,$0
	// ARM has no multiword division, so use portable code.
	B ·DivWW_g(SB)


// func MulWW(x, y Word) (z1, z0 Word)
TEXT ·MulWW(SB),NOSPLIT,$0
	MOVW	x+0(FP), R1
	MOVW	y+4(FP), R2
	MULLU	R1, R2, (R4, R3)
	MOVW	R4, z1+8(FP)
	MOVW	R3, z0+12(FP)
	RET

