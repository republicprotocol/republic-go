package stackint

import (
	"math/bits"
)

// ShiftLeft returns x<<n
func (x *Int1024) ShiftLeft(n uint) Int1024 {
	z := x.Clone()
	z.ShiftLeftInPlace(n)
	return z
}

// ShiftLeftInPlace shifts to the left x by one
func (x *Int1024) ShiftLeftInPlace(n uint) {

	// expected := big.NewInt(0).Lsh(x.ToBigInt(), n)

	// If n > 64, first, shift entire words
	div := n / WORDSIZE
	if div > 0 {
		words := min(INT1024WORDS, x.length+uint16(div))
		var firstPositive uint16
		var i uint
		for i = uint(words) - 1; i >= div; i-- {
			x.words[i] = x.words[i-div]
			if x.words[i] != 0 && firstPositive == 0 {
				firstPositive = uint16(i)
			}
		}
		for i = 0; i < div; i++ {
			x.words[i] = 0
		}
		x.length = firstPositive + 1
		n = n - div*WORDSIZE
	}

	if n == 1 {
		x.shiftleftone()
	} else {
		x.shiftleft(n)
	}

	// actual := x.ToBigInt()
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
	// 	panic(fmt.Sprintf("Shiftleft failed!\nFor %v << %v\n.\n\nExp: %b\n\nGot: %b", x, n, expected, actual))
	// }

}

func (x *Int1024) shiftleft(n uint) {
	if n >= SIZE {
		panic("shifting by more than a word")
	}
	var overflow uint64
	var shift uint64 = (1<<n - 1)
	// fmt.Println(shift)
	var firstPositive uint16
	var i uint16
	for i = 0; i < x.length; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - n)) & shift
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] << n) | overflow
		if x.words[i] != 0 {
			firstPositive = i
		}
		overflow = newOverflow
	}
	if overflow != 0 && x.length < INT1024WORDS {
		x.length++
		x.words[x.length-1] = overflow
	} else {
		x.length = firstPositive + 1
	}
}

func (x *Int1024) shiftleftone() {
	var overflow uint64
	var firstPositive uint16
	var i uint16
	for i = 0; i < x.length; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - 1)) & 1
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] << 1) | overflow
		if x.words[i] != 0 {
			firstPositive = i
		}
		overflow = newOverflow
	}
	if overflow != 0 && x.length < INT1024WORDS {
		firstPositive = i
		x.words[i] = overflow
	}
	x.length = firstPositive + 1

}

// ShiftRight returns x>>n
func (x *Int1024) ShiftRight(n uint) Int1024 {
	z := x.Clone()
	z.ShiftRightInPlace(n)
	return z
}

// ShiftRightInPlace shifts to the right x by one
func (x *Int1024) ShiftRightInPlace(n uint) {

	// expected := big.NewInt(0).Rsh(x.ToBigInt(), n)

	// If n > 64, first, shift entire words
	div := n / WORDSIZE
	if div > 0 {
		previous := x.length
		// uint overflows
		if uint16(div) > x.length {
			x.length = 0
		} else {
			x.length = x.length - uint16(div)
		}
		var i uint16
		for i = 0; i < x.length; i++ {
			x.words[i] = x.words[i+uint16(div)]
		}
		for i = x.length; i < previous; i++ {
			x.words[i] = 0
		}
		if x.length == 0 {
			x.length = 1
		}
		n = n - div*WORDSIZE
	}

	if n == 1 {
		x.shiftrightone()
	} else {
		x.shiftright(n)
	}

	// actual := x.ToBigInt()
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
	// 	panic(fmt.Sprintf("Shiftleft failed!\nFor %v << %v\n.\n\nExp: %b\n\nGot: %b", x, n, expected, actual))
	// }

}

func (x *Int1024) shiftright(n uint) {
	if n >= SIZE {
		panic("shifting by more than a word")
	}
	var overflow uint64
	var shift uint64 = (1<<n - 1)
	for i := int16(x.length - 1); i >= 0; i-- {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] & shift) << (WORDSIZE - n)
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] >> n) | overflow
		overflow = newOverflow
	}
	if x.words[x.length-1] == 0 && x.length > 1 {
		x.length--
	}
}

func (x *Int1024) shiftrightone() {
	overflow := uint64(0)
	for i := int16(x.length - 1); i >= 0; i-- {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] & 1) << (WORDSIZE - 1)
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] >> 1) | overflow
		overflow = newOverflow
	}
	if x.words[x.length-1] == 0 && x.length > 1 {
		x.length--
	}

}

func (x *Int1024) IsBitSet(n int) bool {
	if n > (SIZE-1) || n < 0 {
		return false
	}
	word := n / WORDSIZE
	if uint16(word) >= x.length {
		return false
	}
	bit := uint(n % WORDSIZE)
	return x.words[word]&(1<<bit) != 0
}

// AND returns x&y
func (x *Int1024) AND(y *Int1024) Int1024 {
	min := x.length
	yLen := y.length
	if yLen < min {
		min = yLen
	}

	z := Zero()
	var firstPositive uint16
	var i uint16
	for i = 0; i < min; i++ {
		z.words[i] = x.words[i] & y.words[i]
		if z.words[i] != 0 {
			firstPositive = i
		}
	}
	z.length = firstPositive + 1
	return z
}

// OR returns x|y
func (x *Int1024) OR(y *Int1024) Int1024 {
	z := x.Clone()
	z.ORInPlace(y)
	return z
}

// ORInPlace sets x to x|y
func (x *Int1024) ORInPlace(y *Int1024) {

	min := x.length
	max := y.length
	maxi := y
	if max < min {
		min, max = max, min
		maxi = x
	}

	var i uint16
	var firstPositive uint16
	for i = 0; i < min; i++ {
		x.words[i] = x.words[i] | y.words[i]
		if x.words[i] != 0 {
			firstPositive = i
		}
	}
	for i = min; i < max; i++ {
		x.words[i] = maxi.words[i]
		if x.words[i] != 0 {
			firstPositive = i
		}
	}

	x.length = firstPositive + 1
}

// XOR returns x&y
func (x *Int1024) XOR(y *Int1024) Int1024 {
	min := x.length
	max := y.length
	maxi := y
	if max < min {
		min, max = max, min
		maxi = x
	}

	var words [INT1024WORDS]uint64
	var i uint16
	var firstPositive uint16
	for i = 0; i < min; i++ {
		words[i] = x.words[i] ^ y.words[i]
		if words[i] != 0 {
			firstPositive = i
		}
	}
	for i = min; i < max; i++ {
		words[i] = maxi.words[i]
		if words[i] != 0 {
			firstPositive = i
		}
	}
	length := firstPositive + 1
	return Int1024{
		words:  words,
		length: length,
	}
}

// NOT returns ~x
func (x *Int1024) NOT() Int1024 {
	z := Zero()
	var i uint16
	for i = 0; i < x.length; i++ {
		z.words[i] = ^x.words[i]
	}
	return z
}

// BitLength returns the number bits required to represent x (equivalent to len(x.ToBinary()))
func (x *Int1024) BitLength() int {
	if x.length == 1 && x.words[0] == 0 {
		return 1
	}
	return (int(x.length-1))*64 + bits.Len64(x.words[x.length-1])
}

func (x *Int1024) Mask(n uint) Int1024 {
	z := x.Clone()
	z.MaskInPlace(n)
	return z
}

func (x *Int1024) MaskInPlace(n uint) {
	nn := uint16(n / WORDSIZE)
	if n%WORDSIZE != 0 {
		panic("not implemented")
	} else {
		if x.length > nn {
			for i := nn; i < x.length; i++ {
				x.words[i] = 0
			}
			firstPositive := uint16(0)
			for i := nn - 1; i > 0; i-- {
				if x.words[i] != 0 {
					firstPositive = uint16(i)
					break
				}
			}
			x.length = firstPositive + 1
		}
	}
}
