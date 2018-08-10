package testutils

import (
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/oracle"
)

// RandMidpointPrice returns a random MidpointPrice.
func RandMidpointPrice() oracle.MidpointPrice {
	prices := make(map[uint64]uint64, 10)
	for i := range prices {
		prices[i] = rand.Uint64()
	}

	return oracle.MidpointPrice{
		Signature: []byte{},
		Prices:    prices,
		Nonce:     uint64(time.Now().Unix()),
	}
}
