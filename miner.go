package miner

import (
	"log"
	"runtime"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order-compute"
)

type Miner struct {
	network.Node
	compute.ComputationMatrix
	compute.ReconstructionMatrix
}

func NewMiner(config Config) Miner {
	miner := Miner{}
	miner.Node = network.NewNode(config.MultiAddress, config.MultiAddresses, miner)
	return miner
}

func (miner Miner) OnPingReceived(peer identity.MultiAddress) {
}

func (miner Miner) OnOrderFragmentReceived(orderFragment compute.OrderFragment) {
	miner.ComputationMatrix.FillComputations(&orderFragment)
}

func (miner Miner) OnComputedOrderFragmentReceived(orderFragment compute.OrderFragment) {
}

func (miner Miner) Mine(quit chan struct{}) {
	go func() {
		if err := miner.Serve(); err != nil {
			// TODO: Do something other than die.
			log.Fatal(err)
		}
	}()
	for {
		select {
		case <-quit:
			miner.Stop()
			return
		default:
			do.CoBegin(
				func() do.Option {
					miner.ComputeAll()
					return do.Ok(nil)
				},
				func() do.Option {
					miner.ReconstructAll()
					return do.Ok(nil)
				},
			)
		}
	}
}

func (miner Miner) ComputeAll() {
	numberOfCPUs := runtime.NumCPU()
	computations := miner.ComputationMatrix.WaitForComputations(numberOfCPUs)
	do.CoForAll(computations, func(i int) {
		miner.Compute(computations[i])
	})
}

// Compute the required computation on two OrderFragments and send the result
// to all Miners in the M Network.
// TODO: Send computed order fragments to the M Network instead of all peers.
func (miner Miner) Compute(com *Computation) {
	com.Sub()
	miner.JoinMatrix.FillJoins(com.Out)
	for _, multi := range miner.DHT.MultiAddresses() {
		network.RPCSendComputedOrderFragment(multi, com.Out)
	}
}

func (miner Miner) ReconstructAll() {
	numberOfCPUs := runtime.NumCPU()
	reconstructables := miner.ReconstructionMatrix.WaitForJoins(numberOfCPUs)
	do.CoForAll(joins, func(i int) {
		miner.Reconstruct(reconstructables[i])
	})
}

func (miner Miner) Reconstruct(orderFragments []*compute.OrderFragment) {
	match, err := compute.IsMatch(orderFragments)
	if err != nil {
		return
	}
	if match {
		// TODO: Do something other than logging to Stdout.
		log.Println("match", orderFragments[0].OrderIDs)
	}
}
