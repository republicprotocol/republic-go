package testutils

import (
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/oracle"
)

// RandMidpointPrice returns a random MidpointPrice.
func RandMidpointPrice() oracle.MidpointPrice {
	tokenPairs, prices := make([]uint64, 10), make([]uint64, 10)
	for i := range tokenPairs {
		tokenPairs[i] = rand.Uint64()
		prices[i] = rand.Uint64()
	}

	return oracle.MidpointPrice{
		Signature:  []byte{},
		TokenPairs: tokenPairs,
		Prices:     prices,
		Nonce:      uint64(time.Now().Unix()),
	}
}
