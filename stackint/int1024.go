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
	z := zero()
	z.words[INT1024WORDS-1] = Word(n)
	return z
}

// FromString returns a new Int1024 from a string
func FromString(number string) Int1024 {
	self := zero()

	// Length of string
	length := len(number)

	// Break up into blocks of size 19 (log10(2 ** 64))
	blockSize := 19

	// Number of blocks
	count := (length / blockSize)
	if length%blockSize != 0 {
		count++
	}

	// TODO: Replace with 10.Pow(blockSize)
	shift := FromUint64(10)
	blockSizeInt := FromUint64(uint64(blockSize))
	shift = shift.Exp(&blockSizeInt)
	shiftAcc := ONE

	// Loop through each block. Multiply block by (10**19)**i and add to number.
	for i := 0; i < count; i++ {
		end := length - i*blockSize
		start := length - (i+1)*blockSize
		if start < 0 {
			start = 0
		}
		word, err := strconv.ParseUint(number[start:end], 10, 64)
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

	x := zero()

	numWords := len(bytes128) / 8
	if len(bytes128)%8 != 0 {
		numWords++
	}

	for i := 0; i < numWords; i++ {
		k := (numWords - 1) - i
		b8 := bytes128[k*8 : (k+1)*8]
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

	x := zero()

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

// Words returns the internal [16]Word that stores the value of x
func (x *Int1024) Words() [INT1024WORDS]Word {
	return x.words
}

// ToBinary returns the binary representation of x as a string
func (x *Int1024) ToBinary() string {
	str := ""
	for i := 0; i < INT1024WORDS; i++ {
		str = fmt.Sprintf(str+"%064b", x.words[i])
	}
	stripped := strings.TrimLeft(str, "0")
	if stripped == "" {
		return "0"
	}
	return stripped
}

/* CONSTANTS */

// zero returns a new Int1024 representing 0
func zero() Int1024 {
	var words [INT1024WORDS]Word
	// for i := 0; i < INT1024WORDS; i++ {
	// 	words[i] = 0
	// }
	return Int1024{
		words: words,
	}
}

// ZERO is the Int1024 that represents 0
var ZERO = zero()

// ONE is the Int1024 that represents 1
var ONE = FromUint64(1)

// TWO is the Int1024 that represents 2
var TWO = FromUint64(2)

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
var MAXINT1024 = maxInt()

// maxInt returns a new Int1024 representing 2**1024 - 1
func twoPow1023() Int1024 {
	var words [INT1024WORDS]Word
	words[0] = 1 << 63
	return Int1024{
		words: words,
	}
}

// TWOPOW1023 is the Int1024 that represents 2**1023
var TWOPOW1023 = twoPow1023()
