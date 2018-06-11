package stackint

import "github.com/republicprotocol/republic-go/stackint/asm"

// IsZero returns true of x == 0
func (x *Int1024) IsZero() bool {
	return x.length == 0
}

// EqualsWord returns true of x represents the Word n
func (x *Int1024) EqualsWord(n asm.Word) bool {
	if n == 0 {
		return x.length == 0
	}
	return x.length == 1 && x.words[0] == n
}

// Cmp returns -1 if x<y, 0 if x=y, 1 if x>y
func (x *Int1024) Cmp(y *Int1024) int {

	if x.length < y.length {
		return -1
	} else if x.length > y.length {
		return 1
	} else {
		var i int16
		for i = int16(x.length) - 1; i >= 0; i-- {
			if x.words[i] < y.words[i] {
				return -1
			}
			if x.words[i] > y.words[i] {
				return +1
			}
		}
		return 0
	}
}

// Equals returns true of x and y represent the same Int1024
func (x *Int1024) Equals(y *Int1024) bool {
	return x.Cmp(y) == 0
}

// LessThan returns x<y
func (x *Int1024) LessThan(y *Int1024) bool {
	return x.Cmp(y) < 0
}

// GreaterThan returns x>y
func (x *Int1024) GreaterThan(y *Int1024) bool {
	return x.Cmp(y) > 0
}

// LessThanOrEqual returns x<=y
func (x *Int1024) LessThanOrEqual(y *Int1024) bool {
	return x.Cmp(y) <= 0
}

// GreaterThanOrEqual returns x>=y
func (x *Int1024) GreaterThanOrEqual(y *Int1024) bool {
	return x.Cmp(y) >= 0
}

// IsEven returns (x%2)==0
func (x *Int1024) IsEven() bool {
	return (x.words[0] & 1) == 0
}
