package stackint

// ShiftLeft returns x<<n
func (x *Int1024) ShiftLeft(n int) Int1024 {
	z := x.Clone()
	for i := 0; i < n; i++ {
		z.overwritingShiftLeftByOne()
	}
	return z
}

// ShiftRight returns x>>n
func (x *Int1024) ShiftRight(n int) Int1024 {
	z := x.Clone()
	for i := 0; i < n; i++ {
		z.overwritingShiftRightByOne()
	}
	return z
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
