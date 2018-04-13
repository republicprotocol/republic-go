package darknode

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/smpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"
)

type DarkNodes []DarkNode

type DarkNode struct {
	Config

	darkNodeRegistry contracts.DarkNodeRegistry
	darkOcean        *darkocean.Ocean
	Logger           *logger.Logger

	Server *grpc.Server
	Relay  *rpc.RelayService

	Computer smpc.Computer
}

func NewDarkNode(config Config) (DarkNode, error) {
	node := new(DarkNode)
	node.Config = config

	// Connect to Ethereum
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client, err := client.Connect(
		config.Ethereum.URI,
		client.Network(config.Ethereum.Network),
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		return DarkNode{}, err
	}
	darkNodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return DarkNode{}, err
	}
	node.darkNodeRegistry = darkNodeRegistry

	// Create dark ocean.
	darkOcean, err := darkocean.NewOcean(darkNodeRegistry)
	if err != nil {
		return DarkNode{}, err
	}
	node.darkOcean = darkOcean

	// Create the logger and start all plugins
	node.Logger, err = logger.NewLogger(config.LoggerOptions)
	if err != nil {
		return DarkNode{}, err
	}
	node.Logger.Start()

	// Initialize RPC server and services
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	node.Relay = rpc.NewRelayService(node.NetworkOption, node, node.Logger)
	node.Computer = rpc.NewComputerService()

	return *node, nil
}

// Stop the DarkNode.
func (node *DarkNode) Stop() {
	// Stop serving gRPC services
	node.Server.Stop()

	// Stop the logger
	node.Logger.Stop()

	// Force the GC to run
	runtime.GC()
}

func (node *DarkNode) Run(ctx context.Context) {
	errChs := [2]<-chan error{}

	errChs[0] = node.RunRPC(ctx)
	errChs[1] = node.RunDarkOcean(ctx)

	for err := range dispatch.MergeErrors(errChs[:]...) {
		log.Println(err)
	}
}

func (node *DarkNode) RunRPC(ctx context.Context) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		// Turn the gRPC server on.
		node.Logger.Network(logger.Info, fmt.Sprintf("gRPC services listening on %s:%s", node.Host, node.Port))
		node.Relay.Register(node.Server)
		node.Computer.Register(node.Server)

		listener, err := net.Listen("tcp", node.Host+":"+node.Port)
		if err != nil {
			errCh <- err
			return
		}
		if err := node.Server.Serve(listener); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		<-ctx.Done()
		node.Server.Stop()
	}()

	return errCh
}

func (node *DarkNode) RunDarkOcean(ctx context.Context) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		epoch, err := node.darkNodeRegistry.CurrentEpoch()
		if err != nil {
			errCh <- err
			return
		}
		minimumEpochIntervalBig, err := node.darkNodeRegistry.MinimumEpochInterval()
		if err != nil {
			errCh <- err
			return
		}
		minimumEpochInterval, err := minimumEpochIntervalBig.ToUint()
		if err != nil {
			errCh <- err
			return
		}

		t := time.NewTicker(time.Duration(minimumEpochInterval*1000/24) * time.Millisecond)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-t.C:
				nextEpoch, err := node.darkNodeRegistry.CurrentEpoch()
				if err != nil {
					errCh <- err
					continue
				}
				if bytes.Equal(epoch.Blockhash[:], nextEpoch.Blockhash[:]) {
					continue
				}
				epoch = nextEpoch
				if err := node.darkOcean.Update(); err != nil {
					errCh <- fmt.Errorf("cannot update dark ocean: %v", err)
				}
			}
		}
	}()

	return errCh
}

func (node *DarkNode) DarkOcean() *darkocean.Ocean {
	return node.darkOcean
}

func (node *DarkNode) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {

	log.Printf( "order %s received from the %s", orderFragment.OrderID.String(), from.ID().String())
}
