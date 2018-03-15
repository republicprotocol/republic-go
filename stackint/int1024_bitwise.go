package stackint

// ShiftLeft returns x<<n
func (x *Int1024) ShiftLeft(n int) Int1024 {
	z := x.Clone()
	for i := 0; i < n; i++ {
		z.overwritingShiftLeftByOne()
	}
	return z
}

// shiftLeftByOne shifts to the left x by one
func (x *Int1024) overwritingShiftLeftByOne() {
	overflow := Word(0)
	for i := INT1024WORDS - 1; i >= 0; i-- {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - 1)) & 1
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] << 1) | overflow
		overflow = newOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured (not important for Shift)
	}
}

// ShiftRight returns x>>n
func (x *Int1024) ShiftRight(n int) Int1024 {
	z := x.Clone()
	for i := 0; i < n; i++ {
		z.overwritingShiftRightByOne()
	}
	return z
}

// shiftRightByOne shifts to the right x by one
func (x *Int1024) overwritingShiftRightByOne() {
	overflow := Word(0)
	for i := 0; i < INT1024WORDS; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] & 1) << (WORDSIZE - 1)
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] >> 1) | overflow
		overflow = newOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured (not important for Shift)
	}
}

// AND returns x&y
func (x *Int1024) AND(y *Int1024) Int1024 {
	z := zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = x.words[i] & y.words[i]
	}
	return z
}

// OR returns x|y
func (x *Int1024) OR(y *Int1024) Int1024 {
	z := zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = x.words[i] | y.words[i]
	}
	return z
}

// XOR returns x&y
func (x *Int1024) XOR(y *Int1024) Int1024 {
	z := zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = x.words[i] ^ y.words[i]
	}
	return z
}

// NOT returns ~x
func (x *Int1024) NOT() Int1024 {
	z := zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = ^x.words[i]
	}
	return z
}
