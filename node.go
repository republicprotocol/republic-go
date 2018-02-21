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
		HiddenOrderBook: compute.NewHiddenOrderBook(config.ComputationShardSize),
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
			node.DarkPool = append(node.DarkPool, target)
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
	// TODO: Log metrics for the ping.
}

// OnQueryCloserPeersReceived implements the swarm.Delegate interface. It is
// used by the underlying swarm.Node whenever the Miner has handled a
// QueryCloserPeers RPC.
func (node *DarkNode) OnQueryCloserPeersReceived(peer identity.MultiAddress) {
	// TODO: Log metrics for the query.
}

// OnQueryCloserPeersOnFrontierReceived implements the swarm.Delegate
// interface. It is called by the underlying swarm.Node whenever the Miner
// has handled a QueryCloserPeersOnFrontier RPC.
func (node *DarkNode) OnQueryCloserPeersOnFrontierReceived(peer identity.MultiAddress) {
	// TODO: Log metrics for the deep query.
}

// OnOrderFragmentReceived implements the xing.Delegate interface. It is called
// by the underlying xing.Node whenever the Miner receives a
// compute.OrderFragment that it must process.
func (node *DarkNode) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	node.HiddenOrderBook.AddOrderFragment(orderFragment)
}

// OnElectShard is a delegate method that is called when the DarkNode has
// received an RPC for electing a Shard. To finish the election, the DarkNode
// should filter the Shard is received and return it. The DarkNode should
// filter out the deltas and residues that it does not have access to.
func (node *DarkNode) OnElectShard(from identity.MultiAddress, shard compute.ComputationShard) compute.ComputationShard {
	// TODO: Elect the shard. Check which deltas and residues the DarkNode has
	// and remove all others from the shard before returning it.
	panic("unimplemented")
}

// OnComputeShard is a delegate method that is called when the DarkNode has
// received an RPC for computing a Shard. The process for computing a Shard is
// to compute ([[a]] - [[b]]) * [[r^2]] using the ([[a]] - [[b]]) deltas and
// the [[r^2]] residues in the Shard.
func (node *DarkNode) OnComputeShard(from identity.MultiAddress, shard compute.ComputationShard) {
	// TODO: Compute the shard. At this stage of the implementation, no
	// copmutation is actually needed and the Shard can be finalized
	// immediately.
	panic("unimplemented")
}

// OnFinalizeShard is a delegate method that is called when the DarkNode has
// received an RPC for finalizing a Shard. The process for finalizing a Shard
// is to store all computation shares from within the different Shards that are
// finalized. Eventually, enough computation shares will be acquired and the
// computation proper can be reconstructed.
func (node *DarkNode) OnFinalizeShard(from identity.MultiAddress, shard compute.ComputationShard) {
	// TODO: Store the shares in a map until we have enough to reconstruct the
	// computation proper. After reconstruction, finalize the computation and
	// stop all processing for it (elections, computations, finalizations).
	panic("unimplemented")

	// FIXME: After reconstruction there should be some interaction with the
	// traders.
}

func (node *DarkNode) RunComputationShardPacker() {
	go func() {
		running := int64(1)

		computationShardChan := make(chan compute.ComputationShard)
		go func() {
			for atomic.LoadInt64(&running) != 0 {
				computationShardChan <- node.HiddenOrderBook.WaitForComputationShard()
			}
		}()

		for atomic.LoadInt64(&running) != 0 {
			select {
			case computationShard := <-computationShardChan:
				go func() {
					computationShardConsensus := node.BroadcastComputationShard(computationShard)
					node.BroadcastComputationShardConsensus(computationShardConsensus)
				}()
			case <-time.Tick(time.Duration(node.Configuration.ComputationShardInterval) * time.Second):
				preemptedComputationShard := node.HiddenOrderBook.PreemptComputationShard()
				if len(preemptedComputationShard.Computations) > 0 {
					go func() {
						computationShardConsensus := node.BroadcastComputationShard(preemptedComputationShard)
						node.BroadcastComputationShardConsensus(computationShardConsensus)
					}()
				}
			case <-node.quitPacker:
				atomic.StoreInt64(&running, 0)
			}
		}
	}()
}

func (node *DarkNode) BroadcastComputationShard(computationShard compute.ComputationShard) {
	// Track all of the no bids on computations.
	nosMu := new(sync.Mutex)
	nos := map[string]int{}

	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		// TODO: send computation shard and get response
		// computationShardBids := peer.SendComputationShardProposal(computationShard)
		computationShardBids := struct{}{}

		func() {
			nosMu.Lock()
			defer nosMu.Unlock()
			for computationID, bid := range computationShardBids {
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
	computations := make([]*compute.Computation, 0, len(computationShard.Computations))
	for _, computation := range computationShard.Computations {
		if noBids[string(computation.ID)] >= node.DarkPoolLimit {
			continue
		}
		computations = append(computations, computation)
	}

	computationShardConsensus := compute.NewComputationShard(computations)
	return computationShardConsensus
}

func (node *DarkNode) BroadcastComputationShardConsensus(computationShard compute.ComputationShard) {
	go func() {
		node.Compute(computationShard)
	}()
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		// TODO: send computation chunk and get response
		// peer.SendComputationShardConsensus(computationShard)
	})
}

func (node *DarkNode) Compute(computationShard compute.ComputationShard) {
	resultFragments := computationShard.Compute(node.Configuration.Prime)
	do.ForAll(node.DarkPool, func(i int) {
		peer := node.DarkPool[i]
		rpc.SendResultFragmentToTarget() // FIXME: Finish calling this RPC
	})
}

func (node *DarkNode) BidOnComputationShard(computationShard compute.ComputationShard) compute.ComputationShardBid {
	computationShardBid := compute.ComputationShardBid{
		ID:   computationShard.ID,
		Bids: map[string]compute.ComputationBid{},
	}
	for _, computation := range computationShard.Computations {
		computationShardBid.Bids[string(computation.ID)] = compute.ComputationBidNo
	}
	// FIXME:
	return computationShardBid
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
	return userConnection.IsDarkNodeRegistered(idInBytes)
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
