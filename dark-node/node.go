package node

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
)

const defaultTimeout = 5 * time.Second

// To be retrieved from the Registrar contract
var (
	// Prime ...
	Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
)

// LogQueue allows multiple clients to receive logs from a node
type LogQueue struct {
	do.GuardedObject

	channels []chan do.Option
}

// NewLogQueue returns a new LogQueue
func NewLogQueue() *LogQueue {
	logQueue := new(LogQueue)
	logQueue.GuardedObject = do.NewGuardedObject()
	logQueue.channels = nil
	return logQueue
}

// Publish allows a node to push a log to each client
func (logQueue *LogQueue) Publish(val do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()

	var logQueueLength = len(logQueue.channels)
	for i := 0; i < logQueueLength; i++ {
		timer := time.NewTicker(10 * time.Second)
		defer timer.Stop()
		select {
		case logQueue.channels[i] <- val:
		case <-timer.C:
			// Deregister the channel
			logQueue.channels[i] = logQueue.channels[logQueueLength-1]
			logQueue.channels = logQueue.channels[:logQueueLength-1]
			logQueueLength--
			i--
		}
	}
}

// Subscribe allows a new client to listen to events from a node
func (logQueue *LogQueue) Subscribe(channel chan do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()

	logQueue.channels = append(logQueue.channels, channel)
}

// Unsubscribe ...
func (logQueue *LogQueue) Unsubscribe(channel chan do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()
	length := len(logQueue.channels)
	for i := 0; i < length; i++ {
		// https://golang.org/ref/spec#Comparison_operators
		// Two channel values are equal if they were created by the same call to make
		// or if both have value nil.
		if logQueue.channels[i] == channel {
			logQueue.channels[i] = logQueue.channels[length-1]
			logQueue.channels = logQueue.channels[:length-1]
			break
		}
	}
}

// DarkNode ...
type DarkNode struct {
	Server        *grpc.Server
	Swarm         *network.SwarmService
	Dark          *network.DarkOceanService
	Configuration *Config
	Registrar     *dnr.DarkNodeRegistrar

	DeltaBuilder        *compute.DeltaBuilder
	DeltaFragmentMatrix *compute.DeltaFragmentMatrix
	DarkPoolLimit       int64
	DarkPool            dnr.DarkPool
	DarkOceanOverlay    *dnr.DarkOceanOverlay

	EpochBlockhash [32]byte

	logQueue *LogQueue
}

// NewDarkNode creates a new DarkNode, a new network.Node and network.Node and assigns the
// new DarkNode as the delegate for both. Returns the new DarkNode, or an error.
func NewDarkNode(config *Config) (*DarkNode, error) {
	if config.Prime == nil {
		config.Prime = Prime
	}
	if config.ComputationShardInterval == 0 {
		config.ComputationShardInterval = 5
	}
	if config.ComputationShardSize == 0 {
		config.ComputationShardSize = 10
	}

	ethereumKeyPair, err := config.EthereumKeyPair()
	if err != nil {
		return nil, err
	}
	registrar, err := ConnectToRegistrar(ethereumKeyPair)
	if err != nil {
		return nil, err
	}

	node := &DarkNode{
		Server:        grpc.NewServer(grpc.ConnectionTimeout(time.Minute)),
		Configuration: config,
		Registrar:     registrar,
		logQueue:      NewLogQueue(),
	}
	// node.DarkPool = config.BootstrapMultiAddresses
	node.DarkPoolLimit = 5
	node.DeltaBuilder = compute.NewDeltaBuilder(node.DarkPoolLimit, config.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(config.Prime)

	swarmOptions := network.Options{
		MultiAddress:            config.MultiAddress,
		BootstrapMultiAddresses: config.BootstrapMultiAddresses,
		Debug:           network.DebugMedium,
		Alpha:           3,
		MaxBucketLength: 20,
		Timeout:         30 * time.Second,
		TimeoutStep:     30 * time.Second,
		TimeoutRetries:  3,
		Concurrent:      true,
	}
	swarmNode := network.NewNode(node.Server, node, swarmOptions)
	node.Swarm = swarmNode

	darkOptions := network.Options{
		Address:        config.MultiAddress.Address(),
		Debug:          network.DebugMedium,
		Timeout:        30 * time.Second,
		TimeoutStep:    30 * time.Second,
		TimeoutRetries: 3,
		Concurrent:     true,
	}
	darkNode := network.NewNode(node.Server, node, darkOptions)
	node.Dark = darkNode

	err = node.Configuration.Logger.Start()
	if err != nil {
		return node, err
	}

	return node, nil
}

// Start mining for compute.Orders that are matched. It establishes connections
// to other peers in the swarm network by bootstrapping against a set of
// bootstrap network.Nodes.
func (node *DarkNode) Start() error {
	go func() {
		for {
			time.Sleep(20 * time.Second)
			node.Usage()
		}
	}()

	isRegistered := node.IsRegistered()
	for !isRegistered {
		timeout := 60 * time.Second
		log.Printf("%v not registered. Sleeping for %v seconds.", node.Configuration.MultiAddress.Address(), timeout.Seconds())
		time.Sleep(timeout)
		isRegistered = node.IsRegistered()
	}

	// Bootstrap the connections in the network.
	node.Swarm.Bootstrap()

	go node.WatchForEpoch()

	return nil
}

func (node *DarkNode) Error(err error) {
	node.Configuration.Logger.Error(err)
}

func (node *DarkNode) Info(info string) {
	node.Configuration.Logger.Info(info)
}

func (node *DarkNode) Warning(warning string) {
	node.Configuration.Logger.Warning(warning)
}

// StartListening starts listening for rpc calls
func (node *DarkNode) StartListening() error {
	node.Info(fmt.Sprintf("Listening on %s:%s\n", node.Configuration.Host, node.Configuration.Port))
	node.Swarm.Register()
	node.Dark.Register()
	listener, err := net.Listen("tcp", node.Configuration.Host+":"+node.Configuration.Port)
	if err != nil {
		return err
	}
	return node.Server.Serve(listener)
}

// StopListening stops listening for rpc calls
func (node *DarkNode) StopListening() {
	node.Server.Stop()
}

func (node *DarkNode) log(kind, message string) {
	node.logQueue.Publish(do.Ok(&rpc.LogEvent{Type: []byte(kind), Message: []byte(message)}))
}

// IsRegistered returns true if the dark node is registered for the current epoch
func (node *DarkNode) IsRegistered() bool {
	registered, err := node.Registrar.IsDarkNodeRegistered(node.Configuration.MultiAddress.ID())
	if err != nil {
		return false
	}
	return registered
}

// IsPendingRegistration returns true if the dark node will be registered in the next epoch
func (node *DarkNode) IsPendingRegistration() bool {
	registered, err := node.Registrar.IsDarkNodePendingRegistration(node.Configuration.MultiAddress.ID())
	if err != nil {
		return false
	}
	return registered
}

// Register the node on the registrar smart contract .
func (node *DarkNode) Register() error {
	registered := node.IsRegistered()
	if registered {
		return nil
	}
	publicKey := append(node.Configuration.RepublicKeyPair.PublicKey.X.Bytes(), node.Configuration.RepublicKeyPair.PublicKey.Y.Bytes()...)
	_, err := node.Registrar.Register(node.Configuration.MultiAddress.ID(), publicKey)
	if err != nil {
		return err
	}
	err = node.Registrar.WaitTillRegistration(node.Configuration.MultiAddress.ID())
	return err
}

// Deregister the node on the registrar smart contract
func (node *DarkNode) Deregister() error {
	_, err := node.Registrar.Deregister(node.Configuration.MultiAddress.ID())
	if err != nil {
		return err
	}
	return nil
}

// PingDarkPool call rpc.PingTarget on each node in a dark pool
func (node *DarkNode) PingDarkPool(ids dnr.IDDarkPool) (dnr.DarkPool, dnr.IDDarkPool) {

	darkpool := make(identity.MultiAddresses, 0)
	disconnectedDarkPool := make(dnr.IDDarkPool, 0)

	for _, id := range ids {
		target, err := node.Swarm.FindNode(id)
		if err != nil || target == nil {
			log.Printf("%v couldn't find pool peer %v: %v", node.Configuration.MultiAddress.Address(), id, err)
			disconnectedDarkPool = append(disconnectedDarkPool, id)
			continue
		}

		darkpool = append(darkpool, *target)

		err = rpc.PingTarget(*target, node.Swarm.MultiAddress(), defaultTimeout)
		if err != nil {
			log.Printf("%v couldn't ping pool peer %v: %v", node.Configuration.MultiAddress.Address(), target, err)
			continue
		}

		err = node.Swarm.DHT.UpdateMultiAddress(*target)
		if err != nil {
			log.Printf("%v coudln't update DHT for pool peer %v: %v", node.Configuration.MultiAddress.Address(), target, err)
			continue
		}
	}
	return dnr.DarkPool{Nodes: darkpool}, disconnectedDarkPool
}

// RepingDarkPool will continually attempt to connect to a set of nodes
// in a darkpool until they are all connected
// Call in a goroutine
func (node *DarkNode) RepingDarkPool(ids dnr.IDDarkPool) {
	currentBlockhash := node.EpochBlockhash

	for len(ids) > 0 {
		if node.EpochBlockhash != currentBlockhash {
			return
		}
		log.Printf("Attempting to re-ping nodes!!!: %v", ids)
		i := 0
		for i < len(ids) {
			id := ids[i]
			target, err := node.Swarm.FindNode(id)
			if err != nil || target == nil {
				log.Printf("%v couldn't find pool peer %v: %v", node.Configuration.MultiAddress.Address(), id, err)
				// We couldn't find this node so we move on to the next one
				i++
				continue
			}

			node.DarkPool.Add(*target)

			// Remove id from disconnected ids
			ids[i] = ids[len(ids)-1]
			ids = ids[:len(ids)-1]
			// Because ids is now shorter, we don't increment i

			err = rpc.PingTarget(*target, node.Swarm.MultiAddress(), defaultTimeout)
			if err != nil {
				log.Printf("%v couldn't ping pool peer %v: %v", node.Configuration.MultiAddress.Address(), target, err)
				continue
			}

			err = node.Swarm.DHT.UpdateMultiAddress(*target)
			if err != nil {
				log.Printf("%v coudln't update DHT for pool peer %v: %v", node.Configuration.MultiAddress.Address(), target, err)
				continue
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// WatchForEpoch will check if a new epoch has been triggered and then sleep for 5 minutes
// Should be called in a goroutine
func (node *DarkNode) WatchForEpoch() {
	for {
		epoch, err := node.Registrar.CurrentEpoch()
		if err != nil {
			log.Printf("%v errored when checking epoch: %v", node.Configuration.MultiAddress.Address(), err)
		}

		if epoch.Blockhash != node.EpochBlockhash {
			log.Printf("%v new epoch!", node.Configuration.MultiAddress.Address())
			node.EpochBlockhash = epoch.Blockhash
			node.AfterEachEpoch()
		}
		time.Sleep(5 * 60 * time.Second)
	}
}

// AfterEachEpoch should be run after each new epoch
func (node *DarkNode) AfterEachEpoch() error {
	log.Printf("%v is pinging dark pool\n", node.Configuration.MultiAddress.Address())

	darkOceanOverlay, err := node.Registrar.GetDarkPools()
	if err != nil {
		log.Fatalf("%v couldn't get dark pools: %v", node.Configuration.MultiAddress.Address(), err)
	}
	node.DarkOceanOverlay = darkOceanOverlay

	idPool, err := node.DarkOceanOverlay.FindDarkPool(node.Configuration.MultiAddress.ID())
	if err != nil {
		return err
	}

	connectedDarkPool, disconnectedDarkPool := node.PingDarkPool(idPool)
	node.DarkPool = connectedDarkPool

	log.Printf("%v connected to dark pool: %v", node.Configuration.MultiAddress.Address(), node.DarkPool)

	go node.RepingDarkPool(disconnectedDarkPool)

	return nil
}

// ConnectToRegistrar will connect to the registrar using the given private key to sign transactions
func ConnectToRegistrar(keypair identity.KeyPair) (*dnr.DarkNodeRegistrar, error) {
	// todo : hard code the ciphertext for now
	auth := bind.NewKeyedTransactor(keypair.PrivateKey)
	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)
	client := dnr.Ropsten("https://ropsten.infura.io/")
	contractAddress := common.HexToAddress("0x6e48bdd8949d0c929e9b5935841f6ff18de0e613")
	renContract := common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a")
	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), &client, auth, &bind.CallOpts{}, contractAddress, renContract, nil)
	return userConnection, nil
}

func (node *DarkNode) Usage() {
	// memory
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		node.Error(err)
	}
	node.Info(fmt.Sprintf("%d", vmStat.Used))
	log.Print("mem : ", vmStat.Used)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	if err != nil {
		node.Error(err)
	}
	node.Info(fmt.Sprintf("%d", cpuStat[0].CacheSize))
	log.Print("cpu : ", cpuStat[0].CacheSize)

}
