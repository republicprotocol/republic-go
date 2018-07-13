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
