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

	"encoding/hex"
	"encoding/json"

	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/rs/cors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	ionet "github.com/shirou/gopsutil/net"
	"google.golang.org/grpc"
)

// Prime is the default prime number used to define the finite field.
var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

type DarkNode struct {
	Config
	identity.KeyPair
	identity.ID
	identity.Address

	DeltaNotifications chan *compute.Delta

	LoggerLastNetworkUsage uint64
	Logger                 *logger.Logger
	ClientPool             *rpc.ClientPool
	DHT                    *dht.DHT

	DeltaBuilder                      *compute.DeltaBuilder
	DeltaFragmentMatrix               *compute.DeltaFragmentMatrix
	OrderFragmentWorkerQueue          chan *order.Fragment
	OrderFragmentWorker               *OrderFragmentWorker
	DeltaFragmentBroadcastWorkerQueue chan *compute.DeltaFragment
	DeltaFragmentBroadcastWorker      *DeltaFragmentBroadcastWorker
	DeltaFragmentWorkerQueue          chan *compute.DeltaFragment
	DeltaFragmentWorker               *DeltaFragmentWorker
	DeltaQueue                        chan *compute.Delta
	DeltaMatchWorker                  *DeltaMatchWorker

	Server *grpc.Server
	Swarm  *network.SwarmService
	Dark   *network.DarkService

	DarkNodeRegistrar dnr.DarkNodeRegistrar
	DarkOcean         *dark.Ocean
	DarkPool          *dark.Pool
	EpochBlockhash    [32]byte
}

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewDarkNode(config Config, darkNodeRegistrar dnr.DarkNodeRegistrar) (*DarkNode, error) {
	var err error
	node := &DarkNode{
		Config:             config,
		KeyPair:            config.KeyPair,
		ID:                 config.KeyPair.ID(),
		Address:            config.KeyPair.Address(),
		DeltaNotifications: make(chan *compute.Delta, 100),
	}

	// Create the logger and start all plugins
	node.Logger, err = logger.NewLogger(config.LoggerOptions)
	if err != nil {
		return nil, err
	}
	node.Logger.Start()

	// Load the dark ocean and the dark pool for this node
	node.DarkNodeRegistrar = darkNodeRegistrar
	node.DarkOcean, err = dark.NewOcean(node.Logger, darkNodeRegistrar)
	if err != nil {
		return nil, err
	}
	node.DarkPool = dark.NewPool()
	if darkPool := node.DarkOcean.FindPool(node.ID); darkPool != nil {
		node.DarkPool = darkPool
	}
	k := int64(node.DarkPool.Size()*2/3 + 1)

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

	// Create all background workers that will do all of the actual work
	node.DeltaBuilder = compute.NewDeltaBuilder(k, Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(Prime)
	node.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	node.OrderFragmentWorker = NewOrderFragmentWorker(node.Logger, node.DeltaFragmentMatrix, node.OrderFragmentWorkerQueue)
	node.DeltaFragmentBroadcastWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentBroadcastWorker = NewDeltaFragmentBroadcastWorker(node.Logger, node.ClientPool, node.DarkPool, node.DeltaFragmentBroadcastWorkerQueue)
	node.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentWorker = NewDeltaFragmentWorker(node.Logger, node.DeltaBuilder, node.DeltaFragmentWorkerQueue)
	node.DeltaQueue = make(chan *compute.Delta, 100)
	node.DeltaMatchWorker = NewDeltaMatchWorker(node.Logger, node.DeltaFragmentMatrix, node.DeltaQueue)

	return node, nil
}

func (node *DarkNode) StartBackgroundWorkers() {
	// Usage logger
	go func() {
		for {
			time.Sleep(10 * time.Second)
			node.Usage(10)
		}
	}()

	// Start background workers
	go node.OrderFragmentWorker.Run(node.DeltaFragmentBroadcastWorkerQueue, node.DeltaFragmentWorkerQueue)
	go node.DeltaFragmentBroadcastWorker.Run()
	go node.DeltaFragmentWorker.Run(node.DeltaQueue)
	go node.DeltaMatchWorker.Run(node.DeltaNotifications)
}

func (node *DarkNode) StartServices() {
	node.Logger.Network(logger.Info, fmt.Sprintf("gRPC services listening on %s:%s", node.Host, node.Port))

	node.Swarm.Register(node.Server)
	node.Dark.Register(node.Server)
	listener, err := net.Listen("tcp", node.Host+":"+node.Port)
	if err != nil {
		node.Logger.Error(err.Error())
	}
	if err := node.Server.Serve(listener); err != nil {
		node.Logger.Error(err.Error())
	}
}

func (node *DarkNode) StartUI() {

	host, err := node.NetworkOptions.MultiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		node.Logger.Network(logger.Error, "UI host unknown: "+err.Error())
		return
	}
	port, err := node.NetworkOptions.MultiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		node.Logger.Network(logger.Error, "UI port unknown: "+err.Error())
		return
	}

	node.Logger.Network(logger.Info, fmt.Sprintf("UI listening on %s:%s", node.Host, "3000"))
	http.Handle("/config", cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"nodeID":     "0x" + hex.EncodeToString(node.ID),
			"publicKey":  "0x" + hex.EncodeToString(append(node.KeyPair.PublicKey.X.Bytes(), node.KeyPair.PublicKey.Y.Bytes()...)),
			"address":    node.EthereumKey.Address.String(),
			"republicID": node.ID.String(),
			"ui": map[string]interface{}{
				"host": host,
				"port": port,
			},
			"websocket": map[string]interface{}{
				"host": host,
				"port": "18515",
			},
			"contracts": map[string]interface{}{
				"darkNodeRegistrar": "0xf178237e7d1131b7924435aa8d02B8Ab4d308AFf",
			},
		})
	})))
	http.Handle("/", http.FileServer(http.Dir("/home/.darknode/ui")))
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		node.Logger.Error(err.Error())
	}
}

func (node *DarkNode) Bootstrap() {
	node.Swarm.Bootstrap()
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
	close(node.DeltaQueue)
	close(node.DeltaNotifications)

	// Stop the logger
	node.Logger.Stop()

	// Force the GC to run
	runtime.GC()
}

// WatchDarkOcean for changes. When a change happens, find the dark pool for
// this DarkNode and reconnect to all of the nodes in the pool.
func (node *DarkNode) WatchDarkOcean() {
	// Block until the node is registered
	err := node.DarkNodeRegistrar.WaitUntilRegistration(node.ID)
	for err != nil {
		node.Logger.Error(fmt.Sprintf("cannot determine registration status: %s", err.Error()))

		// Wait for 5 seconds and try again
		time.Sleep(5 * time.Second)
		err = node.DarkNodeRegistrar.WaitUntilRegistration(node.ID)
	}

	changes := make(chan struct{})
	go func() {
		defer close(changes)
		for {
			// Wait for a change to the ocean
			<-changes
			node.Logger.Info("dark ocean change detected")

			// Find the dark pool for this node and connect to all of the dark
			// nodes in the pool
			node.DarkPool.RemoveAll()
			darkPool := node.DarkOcean.FindPool(node.ID)
			if darkPool != nil {
				k := int64((darkPool.Size() * 2 / 3) + 1)
				node.DeltaBuilder.SetK(k)
			}
			node.ConnectToDarkPool(darkPool)
		}
	}()

	// Check for changes every minute
	node.DarkOcean.Watch(time.Minute, changes)
}

// ConnectToDarkPool and return the connected nodes and disconnected nodes
// separately.
func (node *DarkNode) ConnectToDarkPool(darkPool *dark.Pool) {
	// Terminate if the dark pool is no longer relevant
	if darkPool == nil {
		return
	}

	darkPool.ForAll(func(n *dark.Node) {
		if bytes.Equal(node.ID, n.ID) {
			return
		}
		multiAddress := n.MultiAddress()
		if multiAddress != nil {
			return
		}
		// Find the dark node
		multiAddress, err := node.Swarm.FindNode(n.ID)
		if err != nil {
			node.Logger.Error(fmt.Sprintf("cannot find dark node %v: %s", n.ID.Address(), err.Error()))
			return
		} else if multiAddress == nil {
			node.Logger.Warn(fmt.Sprintf("cannot find dark node %v: %s", n.ID.Address(), err.Error()))
			return
		}

		// Ping the dark node to test the connection
		node.ClientPool.Ping(*multiAddress)
		if err != nil {
			node.Logger.Warn(fmt.Sprintf("cannot ping to dark node %v: %s", n.ID.Address(), err.Error()))
			return
		}

		// Update the DHT
		err = node.Swarm.DHT.UpdateMultiAddress(*multiAddress)
		if err != nil {
			node.Logger.Warn(fmt.Sprintf("cannot update DHT with dark node %v: %s", n.ID.Address(), err.Error()))
			return
		}

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
		node.Logger.OrderReceived(logger.Info, orderFragment.OrderID.String(), orderFragment.ID.String())
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

// Usage logs memory and cpu usage
func (node *DarkNode) Usage(seconds uint64) {

	// Get CPU usage
	_, err := cpu.Info()
	if err != nil {
		node.Logger.Error(err.Error())
	}
	var cpuPercentage float64
	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		node.Logger.Error(err.Error())
	}
	if len(cpuPercentages) > 0 {
		cpuPercentage = cpuPercentages[0]
	}

	// Get memory usage
	stat, err := mem.VirtualMemory()
	if err != nil {
		node.Logger.Error(err.Error())
	}
	memoryPercentage := stat.UsedPercent

	ios, err := ionet.IOCounters(false)
	if err != nil {
		node.Logger.Error(err.Error())
	}

	// Get network usage
	var networkUsage uint64
	if len(ios) > 0 {
		networkUsage += ios[0].BytesSent
		networkUsage += ios[0].BytesRecv
	}
	if node.LoggerLastNetworkUsage == 0 {
		node.LoggerLastNetworkUsage = networkUsage
	}

	node.Logger.Usage(cpuPercentage, memoryPercentage, (networkUsage-node.LoggerLastNetworkUsage)/seconds)
	node.LoggerLastNetworkUsage = networkUsage
}
