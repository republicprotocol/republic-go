package node

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-do"
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
	DarkPool        identity.MultiAddresses
	DarkPoolLimit   int

	quitServer chan struct{}
	quitPacker chan struct{}
}

// NewDarkNode creates a new DarkNode, a new swarm.Node and dark.Node and assigns the
// new DarkNode as the delegate for both. Returns the new DarkNode, or an error.
func NewDarkNode(config *Config) (*DarkNode, error) {
	node := &DarkNode{
		HiddenOrderBook: compute.NewHiddenOrderBook(config.ComputationChunkSize),
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
		log.Printf("Listening on %s:%s\n", node.Configuration.Host, node.Configuration.Port)
		node.Swarm.Register()
		node.Dark.Register()
		listener, err := net.Listen("tcp", node.Configuration.Host+":"+node.Configuration.Port)
		if err != nil {
			log.Fatal(err)
		}
		if err := node.Server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the server to start and bootstrap the connections in the swarm.
	time.Sleep(time.Second)
	node.Swarm.Bootstrap()

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
func (node *DarkNode) OnPingReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersReceived implements the swarm.Delegate interface. It is
// used by the underlying swarm.Node whenever the Miner has handled a
// QueryCloserPeers RPC.
func (node *DarkNode) OnQueryCloserPeersReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersOnFrontierReceived implements the swarm.Delegate
// interface. It is called by the underlying swarm.Node whenever the Miner
// has handled a QueryCloserPeersOnFrontier RPC.
func (node *DarkNode) OnQueryCloserPeersOnFrontierReceived(peer identity.MultiAddress) {
}

// OnOrderFragmentReceived implements the xing.Delegate interface. It is called
// by the underlying xing.Node whenever the Miner receives a
// compute.OrderFragment that it must process.
func (node *DarkNode) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	node.HiddenOrderBook.AddOrderFragment(orderFragment)
}

// OnResultFragmentReceived implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.ResultFragment that it must process.
func (node *DarkNode) OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment) {
}

// OnOrderFragmentForwarding implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.OrderFragment that it must forward to the correct xing.Node.
func (node *DarkNode) OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment) {
}

// OnResultFragmentForwarding implements the xing.Delegate interface. It is
// called by the underlying xing.Node whenever the Miner receives a
// compute.ResultFragment that it must forward to the correct xing.Node.
func (node *DarkNode) OnResultFragmentForwarding(to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment) {
}

func (node *DarkNode) RunComputationChunkPacker() {
	go func() {
		running := int64(1)

		computationChunkChan := make(chan compute.ComputationChunk)
		go func() {
			for atomic.LoadInt64(&running) != 0 {
				computationChunkChan <- node.HiddenOrderBook.WaitForComputationChunk()
			}
		}()

		for atomic.LoadInt64(&running) != 0 {
			select {
			case computationChunk := <-computationChunkChan:
				go func() {
					computationChunkConsensus := node.BroadcastComputationChunk(computationChunk)
					node.BroadcastComputationChunkConsensus(computationChunkConsensus)
				}()
			case <-time.Tick(time.Second * node.Configuration.ComputationChunkInterval):
				preemptedComputationChunk := node.HiddenOrderBook.PreemptComputationChunk()
				if len(preemptedComputationChunk) > 0 {
					go func() {
						computationChunkConsensus := node.BroadcastComputationChunk(preemptedComputationChunk)
						node.BroadcastComputationChunkConsensus(computationChunkConsensus)
					}()
				}
			case <-node.quitPacker:
				atomic.StoreInt64(&running, 0)
			}
		}
	}()
}

func (node *DarkNode) BroadcastComputationChunk(computationChunk compute.ComputationChunk) {
	// Track all of the no bids on computations.
	nosMu := new(sync.Mutex)
	nos := map[string]int{}

	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		// TODO: send computation chunk and get response
		// computationChunkBids := peer.SendComputationChunkProposal(computationChunk)
		computationChunkBids := struct{}{}

		func() {
			nosMu.Lock()
			defer nosMu.Unlock()
			for computationID, bid := range computationChunkBids {
				if bid == compute.ComputationBidYes {
					continue
				}
				nos[computationID]++
			}
		}()

		// FIXME: if a node timeouts or doesn't respond then classify that as a
		// no for all computations.
	})

	// Create a new slice of compute.Computations for which more than
	// 2/3rds of the dark pool can participate.
	computations := make([]*compute.Computation, 0, len(computationChunk.Computations))
	for _, computation := range computationChunk.Computations {
		if noBids[string(computation.ID)] >= node.DarkPoolLimit {
			continue
		}
		computations = append(computations, computation)
	}

	computationChunkConsensus := compute.NewComputationChunk(computations)
	return computationChunkConsensus
}

func (node *DarkNode) BroadcastComputationChunkConsensus(computationChunk compute.ComputationChunk) {
	go func() {
		node.Compute(computationChunk)
	}()
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		// TODO: send computation chunk and get response
		// peer.SendComputationChunkConsensus(computationChunk)
	})
}

func (node *DarkNode) Compute(computationChunk compute.ComputationChunk) {
	resultFragments := computationChunk.Compute(node.Configuration.Prime)
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		rpc.SendResultFragmentToTarget() // FIXME: Finish calling this RPC
	})
}
