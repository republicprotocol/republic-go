package testutils

import (
	"math/big"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/registry"
)

// RandomEpoch returns a new epoch with only one pod and one darknode.
func RandomEpoch(i int) (identity.Address, registry.Epoch, error) {
	addr, err := RandomAddress()
	if err != nil {
		return identity.Address(""), registry.Epoch{}, err
	}
	return addr, registry.Epoch{
		Hash: Random32Bytes(),
		Pods: []registry.Pod{
			{
				Position:  0,
				Hash:      Random32Bytes(),
				Darknodes: []identity.Address{addr},
			},
		},
		Darknodes:     []identity.Address{addr},
		BlockNumber:   big.NewInt(int64(i)),
		BlockInterval: big.NewInt(int64(2)),
	}, nil
}
