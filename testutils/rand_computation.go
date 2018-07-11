package testutils

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/ome"
)

// RandomComputation generates a random computation with empty epoch hash.
func RandomComputation() (ome.Computation, error) {
	buy, sell := RandomBuyOrder(), RandomSellOrder()
	buyFragments, err := buy.Split(24, 16)
	if err != nil {
		return ome.Computation{}, err
	}
	sellFragments, err := sell.Split(24, 16)
	if err != nil {
		return ome.Computation{}, err
	}
	comp := ome.Computation{
		Buy:  buyFragments[0],
		Sell: sellFragments[0],
	}
	copy(comp.ID[:], crypto.Keccak256(buy.ID[:], sell.ID[:]))
	return comp, nil
}
