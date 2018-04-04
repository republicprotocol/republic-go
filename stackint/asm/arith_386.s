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
	MOVL x+0(FP), AX
	MULL y+4(FP)
	MOVL DX, z1+8(FP)
	MOVL AX, z0+12(FP)
	RET


// func DivWW(x1, x0, y Word) (q, r Word)
TEXT ·DivWW(SB),NOSPLIT,$0
	MOVL x1+0(FP), DX
	MOVL x0+4(FP), AX
	DIVL y+8(FP)
	MOVL AX, q+12(FP)
	MOVL DX, r+16(FP)
	RET


// func AddVV(z, x, y []Word) (c Word)
TEXT ·AddVV(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), CX
	MOVL z_len+4(FP), BP
	MOVL $0, BX		// i = 0
	MOVL $0, DX		// c = 0
	JMP E1

L1:	MOVL (SI)(BX*4), AX
	ADDL DX, DX		// restore CF
	ADCL (CX)(BX*4), AX
	SBBL DX, DX		// save CF
	MOVL AX, (DI)(BX*4)
	ADDL $1, BX		// i++

E1:	CMPL BX, BP		// i < n
	JL L1

	NEGL DX
	MOVL DX, c+36(FP)
	RET


// func SubVV(z, x, y []Word) (c Word)
// (same as AddVV except for SBBL instead of ADCL and label names)
TEXT ·SubVV(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), CX
	MOVL z_len+4(FP), BP
	MOVL $0, BX		// i = 0
	MOVL $0, DX		// c = 0
	JMP E2

L2:	MOVL (SI)(BX*4), AX
	ADDL DX, DX		// restore CF
	SBBL (CX)(BX*4), AX
	SBBL DX, DX		// save CF
	MOVL AX, (DI)(BX*4)
	ADDL $1, BX		// i++

E2:	CMPL BX, BP		// i < n
	JL L2

	NEGL DX
	MOVL DX, c+36(FP)
	RET


// func AddVW(z, x []Word, y Word) (c Word)
TEXT ·AddVW(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), AX	// c = y
	MOVL z_len+4(FP), BP
	MOVL $0, BX		// i = 0
	JMP E3

L3:	ADDL (SI)(BX*4), AX
	MOVL AX, (DI)(BX*4)
	SBBL AX, AX		// save CF
	NEGL AX
	ADDL $1, BX		// i++

E3:	CMPL BX, BP		// i < n
	JL L3

	MOVL AX, c+28(FP)
	RET


// func SubVW(z, x []Word, y Word) (c Word)
TEXT ·SubVW(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), AX	// c = y
	MOVL z_len+4(FP), BP
	MOVL $0, BX		// i = 0
	JMP E4

L4:	MOVL (SI)(BX*4), DX
	SUBL AX, DX
	MOVL DX, (DI)(BX*4)
	SBBL AX, AX		// save CF
	NEGL AX
	ADDL $1, BX		// i++

E4:	CMPL BX, BP		// i < n
	JL L4

	MOVL AX, c+28(FP)
	RET


// func ShlVU(z, x []Word, s uint) (c Word)
TEXT ·ShlVU(SB),NOSPLIT,$0
	MOVL z_len+4(FP), BX	// i = z
	SUBL $1, BX		// i--
	JL X8b			// i < 0	(n <= 0)

	// n > 0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL s+24(FP), CX
	MOVL (SI)(BX*4), AX	// w1 = x[n-1]
	MOVL $0, DX
	SHLL CX, DX:AX		// w1>>ŝ
	MOVL DX, c+28(FP)

	CMPL BX, $0
	JLE X8a			// i <= 0

	// i > 0
L8:	MOVL AX, DX		// w = w1
	MOVL -4(SI)(BX*4), AX	// w1 = x[i-1]
	SHLL CX, DX:AX		// w<<s | w1>>ŝ
	MOVL DX, (DI)(BX*4)	// z[i] = w<<s | w1>>ŝ
	SUBL $1, BX		// i--
	JG L8			// i > 0

	// i <= 0
X8a:	SHLL CX, AX		// w1<<s
	MOVL AX, (DI)		// z[0] = w1<<s
	RET

X8b:	MOVL $0, c+28(FP)
	RET


// func ShrVU(z, x []Word, s uint) (c Word)
TEXT ·ShrVU(SB),NOSPLIT,$0
	MOVL z_len+4(FP), BP
	SUBL $1, BP		// n--
	JL X9b			// n < 0	(n <= 0)

	// n > 0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL s+24(FP), CX
	MOVL (SI), AX		// w1 = x[0]
	MOVL $0, DX
	SHRL CX, DX:AX		// w1<<ŝ
	MOVL DX, c+28(FP)

	MOVL $0, BX		// i = 0
	JMP E9

	// i < n-1
L9:	MOVL AX, DX		// w = w1
	MOVL 4(SI)(BX*4), AX	// w1 = x[i+1]
	SHRL CX, DX:AX		// w>>s | w1<<ŝ
	MOVL DX, (DI)(BX*4)	// z[i] = w>>s | w1<<ŝ
	ADDL $1, BX		// i++
	
E9:	CMPL BX, BP
	JL L9			// i < n-1

	// i >= n-1
X9a:	SHRL CX, AX		// w1>>s
	MOVL AX, (DI)(BP*4)	// z[n-1] = w1>>s
	RET

X9b:	MOVL $0, c+28(FP)
	RET


// func MulAddVWW(z, x []Word, y, r Word) (c Word)
TEXT ·MulAddVWW(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), BP
	MOVL r+28(FP), CX	// c = r
	MOVL z_len+4(FP), BX
	LEAL (DI)(BX*4), DI
	LEAL (SI)(BX*4), SI
	NEGL BX			// i = -n
	JMP E5

L5:	MOVL (SI)(BX*4), AX
	MULL BP
	ADDL CX, AX
	ADCL $0, DX
	MOVL AX, (DI)(BX*4)
	MOVL DX, CX
	ADDL $1, BX		// i++

E5:	CMPL BX, $0		// i < 0
	JL L5

	MOVL CX, c+32(FP)
	RET


// func AddMulVVW(z, x []Word, y Word) (c Word)
TEXT ·AddMulVVW(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL x+12(FP), SI
	MOVL y+24(FP), BP
	MOVL z_len+4(FP), BX
	LEAL (DI)(BX*4), DI
	LEAL (SI)(BX*4), SI
	NEGL BX			// i = -n
	MOVL $0, CX		// c = 0
	JMP E6

L6:	MOVL (SI)(BX*4), AX
	MULL BP
	ADDL CX, AX
	ADCL $0, DX
	ADDL AX, (DI)(BX*4)
	ADCL $0, DX
	MOVL DX, CX
	ADDL $1, BX		// i++

E6:	CMPL BX, $0		// i < 0
	JL L6

	MOVL CX, c+28(FP)
	RET


// func DivWVW(z* Word, xn Word, x []Word, y Word) (r Word)
TEXT ·DivWVW(SB),NOSPLIT,$0
	MOVL z+0(FP), DI
	MOVL xn+12(FP), DX	// r = xn
	MOVL x+16(FP), SI
	MOVL y+28(FP), CX
	MOVL z_len+4(FP), BX	// i = z
	JMP E7

L7:	MOVL (SI)(BX*4), AX
	DIVL CX
	MOVL AX, (DI)(BX*4)

E7:	SUBL $1, BX		// i--
	JGE L7			// i >= 0

	MOVL DX, r+32(FP)
	RET

