package stackint

import (
	"io"
)

// Random returns a new int1024 less than max (or equal if max is 0), filled from the io reader
func Random(rand io.Reader, max *Int1024) (Int1024, error) {
	// If max is 0 or 1, then only option is 0
	if max.LessThanOrEqual(&one) {
		return Zero(), nil
	}
	n := max.Sub(&one)
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := n.BitLength()
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, k)

	for {
		_, err := io.ReadFull(rand, bytes)
		if err != nil {
			return Int1024{}, err
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)

		n, err = FromBytes(bytes)
		if err != nil {
			return Int1024{}, err
		}
		if n.LessThan(max) {
			return n, nil
		}
	}
}
