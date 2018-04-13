package darknode

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/smpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/dht"
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
	DHT              *dht.DHT
	Computer         smpc.Computer
	ClientPool       rpc.ClientPool

	Server          *grpc.Server
	RelayService    *rpc.RelayService
	ComputerService *rpc.ComputerService

	orderFragmentCh chan order.Fragment
	orderTupleCh    <-chan smpc.OrderTuple
	deltaFragmentCh chan smpc.DeltaFragment
	deltaCh         <-chan smpc.Delta

	sharedOrderTable   smpc.SharedOrderTable
	sharedDeltaBuilder smpc.SharedDeltaBuilder

	// FIXME: Improve beyond working for the demo.
	routingTableOut map[string]chan *rpc.Computation
	routingTableIn  map[string]<-chan *rpc.Computation
	routingTableErr map[string]<-chan error
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

	// Create the clientPool
	node.ClientPool = *rpc.NewClientPool(node.NetworkOption.MultiAddress)

	// Create DHT
	node.DHT = dht.NewDHT(node.Config.NetworkOption.MultiAddress.Address(), node.Config.NetworkOption.MaxBucketLength)

	// Create the logger and start all plugins
	node.Logger, err = logger.NewLogger(config.LoggerOptions)
	if err != nil {
		return DarkNode{}, err
	}
	node.Logger.Start()

	// Initialize RPC server and services
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))

	node.RelayService = rpc.NewRelayService(node.NetworkOption, node, node.Logger)
	node.ComputerService = rpc.NewComputerService()

	// smpc
	node.orderFragmentCh = make(chan order.Fragment)
	node.deltaFragmentCh = make(chan smpc.DeltaFragment)

	node.sharedOrderTable = smpc.NewSharedOrderTable()
	node.sharedDeltaBuilder = smpc.NewSharedDeltaBuilder(int64(len(node.NetworkOption.BootstrapMultiAddresses)), smpc.Prime)
	node.routingTableOut = map[string]chan *rpc.Computation{}
	node.routingTableIn = map[string]<-chan *rpc.Computation{}
	node.routingTableErr = map[string]<-chan error{}

	for _, multiAddress := range node.NetworkOption.BootstrapMultiAddresses {
		if bytes.Equal(node.NetworkOption.Address.ID(), multiAddress.ID()) {
			continue
		}
		node.routingTableOut[multiAddress.String()] = make(chan *rpc.Computation)
	}

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

	errChs := [5]<-chan error{}

	// FIXME: Integrate correctly.
	node.orderTupleCh, errChs[2] = smpc.ProcessOrderFragments(ctx, node.orderFragmentCh, &node.sharedOrderTable, 1)
	deltaFragmentCh, _ := smpc.ProcessOrderTuples(ctx, node.orderTupleCh, 1)
	go func() {
		for deltaFragment := range deltaFragmentCh {
			node.deltaFragmentCh <- deltaFragment
		}
	}()

	node.deltaCh, errChs[4] = smpc.ProcessDeltaFragments(ctx, node.deltaFragmentCh, &node.sharedDeltaBuilder, 1)

	errChs[0] = node.RunRPC(ctx)
	errChs[1] = node.RunDarkOcean(ctx)

	go func() {
		for _, multiAddress := range node.NetworkOption.BootstrapMultiAddresses {
			go func(multiAddress identity.MultiAddress) {
				if bytes.Equal(node.NetworkOption.Address.ID(), multiAddress.ID()) {
					return
				}
				if bytes.Compare(node.NetworkOption.Address.ID(), multiAddress.ID()) < 0 {
					node.routingTableIn[multiAddress.String()], node.routingTableErr[multiAddress.String()] = node.ClientPool.Compute(ctx, multiAddress, node.routingTableOut[multiAddress.String()])
					log.Println("Connected")

					for computation := range node.routingTableIn[multiAddress.String()] {
						log.Println("Receiving delta fragment")
						deltaFragment, err := rpc.UnmarshalDeltaFragment(computation.DeltaFragment)
						panic(err)
						node.deltaFragmentCh <- deltaFragment
					}
				} else {
					node.routingTableIn[multiAddress.String()], node.routingTableErr[multiAddress.String()] = node.ComputerService.WaitForCompute(multiAddress, node.routingTableOut[multiAddress.String()])
					log.Println("Connected")
					for computation := range node.routingTableIn[multiAddress.String()] {
						log.Println("Receiving delta fragment")
						deltaFragment, err := rpc.UnmarshalDeltaFragment(computation.DeltaFragment)
						panic(err)
						node.deltaFragmentCh <- deltaFragment
					}
				}
			}(multiAddress)
		}
	}()

	go func() {
		for deltaFragment := range node.deltaFragmentCh {
			var wg sync.WaitGroup
			wg.Add(len(node.NetworkOption.BootstrapMultiAddresses))
			for _, multiAddress := range node.NetworkOption.BootstrapMultiAddresses {
				log.Println("Broadcasting delta fragment")
				go func(multiAddress identity.MultiAddress) {
					defer wg.Done()
					node.routingTableOut[multiAddress.String()] <- &rpc.Computation{DeltaFragment: rpc.MarshalDeltaFragment(&deltaFragment)}
				}(multiAddress)
			}
			wg.Wait()
		}
	}()

	go func() {
		for delta := range node.deltaCh {
			if delta.IsMatch(smpc.Prime) {
				log.Printf("match found: buy = %s; sell = %s", delta.BuyOrderID.String(), delta.SellOrderID.String())
			}
		}
	}()

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

		node.RelayService.Register(node.Server)
		node.ComputerService.Register(node.Server)

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
	go func() { node.orderFragmentCh <- *orderFragment }()
}
