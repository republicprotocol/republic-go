package stackint

// ShiftLeft returns x<<n
func (x *Int1024) ShiftLeft(n uint) Int1024 {
	z := x.Clone()
	z.ShiftLeftInPlace(n)
	return z
}

// ShiftLeftInPlace shifts to the left x by one
func (x *Int1024) ShiftLeftInPlace(n uint) {

	// If n > 64, first, shift entire words
	div := n / WORDSIZE
	if div > 0 {
		x.length = min(SIZE, x.length+uint16(div))
		var i uint
		for i = uint(x.length); i >= div; i-- {
			x.words[i] = x.words[i-div]
		}
		for i = 0; i < div; i++ {
			x.words[i] = 0
		}
		n = n - div*WORDSIZE
	}

	if n == 1 {
		x.shiftleftone()
	} else {
		x.shiftleft(n)
	}
}

func (x *Int1024) shiftleft(n uint) {
	if n >= SIZE {
		panic("shifting by more than a word")
	}
	var overflow uint64
	var shift uint64 = (1<<n - 1)
	// fmt.Println(shift)
	var i uint16
	for i = 0; i < x.length; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - n)) & shift
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] << n) | overflow
		overflow = newOverflow
	}
	if overflow != 0 && x.length < SIZE {
		x.length++
		x.words[x.length-1] = overflow
	}
	x.length += uint16(n)
}

func (x *Int1024) shiftleftone() {
	var overflow uint64
	var i uint16
	for i = 0; i < x.length; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - 1)) & 1
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] << 1) | overflow
		overflow = newOverflow
	}
	if overflow != 0 && x.length < SIZE {
		x.length++
		x.words[x.length-1] = overflow
	}
	x.length++
}

// ShiftRight returns x>>n
func (x *Int1024) ShiftRight(n uint) Int1024 {
	z := x.Clone()
	z.ShiftRightInPlace(n)
	return z
}

// ShiftRightInPlace shifts to the right x by one
func (x *Int1024) ShiftRightInPlace(n uint) {
	// If n > 64, first, shift entire words
	div := n / WORDSIZE
	if div > 0 {
		x.length = max(0, x.length-uint16(div))
		var i uint
		for i = 0; i < uint(x.length); i++ {
			x.words[i] = x.words[i+div]
		}
		for i = uint(x.length); i < uint(x.length)+div; i++ {
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
}

func (x *Int1024) shiftright(n uint) {
	if n >= SIZE {
		panic("shifting by more than a word")
	}
	var overflow uint64
	var shift uint64 = (1<<n - 1)
	for i := x.length - 1; i >= 0; i-- {
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
	for i := x.length - 1; i >= 0; i-- {
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
	var i uint16
	for i = 0; i < min; i++ {
		z.words[i] = x.words[i] & y.words[i]
	}
	return z
}

// OR returns x|y
func (x *Int1024) OR(y *Int1024) Int1024 {
	min := x.length
	max := y.length
	maxi := y
	if max < min {
		min, max = max, min
		maxi = x
	}

	z := Zero()
	var i uint16
	for i = 0; i < min; i++ {
		z.words[i] = x.words[i] | y.words[i]
	}
	for i = 0; i < max; i++ {
		z.words[i] = maxi.words[i]
	}

	return z
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

	z := Zero()
	var i uint16
	for i = 0; i < min; i++ {
		z.words[i] = x.words[i] ^ y.words[i]
	}
	for i = 0; i < max; i++ {
		z.words[i] = maxi.words[i]
	}
	return z
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
	word := x.words[x.length-1]
	wordBits := 0
	for word > 0 {
		word /= 2
		wordBits++
	}
	if wordBits == 0 {
		wordBits = 1
	}
	return (int(x.length-1))*64 + wordBits
}
