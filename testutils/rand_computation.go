package testutils

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
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
	buyFragments[0].EpochDepth = order.FragmentEpochDepth(0)
	sellFragments[0].EpochDepth = order.FragmentEpochDepth(0)
	comp := ome.Computation{
		Buy:        buyFragments[0],
		Sell:       sellFragments[0],
		EpochDepth: order.FragmentEpochDepth(0),
	}
	copy(comp.ID[:], crypto.Keccak256(buy.ID[:], sell.ID[:]))
	return comp, nil
}
