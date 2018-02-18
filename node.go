package node

import (
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-swarm-network"
	"google.golang.org/grpc"
)

type DarkNode struct {
	Server        *grpc.Server
	Swarm         *swarm.Node
	Dark          *dark.Node
	Configuration *Config

	HiddenOrderBook *compute.HiddenOrderBook
	DarkPool        MultiAddresses

	quitServer chan struct{}
	quitPacker chan struct{}
}

// NewDarkNode creates a new DarkNode, a new swarm.Node and dark.Node and assigns the
// new DarkNode as the delegate for both. Returns the new DarkNode, or an error.
func NewDarkNode(config *Config) (*DarkNode, error) {
	node := &DarkNode{
		HiddenOrderBook: compute.NewHiddenOrderBook(4),
		Server:          grpc.NewServer(grpc.ConnectionTimeout(time.Minute)),
		Configuration:   config,

		quitServer: make(chan struct{}),
		quitPacker: make(chan struct{}),
	}

	swarmOptions := swarm.Options{
		MultiAddress:            config.MultiAddress,
		BootstrapMultiAddresses: config.BootstrapMultiAddresses,
		Debug:           swarm.DebugHigh,
		Alpha:           3,
		MaxBucketLength: 20,
		Timeout:         30 * time.Second,
		TimeoutStep:     30 * time.Second,
		TimeoutRetries:  3,
		Concurrent:      true,
	}
	swarmNode := swarm.NewNode(node.Server, node, swarmOptions)
	node.Swarm = swarmNode

	darkOptions := dark.Options{
		Address:        config.MultiAddress.Address(),
		Debug:          dark.DebugHigh,
		Timeout:        30 * time.Second,
		TimeoutStep:    30 * time.Second,
		TimeoutRetries: 3,
		Concurrent:     true,
	}
	darkNode := dark.NewNode(node.Server, node, darkOptions)
	node.Dark = darkNode

	return node, nil
}

// Start mining for compute.Orders that are matched. It establishes connections
// to other peers in the swarm network by bootstrapping against a set of
// bootstrap swarm.Nodes.
func (node *DarkNode) Start() {
	// Start both gRPC servers.
	go func() {
		log.Printf("Listening on %s:%s\n", miner.Configuration.Host, miner.Configuration.Port)
		miner.Swarm.Register()
		miner.Xing.Register()
		listener, err := net.Listen("tcp", miner.Configuration.Host+":"+miner.Configuration.Port)
		if err != nil {
			log.Fatal(err)
		}
		if err := miner.Server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the server to start and bootstrap the connections in the swarm.
	time.Sleep(time.Second)
	miner.Swarm.Bootstrap()

	// Wait until the signal for stopping the server is received, and then call
	// Stop.
	<-node.quitServer
	node.Server.Stop()
}

// Stop mining.
func (node *DarkNode) Stop() {
	node.quitServer <- struct{}{}
	node.quitPacker <- struct{}{}
}

// OnPingReceived implements the swarm.Delegate interface. It is used by the
// underlying swarm.Node whenever the Miner has handled a Ping RPC.
func (miner *Miner) OnPingReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersReceived implements the swarm.Delegate interface. It is
// used by the underlying swarm.Node whenever the Miner has handled a
// QueryCloserPeers RPC.
func (miner *Miner) OnQueryCloserPeersReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersOnFrontierReceived implements the swarm.Delegate
// interface. It is called by the underlying swarm.Node whenever the Miner
// has handled a QueryCloserPeersOnFrontier RPC.
func (miner *Miner) OnQueryCloserPeersOnFrontierReceived(peer identity.MultiAddress) {
}

// OnOrderFragmentReceived implements the xing.Delegate interface. It is called
// by the underlying xing.Node whenever the Miner receives a
// compute.OrderFragment that it must process.
func (miner *Miner) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	miner.HiddenOrderBook.AddOrderFragment(orderFragment)
}

// OnResultFragmentReceived implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.ResultFragment that it must process.
func (miner *Miner) OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment) {
	miner.addResultFragments([]*compute.ResultFragment{resultFragment})
}

// OnOrderFragmentForwarding implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.OrderFragment that it must forward to the correct xing.Node.
func (miner *Miner) OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment) {
}

// OnResultFragmentForwarding implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.ResultFragment that it must forward to the correct xing.Node.
func (miner *Miner) OnResultFragmentForwarding(to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment) {
}

func (node *DarkNode) RunPacker() {
	go func() {
		running := int64(1)

		computationBlockChan := make(chan compute.ComputationBlock)
		go func() {
			for atomic.LoadInt64(&running) != 0 {
				computationBlockChan <- node.HiddenOrderBook.WaitForComputationBlock()
			}
		}()

		for atomic.LoadInt64(&running) != 0 {
			select {
			case computationBlock := <-computationBlockChan:
				node.BroadcastComputationBlock(computationBlock)
			case <-node.quitPacker:
				atomic.StoreInt64(&running, 0)
			}
		}
	}()
}

func (node *DarkNode) BroadcastComputationBlock(computationBlock ComputationBlock) {
	for _, multiAddress := range node.DarkPool {
		// TODO: send computation block to each node in our dark pool.
	}
}
