package testutils

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
)

func ComputationID(buy, sell order.ID) [32]byte {
	var id [32]byte
	copy(id[:], crypto.Keccak256(buy[:], sell[:]))
	return id
}
