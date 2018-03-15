package stackint

import (
	"io"
)

// Random returns a new
func Random(rand io.Reader, max *Int1024) (n Int1024, err error) {
	if max.IsZero() {
		return Zero(), nil
	}
	n = max.Sub(&one)
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := n.BitLength()
	if bitLen == 0 {
		// the only valid result is 0
		return Zero(), nil
	}
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, k)

	for {
		_, err = io.ReadFull(rand, bytes)
		if err != nil {
			return Zero(), err
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)

		n = FromBytes(bytes)
		if n.LessThan(max) {
			return
		}
	}
}
