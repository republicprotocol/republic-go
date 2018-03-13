package stackint

// Equals returns true of x and y represent the same Int1024
func (x *Int1024) Equals(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] != y.words[i] {
			return false
		}
	}
	return true
}

// IsZero returns true of x == 0
func (x *Int1024) IsZero() bool {
	return x.Equals(&ZERO)
}

// LessThan returns x<y
func (x *Int1024) LessThan(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] < y.words[i] {
			return true
		}
	}
	return false
}

// LessThanOrEqual returns x<=y
func (x *Int1024) LessThanOrEqual(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] > y.words[i] {
			return false
		}
	}
	return true
}

// GreaterThan returns x>y
func (x *Int1024) GreaterThan(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] > y.words[i] {
			return true
		}
	}
	return false
}

// GreaterThanOrEqual returns x>=y
func (x *Int1024) GreaterThanOrEqual(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] < y.words[i] {
			return false
		}
	}
	return true
}

// IsEven returns (x%2)==0
func (x *Int1024) IsEven() bool {
	return (x.words[INT1024WORDS-1] & 1) == 0
}
