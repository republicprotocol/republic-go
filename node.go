package node

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-atom/ethereum"
	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-dark-node-registrar"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
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

	registered, err := isRegistered(node.Configuration.MultiAddress.ID())
	if err != nil {
		log.Fatal(err)
	}
	if !registered {
		panic("dark node hasn't been registered")
	}

	darkPool, err := getDarkPoolConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the server to start and bootstrap the connections in the swarm.
	time.Sleep(time.Second)
	node.Swarm.Bootstrap()

	//  Ping all nodes in the dark pool
	for _, id := range darkPool {
		target, err := node.Swarm.FindNode(id)
		if err != nil {
			log.Fatal(err)
		}
		// Ignore the node if we can't find it
		if target == nil {
			continue
		}
		err = rpc.PingTarget(target, node.Swarm.MultiAddress(), 5*time.Second)
		// Update the nodes in our DHT if they respond
		if err != nil {
			node.Swarm.DHT.UpdateMultiAddress(target)
		}
	}

	// todo : sync hidden order book with them

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

		computationChunkChan := make(chan compute.ComputationBlock)
		go func() {
			for atomic.LoadInt64(&running) != 0 {
				computationChunkChan <- node.HiddenOrderBook.WaitForComputationBlock()
			}
		}()

		for atomic.LoadInt64(&running) != 0 {
			select {
			case computationChunk := <-computationChunkChan:
				go func() {
					computationChunkConsensus := node.BroadcastComputationChunk(computationChunk)
					node.BroadcastComputationChunkConsensus(computationChunkConsensus)
				}()
			case <-time.Tick(time.Duration(node.Configuration.ComputationChunkInterval) * time.Second):
				preemptedComputationChunk := node.HiddenOrderBook.PreemptComputationBlock()
				if len(preemptedComputationChunk.Computations) > 0 {
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

func (node *DarkNode) BroadcastComputationChunk(computationChunk compute.ComputationBlock) {
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

	computationChunkConsensus := compute.NewComputationBlock(computations)
	return computationChunkConsensus
}

func (node *DarkNode) BroadcastComputationChunkConsensus(computationChunk compute.ComputationBlock) {
	go func() {
		node.Compute(computationChunk)
	}()
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		// TODO: send computation chunk and get response
		// peer.SendComputationChunkConsensus(computationChunk)
	})
}

func (node *DarkNode) Compute(computationChunk compute.ComputationBlock) {
	resultFragments := computationChunk.Compute(node.Configuration.Prime)
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		rpc.SendResultFragmentToTarget() // FIXME: Finish calling this RPC
	})
}

func (node *DarkNode) BidOnComputationChunk(computationChunk compute.ComputationBlock) compute.ComputationBlockBid {
	computationChunkBid := compute.ComputationBlockBid{
		ID:   computationChunk.ID,
		Bids: map[string]compute.ComputationBid{},
	}
	for _, computation := range computationChunk.Computations {
		computationChunkBid.Bids[string(computation.ID)] = compute.ComputationBidNo
	}
	// FIXME:
	return computationChunkBid
}

func isRegistered(id identity.ID) (bool, error) {
	// todo: need to get key from the ethereum private key
	key := `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		return false, err
	}
	client := ethereum.Ropsten("https://ropsten.infura.io/")
	contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")
	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), client, auth, &bind.CallOpts{}, contractAddress, nil)
	idInBytes := [20]byte{}
	copy(idInBytes[:], id)
	return userConnection.IsDarkNodeRegistered(id)
}

func getDarkPoolConfig() ([]identity.ID, error) {
	// todo: need to get key from the ethereum private key
	key := `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		return []identity.ID{}, err
	}
	client := ethereum.Ropsten("https://ropsten.infura.io/")
	contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")
	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), client, auth, &bind.CallOpts{}, contractAddress, nil)
	ids, err := userConnection.GetXingOverlay()
	if err != nil {
		return []identity.ID{}, err
	}
	nodes := make([]identity.ID, len(ids))
	for i := range ids {
		nodes[i] = identity.ID(ids[i][:])
	}
	return nodes, nil
}
