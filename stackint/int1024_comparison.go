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
	return x.EqualsUint64(0)
}

// EqualsUint64 returns true of x represents the uint64 n
func (x *Int1024) EqualsUint64(n uint64) bool {
	for i := 0; i < INT1024WORDS-1; i++ {
		if x.words[i] != 0 {
			return false
		}
	}
	return x.words[INT1024WORDS-1] == n
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

// LessThan returns x<y
func (x *Int1024) LessThan(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] < y.words[i] {
			return true
		}
		if x.words[i] > y.words[i] {
			return false
		}
	}
	return false
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
