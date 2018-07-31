package testutils

import (
	"context"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

type MockMultiAddr struct {
	multiAddr identity.MultiAddress
}

// Swarmer is a mock implementation of the swarm.Swarmer interface.
type Swarmer struct {
	multiAddr    identity.MultiAddress
	multiAddrsMu *sync.Mutex
	multiAddrs   map[identity.Address]MockMultiAddr
}

func NewMockSwarmer(multiAddr identity.MultiAddress) Swarmer {
	return Swarmer{
		multiAddr:    multiAddr,
		multiAddrsMu: new(sync.Mutex),
		multiAddrs:   make(map[identity.Address]MockMultiAddr),
	}
}

func (swarmer *Swarmer) Ping(ctx context.Context) error {
	return nil
}

func (swarmer *Swarmer) Pong(ctx context.Context, to identity.MultiAddress) error {
	return nil
}

func (swarmer *Swarmer) BroadcastMultiAddress(ctx context.Context, multiAddr identity.MultiAddress) error {
	return nil
}

func (swarmer *Swarmer) Query(ctx context.Context, query identity.Address) (identity.MultiAddress, error) {
	return identity.MultiAddress{}, nil
}

func (swarmer *Swarmer) MultiAddress() identity.MultiAddress {
	return swarmer.multiAddr
}

func (swarmer *Swarmer) Peers() (identity.MultiAddresses, error) {
	return make([]identity.MultiAddress, len(swarmer.multiAddrs)), nil
}

func (swarmer *Swarmer) PutMultiAddress(multiAddr identity.MultiAddress) {
	swarmer.multiAddrsMu.Lock()
	defer swarmer.multiAddrsMu.Unlock()
	swarmer.multiAddrs[multiAddr.Address()] = MockMultiAddr{
		multiAddr: multiAddr,
	}
}

func (swarmer *Swarmer) RemoveMultiAddress(multiAddr identity.MultiAddress) {
	swarmer.multiAddrsMu.Lock()
	defer swarmer.multiAddrsMu.Unlock()
	if _, ok := swarmer.multiAddrs[multiAddr.Address()]; ok {
		delete(swarmer.multiAddrs, multiAddr.Address())
	}
}
