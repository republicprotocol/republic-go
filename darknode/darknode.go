package darknode

import (
	"context"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/identity"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

type Darknodes []Darknode

type Darknode struct {
	Config
	Logger *logger.Logger

	darknodeRegistry contracts.DarkNodeRegistry
	ocean            Ocean
}

// NewDarknode returns a new Darknode.
func NewDarknode(config Config) (Darknode, error) {
	node := Darknode{
		Config: config,
		Logger: logger.StdoutLogger,
	}

	// Open a connection to the Ethereum network
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client, err := client.Connect(
		config.Ethereum.URI,
		client.Network(config.Ethereum.Network),
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		return Darknode{}, err
	}

	// Create bindings to the DarknodeRegistry and Ocean
	darknodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return Darknode{}, err
	}
	node.darknodeRegistry = darknodeRegistry
	node.ocean = NewOcean(darknodeRegistry)

	return node, nil
}

// Teardown destroys all resources allocated by the Darknode. The Darknode must
// not be used afterwards.
func (node *Darknode) Teardown() {
	close(node.orderFragments)
	close(node.deltaFragments)
}

// ID returns the ID of the Darknode.
func (node *Darknode) ID() identity.ID {
	key, err := identity.NewKeyPairFromPrivateKey(node.Config.Key.PrivateKey)
	if err != nil {
		panic(err)
	}
	return key.ID()
}

// Ocean returns the Ocean used by this Darknode for computing the Pools and
// its position in them.
func (node *Darknode) Ocean() Ocean {
	return node.ocean
}

// WatchDarknodeRegistryEpoch for changes and connect to Pools as needed.
func (node *Darknode) WatchDarknodeRegistryEpoch(done <-chan struct{}) {

	// Maintain multiple done channels so that multiple epochs can be running
	// in parallel
	var prevDone chan struct{}
	var currDone chan struct{}
	defer func() {
		if prevDone != nil {
			close(prevDone)
		}
		if currDone != nil {
			close(currDone)
		}
	}()

	// Looping until the done channel is closed will recover from errors
	// returned by watching the Ocean
	for {
		select {
		case <-done:
			return
		default:
		}

		epochs, errs := node.ocean.Watch(done)

		for quit := false; !quit; {
			select {

			case <-done:
				return

			case err, ok := <-errs:
				if !ok {
					quit = true
					break
				}
				node.Logger.Network(logger.Error, err.String())

			case epoch, ok := <-epochs:
				if !ok {
					quit = true
					break
				}

				if prevDone != nil {
					close(prevDone)
				}
				prevDone = currDone
				currDone = make(chan struct{})
				node.RunEpoch(currDone, epoch)
			}
		}
	}
}

// RunEpoch until the done channel is closed. Order fragments will be received,
// computed, and broadcast in accordance to the epoch.
func (node *Darknode) RunEpoch(done <-chan struct{}, epoch contracts.Epoch) {

	// Setup a secure multi-party computer
	// FIXME: Calculate n-k threshold correctly
	n := int64(5)
	k := (n + 1) * 2 / 3

	smpcerID := smpc.ComputerID{}
	copy(smpcerID[:], node.ID()[:])
	smpcer = smpc.NewComputer(smpcerID, n, k)

	// Initialize channels
	orderFragments = make(chan order.Fragment)
	deltaFragments = make(chan smpc.DeltaFragment)

}

// Compute begins the Smpcer.
func (node *Darknode) Compute(done <-chan struct{}) {
	deltaFragmentsComputed, deltasComputed := node.smpcer.ComputeOrderMatches(done, node.orderFragments, node.deltaFragments)

	go func() {
		for delta := range deltasComputed {
			if delta.IsMatch(&smpc.Prime) {
				node.Logger.OrderMatch(logger.Info, delta.ID.String(), delta.BuyOrderID.String(), delta.SellOrderID.String())
				node.OrderMatchToHyperdrive(delta)
			}
		}
	}()
}

func (node *Darknode) OrderMatchToHyperdrive(delta smpc.Delta) {
	if !delta.IsMatch(&smpc.Prime) {
		return
	}

	// TODO:
	// 1. Create a Tx for Hyperdrive
	// 2. Gossip the Tx to the Hyperdrive replicas and eventually it will reach
	//    the proposer replica
}
