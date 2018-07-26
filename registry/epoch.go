package registry

import (
	"bytes"
	"math/big"

	"github.com/republicprotocol/republic-go/identity"
)

// An Epoch represents the state of an epoch in the Pod. It stores the
// epoch hash, an ordered list of Pods for the epoch, and all Darknode
// identity.Addresses that are registered for the epoch.
type Epoch struct {
	Hash      [32]byte
	Pods      PodHeap
	Darknodes identity.Addresses

	BlockNumber   *big.Int
	BlockInterval *big.Int
}

// Equal returns true if the hash of two Epochs is equal. Otherwise it returns
// false.
func (ξ *Epoch) Equal(arg *Epoch) bool {
	return bytes.Equal(ξ.Hash[:], arg.Hash[:])
}

// Pod returns the Pod that contains the given identity.Address in this Epoch.
// It returns ErrPodNotFound if the identity.Address is not registered in this
// Epoch.
func (ξ *Epoch) Pod(addr identity.Address) (Pod, error) {
	for _, pod := range ξ.Pods {
		for _, darknodeAddr := range pod.Darknodes {
			if addr == darknodeAddr {
				return pod, nil
			}
		}
	}
	return Pod{}, ErrPodNotFound
}
