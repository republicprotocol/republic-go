package darknode

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"

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

	id           identity.ID
	address      identity.Address
	multiAddress identity.MultiAddress

	darknodeRegistry contracts.DarkNodeRegistry
	relayer          rpc.RelayService
	smpcer           rpc.ComputerService
}

// NewDarknode returns a new Darknode.
func NewDarknode(config Config) (Darknode, error) {
	node := Darknode{
		Config: config,
		Logger: logger.StdoutLogger,
	}

	// Get identity information from the Config
	key, err := identity.NewKeyPairFromPrivateKey(node.Config.Key.PrivateKey)
	if err != nil {
		return node, fmt.Errorf("cannot get ID from private key: %v", err)
	}
	node.id = key.ID()
	node.address = key.Address()
	node.multiAddress = config.Network.MultiAddress

	// Open a connection to the Ethereum network
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client, err := client.Connect(
		config.Ethereum.URI,
		config.Ethereum.Network,
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		println("ERR!", err)
		return node, err
	}

	// Create bindings to the DarknodeRegistry and Ocean
	darknodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return Darknode{}, err
	}
	node.darknodeRegistry = darknodeRegistry
	node.ocean = NewOcean(darknodeRegistry)

	// Create a channel for notifying the Darknode about new epochs
	node.epochRoutes = make(chan EpochRoute, 2)
	node.orderFragments = make(chan order.Fragment)
	node.deltaFragments = make(chan smpc.DeltaFragment)

	return node, nil
}

// ServeRPC services, using a gRPC server, until the done channel is closed or
// an error is encountered.
func (node *Darknode) ServeRPC(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		server := grpc.NewServer()
		go func() {
			defer server.Stop()
			<-done
		}()

		node.relayer = rpc.NewRelayService(rpc.Options{}, node, node.Logger)
		node.relayer.Register(server)

		node.smpcer = rpc.NewComputerService()
		node.smpcer.Register(server)

		listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", node.Config.Host, node.Config.Port))
		if err != nil {
			node.Logger.Network(logger.Error, err.Error())
			errs <- err
			return
		}

		node.Logger.Network(logger.Info, fmt.Sprintf("listening on %v:%v", node.Config.Host, node.Config.Port))
		if err := server.Serve(listener); err != nil {
			node.Logger.Network(logger.Error, err.Error())
			errs <- err
			return
		}
	}()

	time.Sleep(2 * time.Second)
	return errs
}

// RunWatcher will watch for changes to the Ocean and run the secure
// multi-party computation with new Pools. Stops when the done channel is
// closed, and will attempt to recover from errors encountered while
// interacting with the Ocean.
func (node *Darknode) RunWatcher(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		err := node.darknodeRegistry.WaitUntilRegistration(node.ID())
		if err != nil {
			errs <- err
			return
		}

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

			// Start watching epochs
			epochs, errs := RunEpochWatcher(done)
			for quit := false; !quit; {
				select {

				case <-done:
					return

				case err, ok := <-errs:
					if !ok {
						quit = true
						break
					}
					node.Logger.Network(logger.Error, err.Error())

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

					darknodeIDs, err := node.darknodeRegistry.GetAllNodes()
					if err != nil {
						// FIXME: Do not skip the epoch. Retry with a backoff.
						errs <- err
						continue
					}
					darkOcean, err := NewDarkOcean(epoch, darknodeIDs)
					if err != nil {
						// FIXME: Do not skip the epoch. Retry with a backoff.
						errs <- err
						continue
					}

					deltas, errs := RunEpochProcess(currDone, node.ID(), darkOcean, node.router)
					// TODO: Do something with the smpc.Deltas
				}
			}
		}
	}()

	return errs
}

// OrderMatchToHyperdrive converts an order match into a hyperdrive.Tx and
// forwards it to the Hyperdrive.
func (node *Darknode) OrderMatchToHyperdrive(delta smpc.Delta) {
	if !delta.IsMatch(smpc.Prime) {
		return
	}
	// TODO: Implement
}

// OnOpenOrder implements the rpc.RelayDelegate interface.
func (node *Darknode) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
}

// Ocean returns the Ocean used by this Darknode for computing the Pools and
// its position in them.
func (node *Darknode) Ocean() Ocean {
	return node.ocean
}

// ID returns the ID of the Darknode.
func (node *Darknode) ID() identity.ID {
	return node.id
}

// Address returns the Address of the Darknode.
func (node *Darknode) Address() identity.Address {
	return node.address
}

// MultiAddress returns the MultiAddress of the Darknode.
func (node *Darknode) MultiAddress() identity.MultiAddress {
	return node.multiAddress
}
