package stackint

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
