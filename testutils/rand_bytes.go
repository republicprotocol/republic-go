package testutils

import (
	"fmt"
	"math/rand"

	"github.com/republicprotocol/republic-go/crypto"
)

// Random32Bytes creates a random [32]byte.
func Random32Bytes() [32]byte {
	var res [32]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(res[:], hash)
	return res
}

// Random64Bytes creates a random [64]byte.
func Random64Bytes() [64]byte {
	var res [64]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(res[:], hash)
	return res
}
