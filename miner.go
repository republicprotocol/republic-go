package miner

import (
	"math/big"
	"sync"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order"
)

type Miner struct {
	network.Node

	orderFragmentsMu          *sync.Mutex
	orderFragments            map[string]*order.Fragment
	orderFragmentComputations map[string]map[string]*order.Fragment

	done chan struct{}
}

func NewMiner(config Config) Miner {
	miner := Miner{}
	miner.Node = network.NewNode(config.MultiAddress, config.MultiAddresses, miner)
	return miner
}

func (miner Miner) Start() {
	go func() {
		for {
			select {
			case <-miner.done:
				return
			default:
				miner.Mine()
			}
		}
	}()
}

func (miner Miner) Stop() {
	miner.done <- struct{}{}
}

func (miner Miner) Mine() {
	for lhs := range miner.orderFragments {
		for rhs := range miner.orderFragments {
			if lhs == rhs {
				continue
			}
			if result := miner.orderFragmentComputations[lhs][rhs]; result == nil {
				result, err := miner.orderFragments[lhs].Add(miner.orderFragments[rhs], big.NewInt(2))
				if err != nil {
					continue
				}
				miner.orderFragmentComputations[lhs][rhs] = result
			}
		}
	}
}

func (miner Miner) OnPingReceived(peer identity.MultiAddress) {
}

func (miner Miner) OnOrderFragmentReceived(orderFragment order.Fragment) {
	miner.addOrderFragment(orderFragment)
}

func (miner Miner) addOrderFragment(orderFragment order.Fragment) {
	miner.orderFragmentsMu.Lock()
	defer miner.orderFragmentsMu.Unlock()
	miner.orderFragments[string(orderFragment.ID)] = orderFragment
}
