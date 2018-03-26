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

// +build !math_big_pure_go,mips !math_big_pure_go,mipsle

#include "textflag.h"

// This file provides fast assembly versions for the elementary
// arithmetic operations on vectors implemented in arith.go.

TEXT ·mulWW(SB),NOSPLIT,$0
	JMP	·mulWW_g(SB)

TEXT ·divWW(SB),NOSPLIT,$0
	JMP	·divWW_g(SB)

TEXT ·addVV(SB),NOSPLIT,$0
	JMP	·addVV_g(SB)

TEXT ·subVV(SB),NOSPLIT,$0
	JMP	·subVV_g(SB)

TEXT ·addVW(SB),NOSPLIT,$0
	JMP	·addVW_g(SB)

TEXT ·subVW(SB),NOSPLIT,$0
	JMP	·subVW_g(SB)

TEXT ·shlVU(SB),NOSPLIT,$0
	JMP	·shlVU_g(SB)

TEXT ·shrVU(SB),NOSPLIT,$0
	JMP	·shrVU_g(SB)

TEXT ·mulAddVWW(SB),NOSPLIT,$0
	JMP	·mulAddVWW_g(SB)

TEXT ·addMulVVW(SB),NOSPLIT,$0
	JMP	·addMulVVW_g(SB)

TEXT ·divWVW(SB),NOSPLIT,$0
	JMP	·divWVW_g(SB)

