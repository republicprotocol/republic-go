package node

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
)

// Prime is the default prime number used to define the finite field.
var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

type DarkNode struct {
	Config

	TestDeltaNotifications chan *compute.Delta

	Logger     *logger.Logger
	ClientPool *rpc.ClientPool
	DHT        *dht.DHT

	DeltaBuilder                      *compute.DeltaBuilder
	DeltaFragmentMatrix               *compute.DeltaFragmentMatrix
	OrderFragmentWorkerQueue          chan *order.Fragment
	OrderFragmentWorker               *OrderFragmentWorker
	DeltaFragmentBroadcastWorkerQueue chan *compute.DeltaFragment
	DeltaFragmentBroadcastWorker      *DeltaFragmentBroadcastWorker
	DeltaFragmentWorkerQueue          chan *compute.DeltaFragment
	DeltaFragmentWorker               *DeltaFragmentWorker
	GossipWorkerQueue                 chan *compute.Delta
	GossipWorker                      *GossipWorker
	FinalizeWorkerQueue               chan *compute.Delta
	FinalizeWorker                    *FinalizeWorker
	ConsensusWorkerQueue              chan *compute.Delta
	ConsensusWorker                   *ConsensusWorker

	Server *grpc.Server
	Swarm  *network.SwarmService
	Dark   *network.DarkService
	Gossip *network.GossipService

	DarkOcean      *dark.Ocean
	DarkPool       *dark.Pool
	EpochBlockhash [32]byte
}

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewDarkNode(config Config, darkNodeRegistrar dnr.DarkNodeRegistrar) (*DarkNode, error) {
	if config.Prime == nil {
		config.Prime = Prime
	}

	// TODO: This should come from the DNR.
	k := int64(3) // 14)

	var err error
	node := &DarkNode{
		Config:                 config,
		TestDeltaNotifications: make(chan *compute.Delta, 100),
	}

	// Create the logger and start all plugins
	node.Logger, err = logger.NewLogger(config.LoggerOptions)
	if err != nil {
		return nil, err
	}
	node.Logger.Start()

	// Load the dark ocean and the dark pool for this node
	node.DarkPool = dark.NewPool()
	node.DarkOcean, err = dark.NewOcean(node.Logger, darkNodeRegistrar)
	if err != nil {
		return nil, err
	}

	// Create all networking components and services
	node.ClientPool = rpc.NewClientPool(node.NetworkOptions.MultiAddress).
		WithTimeout(node.NetworkOptions.Timeout).
		WithTimeoutBackoff(node.NetworkOptions.TimeoutBackoff).
		WithTimeoutRetries(node.NetworkOptions.TimeoutRetries).
		WithCacheLimit(node.NetworkOptions.ClientPoolCacheLimit)
	node.DHT = dht.NewDHT(node.NetworkOptions.MultiAddress.Address(), node.NetworkOptions.MaxBucketLength)
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	node.Swarm = network.NewSwarmService(node, node.NetworkOptions, node.Logger, node.ClientPool, node.DHT)
	node.Dark = network.NewDarkService(node, node.NetworkOptions, node.Logger)
	node.Gossip = network.NewGossipService(node)

	// Create all background workers that will do all of the actual work
	node.DeltaBuilder = compute.NewDeltaBuilder(k, node.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(node.Prime)
	node.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	node.OrderFragmentWorker = NewOrderFragmentWorker(node.Logger, node.DeltaFragmentMatrix, node.OrderFragmentWorkerQueue)
	node.DeltaFragmentBroadcastWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentBroadcastWorker = NewDeltaFragmentBroadcastWorker(node.Logger, node.ClientPool, node.DarkPool, node.DeltaFragmentBroadcastWorkerQueue)
	node.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentWorker = NewDeltaFragmentWorker(node.Logger, node.DeltaBuilder, node.DeltaFragmentWorkerQueue)
	node.GossipWorkerQueue = make(chan *compute.Delta, 100)
	node.GossipWorker = NewGossipWorker(node.Logger, node.ClientPool, node.NetworkOptions.BootstrapMultiAddresses, node.GossipWorkerQueue)
	node.FinalizeWorkerQueue = make(chan *compute.Delta, 100)
	node.FinalizeWorker = NewFinalizeWorker(node.Logger, k, node.FinalizeWorkerQueue)
	node.ConsensusWorkerQueue = make(chan *compute.Delta, 100)
	node.ConsensusWorker = NewConsensusWorker(node.Logger, node.DeltaFragmentMatrix, node.ConsensusWorkerQueue)

	return node, nil
}

// Start the DarkNode.
func (node *DarkNode) Start() {
	// FIXME: This causes an index out of bounds panic.
	// // Broadcast CPU/Memory/Network usage
	// go func() {
	// 	for {
	// 		node.Usage()
	// 		time.Sleep(20 * time.Second)
	// 	}
	// }()

	// Start background workers
	go node.OrderFragmentWorker.Run(node.DeltaFragmentBroadcastWorkerQueue, node.DeltaFragmentWorkerQueue)
	go node.DeltaFragmentBroadcastWorker.Run()
	go node.DeltaFragmentWorker.Run(node.GossipWorkerQueue, node.TestDeltaNotifications)
	go node.GossipWorker.Run(node.FinalizeWorkerQueue)
	go node.FinalizeWorker.Run(node.ConsensusWorkerQueue)
	go node.ConsensusWorker.Run()

	// Start gRPC services and UI
	go node.StartServices()
	time.Sleep(5 * time.Second)

	//go node.StartUI()
	time.Sleep(time.Second)

	// Bootstrap into the swarm network
	node.Swarm.Bootstrap()
}

func (node *DarkNode) StartServices() {
	node.Logger.Info(logger.TagNetwork, fmt.Sprintf("Listening on %s:%s", node.Host, node.Port))
	node.Swarm.Register(node.Server)
	node.Dark.Register(node.Server)
	node.Gossip.Register(node.Server)
	listener, err := net.Listen("tcp", node.Host+":"+node.Port)
	if err != nil {
		node.Logger.Error(logger.TagNetwork, err.Error())
	}
	if err := node.Server.Serve(listener); err != nil {
		node.Logger.Error(logger.TagNetwork, err.Error())
	}
}

func (node *DarkNode) StartUI() {
	fs := http.FileServer(http.Dir("dark-node-ui"))
	http.Handle("/", fs)
	node.Logger.Info(logger.TagNetwork, "Serving the Dark Node UI")
	err := http.ListenAndServe("0.0.0.0:3000", nil)
	if err != nil {
		node.Logger.Error(logger.TagNetwork, err.Error())
	}
}

// Stop the DarkNode.
func (node *DarkNode) Stop() {
	// Stop serving gRPC services
	node.Server.Stop()
	time.Sleep(time.Second)

	// Stop background workers by closing their job queues
	close(node.OrderFragmentWorkerQueue)
	close(node.DeltaFragmentBroadcastWorkerQueue)
	close(node.DeltaFragmentWorkerQueue)
	close(node.GossipWorkerQueue)
	close(node.FinalizeWorkerQueue)
	close(node.ConsensusWorkerQueue)

	// Stop the logger
	node.Logger.Stop()

	// Force the GC to run
	runtime.GC()
}

// WatchDarkOcean for changes. When a change happens, find the dark pool for
// this DarkNode and reconnect to all of the nodes in the pool.
func (node *DarkNode) WatchDarkOcean() {
	if err := node.DarkOcean.Update(); err != nil {
		node.Logger.Error(logger.TagEthereum, fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
	}
	changes := make(chan struct{})
	go func() {
		defer close(changes)
		for {
			// Wait for a change to the ocean
			<-changes
			node.Logger.Info(logger.TagEthereum, "dark ocean change detected")

			// Find the dark pool for this node and connect to all of the dark
			// nodes in the pool
			node.DarkPool.RemoveAll()
			darkPool := node.DarkOcean.FindPool(node.RepublicKeyPair.ID())
			node.ConnectToDarkPool(darkPool)

		}
	}()
	node.DarkOcean.Watch(5*time.Minute, changes)
}

// ConnectToDarkPool and return the connected nodes and disconnected nodes
// separately.
func (node *DarkNode) ConnectToDarkPool(darkPool *dark.Pool) {
	// Terminate if the dark pool is no longer relevant
	if darkPool == nil {
		return
	}

	darkPool.ForAll(func(n *dark.Node) {
		if bytes.Equal(node.Config.RepublicKeyPair.ID(), n.ID) {
			return
		}
		multiAddress := n.MultiAddress()
		if multiAddress != nil {
			return
		}
		// Find the dark node
		multiAddress, err := node.Swarm.FindNode(n.ID)
		if err != nil {
			node.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot find dark node %v: %s", n.ID.Address(), err.Error()))
			return
		} else if multiAddress == nil {
			// node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("cannot find dark node: %v", n.ID.Address()))
			return
		}

		// Ping the dark node to test the connection
		node.ClientPool.Ping(*multiAddress)
		if err != nil {
			node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("cannot ping to dark node %v: %s", n.ID.Address(), err.Error()))
			return
		}

		// Update the DHT
		err = node.Swarm.DHT.UpdateMultiAddress(*multiAddress)
		if err != nil {
			node.Logger.Warn(logger.TagNetwork, fmt.Sprintf("cannot update DHT with dark node %v: %s", n.ID.Address(), err.Error()))
			return
		}

		node.Logger.Info(logger.TagNetwork, fmt.Sprintf("found dark node: %v", n.ID.Address()))

		// Update the MultiAddress in the node
		n.SetMultiAddress(*multiAddress)
		node.DarkPool.Append(*n)
	})

	// In the background, continue to attempt connections to the disconnected
	// dark nodes in the pool
	go func() {
		time.Sleep(time.Minute)
		hasDisconnectedNodes := int64(0)
		darkPool.For(func(n *dark.Node) {
			if n.MultiAddress() == nil {
				atomic.StoreInt64(&hasDisconnectedNodes, 1)
			}
		})
		if hasDisconnectedNodes != 0 {
			node.ConnectToDarkPool(darkPool)
		}
	}()
}

// OnOpenOrder writes an order fragment that has been received to the
// OrderFragmentWorkerQueue. This is a potentially blocking operation, however
// this delegate method is called on a dedicated goroutine.
func (node *DarkNode) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.OrderFragmentWorkerQueue <- orderFragment
	}()
}

// OnBroadcastDeltaFragment writes a delta fragment that has been received to
// the DeltaFragmentWorkerQueue. This is a potentially blocking operation,
// however this delegate method is called on a dedicated goroutine.
func (node *DarkNode) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.DeltaFragmentWorkerQueue <- deltaFragment
	}()
}

func (node *DarkNode) OnGossip(buyOrderID order.ID, sellOrderID order.ID) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.GossipWorkerQueue <- &compute.Delta{
			ID:          compute.DeltaID(crypto.Keccak256([]byte(buyOrderID), []byte(sellOrderID))),
			BuyOrderID:  buyOrderID,
			SellOrderID: sellOrderID,
		}
	}()
}

func (node *DarkNode) OnFinalize(buyOrderID order.ID, sellOrderID order.ID) {
	// Write to a channel that might be closed
	func() {
		defer func() { recover() }()
		node.FinalizeWorkerQueue <- &compute.Delta{
			ID:          compute.DeltaID(crypto.Keccak256([]byte(buyOrderID), []byte(sellOrderID))),
			BuyOrderID:  buyOrderID,
			SellOrderID: sellOrderID,
		}
	}()
}

// Usage logs memory and cpu usage
func (node *DarkNode) Usage() {
	// memory
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}
	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}
	percentage, err := cpu.Percent(0, false)
	if err != nil {
		node.Logger.Error(logger.TagUsage, err.Error())
	}

	node.Logger.Usage(float32(cpuStat[0].Mhz*percentage[0]/100), int32(vmStat.Used/1024/1024), 0)
}
