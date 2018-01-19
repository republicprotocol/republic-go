package miner

import (
	"log"
	"math/big"
	"runtime"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order-compute"
)

// TODO: Do not make this values constant.
var (
	N        = int64(3)
	K        = int64(2)
	Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
)

type Miner struct {
	*network.Node
	*compute.ComputationMatrix
}

func NewMiner(config *Config) (*Miner, error) {
	miner := &Miner{
		ComputationMatrix: compute.NewComputationMatrix(),
	}
	node, err := network.NewNode(config.Multi, config.BootstrapMultis, miner)
	if err != nil {
		return nil, err
	}
	miner.Node = node
	return miner, nil
}

func (miner *Miner) EstablishConnections() error {
	for _, multi := range miner.DHT.MultiAddresses() {
		log.Println("pinging", multi)
		_, err := miner.RPCPing(multi)
		if err != nil {
			return err
		}
	}
	return nil
}

func (miner *Miner) OnPingReceived(peer identity.MultiAddress) {
}

func (miner *Miner) OnOrderFragmentReceived(orderFragment *compute.OrderFragment) {
	log.Println("received order fragment =", base58.Encode(orderFragment.ID))
	miner.ComputationMatrix.AddOrderFragment(orderFragment)
	log.Println("computation matrix updated")
}

func (miner *Miner) OnResultFragmentReceived(resultFragment *compute.ResultFragment) {
	log.Println("received result fragment =", base58.Encode(resultFragment.ID))
	miner.addResultFragments([]*compute.ResultFragment{resultFragment})
}

func (miner *Miner) Mine(quit chan struct{}) {
	for {
		select {
		case <-quit:
			miner.Stop()
			return
		default:
			// FIXME: If this function call blocks forever then the quit signal
			// will never be received.
			miner.ComputeAll()
		}
	}
}

func (miner Miner) ComputeAll() {
	log.Println("waiting for new computations...")
	numberOfCPUs := runtime.NumCPU()
	computations := miner.ComputationMatrix.WaitForComputations(numberOfCPUs)
	resultFragments := make([]*compute.ResultFragment, len(computations))

	log.Println("executing new computations...")
	do.CoForAll(computations, func(i int) {
		resultFragment, err := miner.Compute(computations[i])
		if err != nil {
			return
		}
		resultFragments[i] = resultFragment
	})
	log.Println("computations done.")
	go func() {
		resultFragmentsOk := make([]*compute.ResultFragment, 0, len(resultFragments))
		for _, resultFragment := range resultFragments {
			if resultFragment != nil {
				resultFragmentsOk = append(resultFragmentsOk, resultFragment)
			}
		}
		miner.addResultFragments(resultFragmentsOk)
	}()
}

// Compute the required computation on two OrderFragments and send the result
// to all Miners in the M Network.
// TODO: Send computed order fragments to the M Network instead of all peers.
func (miner Miner) Compute(computation *compute.Computation) (*compute.ResultFragment, error) {
	resultFragment, err := computation.Sub(Prime)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, multi := range miner.DHT.MultiAddresses() {
			miner.RPCSendResultFragment(multi, resultFragment)
		}
	}()
	return resultFragment, nil
}

func (miner Miner) addResultFragments(resultFragments []*compute.ResultFragment) {
	results, _ := miner.ComputationMatrix.AddResultFragments(K, Prime, resultFragments)
	for _, result := range results {
		if result.IsMatch() {
			log.Println("match found for buy =", base58.Encode(result.BuyOrderID), ", sell =", base58.Encode(result.SellOrderID))
		}
	}
}
