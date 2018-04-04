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

// +build !math_big_pure_go,ppc64 !math_big_pure_go,ppc64le

#include "textflag.h"

// This file provides fast assembly versions for the elementary
// arithmetic operations on vectors implemented in arith.go.

// func MulWW(x, y Word) (z1, z0 Word)
TEXT ·MulWW(SB), NOSPLIT, $0
	MOVD   x+0(FP), R4
	MOVD   y+8(FP), R5
	MULHDU R4, R5, R6
	MULLD  R4, R5, R7
	MOVD   R6, z1+16(FP)
	MOVD   R7, z0+24(FP)
	RET

// func AddVV(z, y, y []Word) (c Word)
// z[i] = x[i] + y[i] for all i, carrying
TEXT ·AddVV(SB), NOSPLIT, $0
	MOVD  z_len+8(FP), R7
	MOVD  x+24(FP), R8
	MOVD  y+48(FP), R9
	MOVD  z+0(FP), R10

	MOVD  R0, R4
	MOVD  R0, R6  // R6 will be the address index
	ADDC R4, R4   // clear CA
	MOVD  R7, CTR

	CMP   R0, R7
	BEQ   done

loop:
	MOVD  (R8)(R6), R11   // x[i]
	MOVD  (R9)(R6), R12   // y[i]
	ADDE  R12, R11, R15   // x[i] + y[i] + CA
	MOVD  R15, (R10)(R6)  // z[i]

	ADD $8, R6
	BC  16, 0, loop	// bdnz

done:
	ADDZE R4
	MOVD  R4, c+72(FP)
	RET

// func SubVV(z, x, y []Word) (c Word)
// z[i] = x[i] - y[i] for all i, carrying
TEXT ·SubVV(SB), NOSPLIT, $0
	MOVD z_len+8(FP), R7
	MOVD x+24(FP), R8
	MOVD y+48(FP), R9
	MOVD z+0(FP), R10

	MOVD  R0, R4  // c = 0
	MOVD  R0, R6
	SUBC R0, R0  // clear CA
	MOVD  R7, CTR

	CMP R0, R7
	BEQ  sublend

// amd64 saves and restores CF, but I believe they only have to do that because all of
// their math operations clobber it - we should just be able to recover it at the end.
subloop:
	MOVD  (R8)(R6), R11 // x[i]
	MOVD  (R9)(R6), R12 // y[i]

	SUBE R12, R11, R15
	MOVD R15, (R10)(R6)

	ADD $8, R6
	BC  16, 0, subloop  // bdnz

sublend:

	ADDZE R4
	XOR   $1, R4
	MOVD  R4, c+72(FP)
	RET

TEXT ·AddVW(SB), NOSPLIT, $0
	BR ·AddVW_g(SB)

TEXT ·SubVW(SB), NOSPLIT, $0
	BR ·SubVW_g(SB)

TEXT ·ShlVU(SB), NOSPLIT, $0
	BR ·ShlVU_g(SB)

TEXT ·ShrVU(SB), NOSPLIT, $0
	BR ·ShrVU_g(SB)

// func MulAddVWW(z, x []Word, y, r Word) (c Word)
TEXT ·MulAddVWW(SB), NOSPLIT, $0
	MOVD z+0(FP), R10	// R10 = z[]
	MOVD x+24(FP), R8	// R8 = x[]
	MOVD y+48(FP), R9	// R9 = y
	MOVD r+56(FP), R4	// R4 = r = c
	MOVD z_len+8(FP), R11	// R11 = z_len

	MOVD R0, R3		// R3 will be the index register
	CMP  R0, R11
	MOVD R11, CTR		// Initialize loop counter
	BEQ  done

loop:
	MOVD   (R8)(R3), R20	// x[i]
	MULLD  R9, R20, R6	// R6 = z0 = Low-order(x[i]*y)
	MULHDU R9, R20, R7	// R7 = z1 = High-order(x[i]*y)
	ADDC   R4, R6		// Compute sum for z1 and z0
	ADDZE  R7
	MOVD   R6, (R10)(R3)	// z[i]
	MOVD   R7, R4		// c
	ADD    $8, R3
	BC  16, 0, loop		// bdnz

done:
	MOVD R4, c+64(FP)
	RET

// func AddMulVVW(z, x []Word, y Word) (c Word)
TEXT ·AddMulVVW(SB), NOSPLIT, $0
	MOVD z+0(FP), R10	// R10 = z[]
	MOVD x+24(FP), R8	// R8 = x[]
	MOVD y+48(FP), R9	// R9 = y
	MOVD z_len+8(FP), R22	// R22 = z_len

	MOVD R0, R3		// R3 will be the index register
	CMP  R0, R22
	MOVD R0, R4		// R4 = c = 0
	MOVD R22, CTR		// Initialize loop counter
	BEQ  done

loop:
	MOVD  (R8)(R3), R20	// Load x[i]
	MOVD  (R10)(R3), R21	// Load z[i]
	MULLD  R9, R20, R6	// R6 = Low-order(x[i]*y)
	MULHDU R9, R20, R7	// R7 = High-order(x[i]*y)
	ADDC   R21, R6		// R6 = z0
	ADDZE  R7		// R7 = z1
	ADDC   R4, R6		// R6 = z0 + c + 0
	ADDZE  R7, R4           // c += z1
	MOVD   R6, (R10)(R3)	// Store z[i]
	ADD    $8, R3
	BC  16, 0, loop		// bdnz

done:
	MOVD R4, c+56(FP)
	RET

// func DivWW(x1, x0, y Word) (q, r Word)
TEXT ·DivWW(SB), NOSPLIT, $0
	MOVD x1+0(FP), R4
	MOVD x0+8(FP), R5
	MOVD y+16(FP), R6

	CMPU R4, R6
	BGE  divbigger

	// from the programmer's note in ch. 3 of the ISA manual, p.74
	DIVDEU R6, R4, R3
	DIVDU  R6, R5, R7
	MULLD  R6, R3, R8
	MULLD  R6, R7, R20
	SUB    R20, R5, R10
	ADD    R7, R3, R3
	SUB    R8, R10, R4
	CMPU   R4, R10
	BLT    adjust
	CMPU   R4, R6
	BLT    end

adjust:
	MOVD $1, R21
	ADD  R21, R3, R3
	SUB  R6, R4, R4

end:
	MOVD R3, q+24(FP)
	MOVD R4, r+32(FP)

	RET

divbigger:
	MOVD $-1, R7
	MOVD R7, q+24(FP)
	MOVD R7, r+32(FP)
	RET

TEXT ·DivWVW(SB), NOSPLIT, $0
	BR ·DivWVW_g(SB)

