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

// Cmp returns -1 if x<y, 0 if x=y, 1 if x>y
func (x *Int1024) Cmp(y *Int1024) int {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] < y.words[i] {
			return -1
		}
		if x.words[i] > y.words[i] {
			return +1
		}
	}
	return 0
}

// GreaterThan returns x>y
func (x *Int1024) GreaterThan(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] > y.words[i] {
			return true
		}
		if x.words[i] < y.words[i] {
			return false
		}
	}
	return false
}

// LessThanOrEqual returns x<=y
func (x *Int1024) LessThanOrEqual(y *Int1024) bool {
	return !x.GreaterThan(y)
}

// GreaterThanOrEqual returns x>=y
func (x *Int1024) GreaterThanOrEqual(y *Int1024) bool {
	return !x.LessThan(y)
}

// IsEven returns (x%2)==0
func (x *Int1024) IsEven() bool {
	return (x.words[INT1024WORDS-1] & 1) == 0
}
