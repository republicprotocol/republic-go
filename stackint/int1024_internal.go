package stackint

// overwritingAdd sets x to x+y (used for multiplication)
func (x *Int1024) overwritingAdd(y *Int1024) {
	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		previousOverflow := overflow
		if x.words[i] > WORDMAX-y.words[i]-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = x.words[i] + y.words[i] + previousOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}
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

// shiftRightByOne shifts to the right x by one
func (x *Int1024) overwritingShiftRightByOne() {
	overflow := Word(0)
	for i := 0; i < INT1024WORDS; i++ {
		// Calculate if word overflows into next word
		newOverflow := (x.words[i] >> (WORDSIZE - 1)) & (1 << (WORDSIZE - 1))
		// Shift word to the right
		// If previous word overflowed, add 1
		x.words[i] = (x.words[i] >> 1) | overflow
		overflow = newOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured (not important for Shift)
	}
}
