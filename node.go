package node

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-dark-node-registrar"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	rpc "github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
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
	Swarm         *swarm.Node
	Dark          *dark.Node
	Configuration *Config
	Registrar     *dnr.DarkNodeRegistrar

	DeltaBuilder        *compute.DeltaBuilder
	DeltaFragmentMatrix *compute.DeltaFragmentMatrix
	DarkPool            identity.MultiAddresses
	DarkPoolLimit       int64

	quitServer chan struct{}
	quitPacker chan struct{}

	logQueue *LogQueue
}

// NewDarkNode creates a new DarkNode, a new swarm.Node and dark.Node and assigns the
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

	registrar, err := ConnectToRegistrar()
	if err != nil {
		return nil, err
	}

	node := &DarkNode{
		Server:        grpc.NewServer(grpc.ConnectionTimeout(time.Minute)),
		Configuration: config,
		Registrar:     registrar,
		quitServer:    make(chan struct{}),
		logQueue:      NewLogQueue(),
	}
	node.DarkPool = config.BootstrapMultiAddresses
	node.DarkPoolLimit = int64((len(node.DarkPool)*2 + 1) / 3)
	node.DeltaBuilder = compute.NewDeltaBuilder(node.DarkPoolLimit, config.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(config.Prime)

	swarmOptions := swarm.Options{
		MultiAddress:            config.MultiAddress,
		BootstrapMultiAddresses: config.BootstrapMultiAddresses,
		Debug:           swarm.DebugMedium,
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
		Debug:          dark.DebugMedium,
		Timeout:        30 * time.Second,
		TimeoutStep:    30 * time.Second,
		TimeoutRetries: 3,
		Concurrent:     true,
	}
	darkNode := dark.NewNode(node.Server, node, darkOptions)
	node.Dark = darkNode

	return node, nil
}

// OnOrderFragmentForwarding ...
func (node *DarkNode) OnOrderFragmentForwarding(address identity.Address, peer identity.MultiAddress, fragment *compute.OrderFragment) {
	// TODO: Log metrics for the deep query.
}

// OnDeltaFragmentForwarding ...
func (node *DarkNode) OnDeltaFragmentForwarding(address identity.Address, peer identity.MultiAddress, fragment *compute.DeltaFragment) {
	// TODO: Log metrics for the deep query.
}

// OnDeltaFragmentReceived ...
func (node *DarkNode) OnDeltaFragmentReceived(peer identity.MultiAddress, fragment *compute.DeltaFragment) {
	// TODO: Log metrics for the deep query.
}

// Start mining for compute.Orders that are matched. It establishes connections
// to other peers in the swarm network by bootstrapping against a set of
// bootstrap swarm.Nodes.
func (node *DarkNode) Start() error {

	//isRegistered := node.IsRegistered()
	//if !isRegistered {
	//	log.Println("You are not registered")
	//	err := node.Register()
	//	if err != nil {
	//		return err
	//	}
	//}

	//darkPool, err := node.Registrar.GetDarkpool()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//darkPool := getDarkPool()

	// Wait for the server to start and bootstrap the connections in the swarm.
	node.Swarm.Bootstrap()

	//for _, id := range darkPool {
	//	target, err := node.Swarm.FindNode(id[:])
	//	if err != nil {
	//		return err
	//	}
	//	// Ignore the node if we can't find it
	//	if target == nil {
	//		continue
	//	}
	//	err = rpc.PingTarget(*target, node.Swarm.MultiAddress(), 5*time.Second)
	//	// Update the nodes in our DHT if they respond
	//	if err == nil {
	//		node.DarkPool = append(node.DarkPool, *target)
	//		node.Swarm.DHT.UpdateMultiAddress(*target)
	//	}
	//}

	////  Ping all nodes in the dark pool
	//for _, id := range darkPool {
	//	target, err := node.Swarm.FindNode(id[:])
	//	if err != nil {
	//		return err
	//	}
	//	// Ignore the node if we can't find it
	//	if target == nil {
	//		continue
	//	}
	//	err = rpc.PingTarget(*target, node.Swarm.MultiAddress(), 5*time.Second)
	//	// Update the nodes in our DHT if they respond
	//	if err == nil {
	//		node.DarkPool = append(node.DarkPool, *target)
	//		node.Swarm.DHT.UpdateMultiAddress(*target)
	//	}
	//}

	<-node.quitServer

	return nil
}

// Stop mining.
func (node *DarkNode) Stop() {
	node.quitServer <- struct{}{}
	node.quitPacker <- struct{}{}
}

// StartListening starts listening for rpc calls
func (node *DarkNode) StartListening() error {
	log.Printf("Listening on %s:%s\n", node.Configuration.Host, node.Configuration.Port)
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

func (node *DarkNode) IsRegistered() bool {
	registered, err := node.Registrar.IsDarkNodeRegistered(node.Configuration.MultiAddress.ID())
	log.Println("is registered ?", registered, "error:", err)
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
	tx, err := node.Registrar.Register(node.Configuration.MultiAddress.ID(), publicKey)
	log.Println(tx, err)
	if err != nil {
		return err
	}
	err = node.Registrar.WaitTillRegistration(node.Configuration.MultiAddress.ID())
	return err
}

func getDarkPool() []identity.ID {
	ids := make([]identity.ID, 8)
	for i := 0; i < 8; i++ {
		config, _ := LoadConfig(fmt.Sprintf("./test_configs/config-%d.json", i))
		ids[i] = config.MultiAddress.ID()
	}
	return ids
}

func ConnectToRegistrar() (*dnr.DarkNodeRegistrar, error) {
	// todo : hard code the ciphertext for now
	key := `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		return nil, err
	}
	client := dnr.Ropsten("https://ropsten.infura.io/")
	contractAddress := common.HexToAddress("0x0B1148699C93cA9Cfa28f11BD581936f673F76ec")
	renContract := common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a")
	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), &client, auth, &bind.CallOpts{}, contractAddress, renContract, nil)
	return userConnection, nil
}
