package stackint

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// go build -a -gcflags='-m -m' int1024.go int1024_arithmetic.go int1024_bitwise.go int1024_comparison.go int1024_internal.go

// SIZE is the number of bits stored by Int1024
const SIZE = 1024

// WORDSIZE is 64 for Word
const WORDSIZE = 64

// Word is the internal type
type Word uint64

// WORDMAX represents the largest word value
const WORDMAX = 1<<WORDSIZE - 1

// INT1024WORDS is 1024 / 64 = 16
const INT1024WORDS = SIZE / WORDSIZE

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
	words [INT1024WORDS]Word
}

// FromUint64 returns a new Int1024 from a Word
func FromUint64(n uint64) Int1024 {
	z := Zero()
	z.words[INT1024WORDS-1] = Word(n)
	return z
}

// ToUint64 converts an Int1024 to a uint64 if it is small enough
func (x *Int1024) ToUint64() uint64 {
	// Check that all other words are zero
	for i := 0; i < INT1024WORDS-2; i++ {
		if x.words[i] != 0 {
			panic("Int1024 is too large to be converted to uint64")
		}
	}
	return uint64(x.words[INT1024WORDS-1])
}

// FromString returns a new Int1024 from a string
func FromString(number string) Int1024 {
	self := Zero()

	// Base to convert from
	base := 10

	if len(number) > 2 {
		if number[0:2] == "0x" {
			number = number[2:]
			base = 16
		} else if number[0:2] == "0b" {
			number = number[2:]
			base = 2
		}
	}

	// Length of string
	length := len(number)

	if length == 0 {
		panic("invalid number")
	}

	// Break up into blocks of size 19 (log10(2 ** 64))
	blockSize := 1
	limit := uint64(1<<63-1) / uint64(base)
	for basePower := uint64(1); basePower < limit; basePower *= uint64(base) {
		blockSize++
	}

	// Number of blocks
	count := (length / blockSize)
	if length%blockSize != 0 {
		count++
	}

	// TODO: Replace with 10.Pow(blockSize)
	shift := FromUint64(uint64(base))
	blockSizeInt := FromUint64(uint64(blockSize))
	shift = shift.Exp(&blockSizeInt)
	shiftAcc := One()

	// Loop through each block. Multiply block by (10**19)**i and add to number.
	for i := 0; i < count; i++ {
		end := length - i*blockSize
		start := length - (i+1)*blockSize
		if start < 0 {
			start = 0
		}
		word, err := strconv.ParseUint(number[start:end], base, 64)
		if err != nil {
			panic(err)
		}

		wordI := FromUint64(word)
		wordI = wordI.Mul(&shiftAcc)
		self = self.Add(&wordI)

		shiftAcc = shiftAcc.Mul(&shift)
	}

	return self
}

func (x *Int1024) String() string {
	blockSize := 19
	blockSize1024 := FromUint64(uint64(blockSize))
	base := FromUint64(10)
	base = base.Exp(&blockSize1024)
	q := *x
	var r Int1024
	ret := ""

	for !q.IsZero() {
		q, r = q.DivMod(&base)
		chars := strconv.FormatUint(uint64(r.words[INT1024WORDS-1]), 10)
		zeroes := strings.Repeat("0", blockSize-len(chars))
		ret = zeroes + chars + ret
	}

	ret = strings.TrimLeft(ret, "0")

	if ret == "" {
		return "0"
	}

	return ret
}

// ToBytes returns an array of 128 bytes (Big Endian)
func (x *Int1024) ToBytes() []byte {

	bytes128 := make([]byte, 128)
	b8 := make([]byte, 8)

	for i := range x.words {
		binary.BigEndian.PutUint64(b8, uint64(x.words[i]))
		for j := 0; j < 8; j++ {
			bytes128[i*8+j] = b8[j]
		}
	}

	return bytes128
}

// FromBytes deserializes an array of 128 bytes to an Int1024 (Big Endian)
func FromBytes(bytes128 []byte) Int1024 {

	x := Zero()

	numWords := len(bytes128) / 8
	if len(bytes128)%8 != 0 {
		numWords++
	}

	// mod := 8 - len(bytes128)%8

	for i := 0; i < numWords; i++ {
		b8 := make([]byte, 8)
		start := len(bytes128) - i*8
		end := start - 8
		if end < 0 {
			end = 0
		}
		for j := 0; j < start-end; j++ {
			b8[7-j] = bytes128[start-j-1]
		}
		x.words[INT1024WORDS-1-i] = Word(binary.BigEndian.Uint64(b8))
	}

	return x
}

// ToLittleEndianBytes returns an array of 128 bytes (Little Endian)
func (x *Int1024) ToLittleEndianBytes() []byte {

	bytes128 := make([]byte, 128)
	b8 := make([]byte, 8)

	for i := range x.words {
		binary.LittleEndian.PutUint64(b8, uint64(x.words[INT1024WORDS-1-i]))
		for j := 0; j < 8; j++ {
			bytes128[i*8+j] = b8[j]
		}
	}

	return bytes128
}

// FromLittleEndianBytes deserializes an array of 128 bytes to an Int1024 (LittleBig Endian)
func FromLittleEndianBytes(bytes128 []byte) Int1024 {

	x := Zero()

	numWords := len(bytes128) / 8

	for i := 0; i < numWords; i++ {
		b8 := bytes128[i*8 : (i+1)*8]
		x.words[INT1024WORDS-1-i] = Word(binary.LittleEndian.Uint64(b8))
	}

	return x
}

// Clone returns a new Int1024 representing the same value as x
func (x *Int1024) Clone() Int1024 {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = x.words[i]
	}
	return Int1024{
		words: words,
	}
}

// Words returns a clone of the [16]Word used by x as its internal representation
func (x *Int1024) Words() [INT1024WORDS]Word {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = x.words[i]
	}
	return words
}

// ToBinary returns the binary representation of x as a string
func (x *Int1024) ToBinary() string {
	str := ""
	started := false
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] == 0 && !started {
			continue
		}
		if !started {
			started = true
			// First time around don't print leading zeros
			str = str + fmt.Sprintf("%b", x.words[i])
		} else {
			str = str + fmt.Sprintf("%064b", x.words[i])
		}
	}
	if str == "" {
		return "0"
	}
	return str
}

/* CONSTANTS */

// Zero returns a new Int1024 representing 0
func Zero() Int1024 {
	var words [INT1024WORDS]Word
	// for i := 0; i < INT1024WORDS; i++ {
	// 	words[i] = 0
	// }
	return Int1024{
		words: words,
	}
}

// One is the Int1024 that represents 1
var One = func() Int1024 { return FromUint64(1) }

// Two is the Int1024 that represents 2
var Two = func() Int1024 { return FromUint64(2) }

// Do not call overwriting functions on these!
var zero = Zero()
var one = One()
var two = Two()

// maxInt returns a new Int1024 representing 2**1024 - 1
func maxInt() Int1024 {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = WORDMAX
	}
	return Int1024{
		words: words,
	}
}

// MAXINT1024 is the Int1024 that represents 2**1024 - 1
var MAXINT1024 = func() Int1024 { return maxInt() }

// maxInt returns a new Int1024 representing 2**1024 - 1
func twoPow1023() Int1024 {
	var words [INT1024WORDS]Word
	words[0] = 1 << 63
	return Int1024{
		words: words,
	}
}

// TWOPOW1023 is the Int1024 that represents 2**1023
var TWOPOW1023 = func() Int1024 { return twoPow1023() }
