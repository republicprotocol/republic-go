package node

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-dark-network"
	dnr "github.com/republicprotocol/go-dark-node-registrar"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
	"google.golang.org/grpc"
)

const defaultTimeout = 5 * time.Second

// To be retrieved from the Registrar contract
var (
	// N is the number of dark nodes in the network
	N = int64(5)
	// K is the number of fragments required to reconstruct the secret
	K = int64(3)
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
	logQueue.channels = make([]chan do.Option, 128)
	return logQueue
}

// Publish allows a node to push a log to each client
func (logQueue *LogQueue) Publish(val do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()

	var logQueueLength = len(logQueue.channels)
	for i := 0; i < logQueueLength; i++ {
		timer := time.Tick(10 * time.Second)
		select {
		case logQueue.channels[i] <- val:
		case <-timer:
			// TODO: deregister the channel
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
	DarkPoolLimit       int

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
		DeltaBuilder:        compute.NewDeltaBuilder(int64(config.ComputationShardSize), config.Prime),
		DeltaFragmentMatrix: &compute.DeltaFragmentMatrix{},
		Server:              grpc.NewServer(grpc.ConnectionTimeout(time.Minute)),
		Configuration:       config,
		Registrar:           registrar,

		quitServer: make(chan struct{}),

		logQueue: NewLogQueue(),
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
	// todo:
	node.DarkPool = config.BootstrapMultiAddresses
	node.DarkPoolLimit = len(node.DarkPool)

	return node, nil
}

// OnSync ...
func (node *DarkNode) OnSync(identity.MultiAddress) chan do.Option {
	// TODO: ...
	panic("uninmplemented")
}

// OnLogs returns a channel for receiving logs from the node
func (node *DarkNode) OnLogs(logs chan do.Option) {
	logChannel := make(chan do.Option, 128)
	node.logQueue.Subscribe(logChannel)
	return logChannel
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

	// TODO
	//darkPool, err := node.Registrar.GetDarkpool()
	//if err != nil {
	//	log.Fatal(err)
	//}
	darkPool := getDarkPool()

	// Wait for the server to start and bootstrap the connections in the swarm.
	node.Swarm.Bootstrap()

	//  Ping all nodes in the dark pool
	for _, id := range darkPool {
		target, err := node.Swarm.FindNode(id[:])
		if err != nil {
			return err
		}
		// Ignore the node if we can't find it
		if target == nil {
			continue
		}
		err = rpc.PingTarget(*target, node.Swarm.MultiAddress(), 5*time.Second)
		// Update the nodes in our DHT if they respond
		if err == nil {
			node.DarkPool = append(node.DarkPool, *target)
			node.Swarm.DHT.UpdateMultiAddress(*target)
		}
	}

	// TODO: Synchronize the hidden order book from other DarkNodes.
	//node.StartSynchronization()

	// Begin electing shards in the background.
	//node.StartShardElections()

	return nil
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
	return nil
}

// Stop mining.
func (node *DarkNode) Stop() {
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

// OnOrderFragmentReceived implements the dark.Delegate interface. It is called
// by the underlying dark.Node whenever the Miner receives a
// compute.OrderFragment that it must process.
func (node *DarkNode) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	node.logQueue.Publish(do.Ok{nil})
	deltaFragments, err := node.DeltaFragmentMatrix.InsertOrderFragment(orderFragment)
	if err != nil {
		log.Println(err)
	}
	for _, multiAddress := range node.DarkPool {
		client, err := rpc.NewClient(multiAddress, node.Swarm.MultiAddress())
		if err != nil {
			log.Println(err)
			continue
		}
		for _, deltaFragment := range deltaFragments {
			err := client.BroadcastDeltaFragment(deltaFragment)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (node *DarkNode) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
	delta, err := node.DeltaBuilder.InsertDeltaFragment(deltaFragment)
	if err != nil {
		log.Println(err)
		return
	}
	if delta == nil {
		return
	}
	if delta.IsMatch() {
		log.Printf("[%v] match found (%v, %v)\n", node.Swarm.Address(), base58.Encode(deltaFragment.BuyOrderID), base58.Encode(deltaFragment.SellOrderID))
		// TODO: Attempt to get consensus on the match and then mark the orders
		// handled if the consensus is won. If the consensus is not won take
		// either the buy, or sell (or both), orders and mark them as completed
		// (this depends on which ones conflicted).
	}
}

// OnSync ...
func (node *DarkNode) OnSync(from identity.MultiAddress) chan do.Option {
	// TODO: ...
	panic("uninmplemented")
}

// OnElectShard is a delegate method that is called when the DarkNode has
// received an RPC for electing a Shard. To finish the election, the DarkNode
// should filter the Shard is received and return it. The DarkNode should
// filter out the deltas and residues that it does not have access to.
func (node *DarkNode) OnElectShard(from identity.MultiAddress, shard compute.Shard) compute.Shard {
	// TODO: Elect the shard. Check which deltas and residues the DarkNode has
	// and remove all others from the shard before returning it.

	returnedShard := compute.Shard{
		Signature: shard.Signature,
	}

	pendingDeltas := node.HiddenOrderBook.PendingDeltaFragments()
	pendingDeltaMap := map[string]bool{}
	for i := range pendingDeltas {
		pendingDeltaMap[string(pendingDeltas[i].ID)] = true
	}

	for i := range shard.Deltas {
		if pendingDeltaMap[string(shard.Deltas[i].ID)] {
			returnedShard.Deltas = append(returnedShard.Deltas, shard.Deltas[i])
		}
	}

	for i := range shard.Residues {
		if false {
			returnedShard.Residues = append(returnedShard.Residues, returnedShard.Residues[i])
		}
	}

	return returnedShard
}

// OnComputeShard is a delegate method that is called when the DarkNode has
// received an RPC for computing a Shard. The process for computing a Shard is
// to compute ([[a]] - [[b]]) * [[r^2]] using the ([[a]] - [[b]]) deltas and
// the [[r^2]] residues in the Shard.
func (node *DarkNode) OnComputeShard(from identity.MultiAddress, shard compute.Shard) {
	node.Compute(shard)
}

// OnFinalizeShard is a delegate method that is called when the DarkNode has
// received an RPC for finalizing a Shard. The process for finalizing a Shard
// is to store all computation shares from within the different Shards that are
// finalized. Eventually, enough computation shares will be acquired and the
// computation proper can be reconstructed.
func (node *DarkNode) OnFinalizeShard(from identity.MultiAddress, deltaShard compute.DeltaShard) {
	for i := range deltaShard.DeltaFragments {
		delta, err := node.DeltaBuilder.InsertDeltaFragment(deltaShard.DeltaFragments[i])
		if err != nil {
			log.Println(err)
		}
		if delta != nil {
			if !delta.IsMatch(node.Configuration.Prime) {
				log.Printf("Mismatch [%v]\n\t(%v, %v)", base58.Encode(delta.ID), base58.Encode(delta.BuyOrderID), base58.Encode(delta.SellOrderID))
				continue
			}
			log.Printf("Match [%v]\n\t(%v, %v)", base58.Encode(delta.ID), base58.Encode(delta.BuyOrderID), base58.Encode(delta.SellOrderID))

			// FIXME: After reconstruction of a matched delta, it should be
			// posted to the Ethereum consensus smart contract. Traders can use
			// this to be notified of the match.
		}
	}
}

// StartShardElections will continue to create Shards and run elections for
// them with other DarkNodes.
func (node *DarkNode) StartShardElections() {
	go func() {
		running := int64(1)

		shardQueue := make(chan compute.Shard)
		go func() {
			for atomic.LoadInt64(&running) != 0 {
				shardQueue <- node.HiddenOrderBook.WaitForShard()
			}
		}()

		timer := time.NewTicker(time.Duration(node.Configuration.ComputationShardInterval) * time.Second)
		defer timer.Stop()

		for atomic.LoadInt64(&running) != 0 {
			select {
			case shard := <-shardQueue:
				go func() {
					shard := node.RunShardElection(shard)
					node.RunShardComputation(shard)
				}()
			case <-timer.C:
				preemptedShard := node.HiddenOrderBook.PreemptShard()
				if len(preemptedShard.Deltas) > 0 {
					go func() {
						shard := node.RunShardElection(preemptedShard)
						node.RunShardComputation(shard)
					}()
				}
			case <-node.quitPacker:
				atomic.StoreInt64(&running, 0)
			}
		}
	}()
}

// RunShardElection by calling the ElectShard RPC on all DarkNodes in the dark
// pool. Collect all responses and create a shard of deltas and residues that
// are held by enough DarkNodes that a computation has a chance of succeeding.
func (node *DarkNode) RunShardElection(shard compute.Shard) compute.Shard {

	deltaVotesMu := new(sync.Mutex)
	deltaVotes := map[string]int{}

	residueVotesMu := new(sync.Mutex)
	residueVotes := map[string]int{}

	do.ForAll(node.DarkPool, func(i int) {
		target := node.DarkPool[i]

		newShard, err := rpc.StartElectShard(target, node.Swarm.MultiAddress(), shard, 5*time.Second)
		if err != nil {
			log.Println(err)
			return
		}

		deltaVotesMu.Lock()
		defer deltaVotesMu.Unlock()

		residueVotesMu.Lock()
		defer residueVotesMu.Unlock()

		for j := range newShard.Deltas {
			deltaVotes[string(newShard.Deltas[j])]++
		}

		for j := range newShard.Residues {
			residueVotes[string(newShard.Residues[j])]++
		}
	})

	returnedShard := compute.Shard{
		Signature: shard.Signature,
	}
	for i := range shard.Deltas {
		delta := shard.Deltas[i]
		votes := deltaVotes[string(delta.ID)]
		if votes > node.Configuration.ComputationShardSize {
			returnedShard.Deltas = append(returnedShard.Deltas, delta)
		}
	}

	for i := range shard.Residues {
		residue := shard.Residues[i]
		votes := residueVotes[string(residue.ID)]
		if votes > node.Configuration.ComputationShardSize {
			returnedShard.Residues = append(returnedShard.Residues, residue)
		}
	}

	return returnedShard
}

// RunShardComputation by calling the ComputeShard RPC on all DarkNodes in the
// dark pool.
func (node *DarkNode) RunShardComputation(shard compute.Shard) {
	go node.Compute(shard)
	do.ForAll(node.DarkPool, func(i int) {
		// TODO: Call the ComputeShard RPC on all peers in the dark pool.
		rpc.AskToComputeShard(node.DarkPool[i], node.Swarm.MultiAddress(), shard, defaultTimeout)
	})
}

// Compute ...
func (node *DarkNode) Compute(shard compute.Shard) {
	//finalShard := shard.Compute(node.Configuration.Prime)
	finalShard := shard.Compute()
	do.ForAll(node.DarkPool, func(i int) {
		rpc.FinalizeShard(node.DarkPool[i], node.Swarm.MultiAddress(), finalShard, defaultTimeout)
	})
}

func ConnectToRegistrar() (*dnr.DarkNodeRegistrar, error) {
	// todo : hard code the ciphertext for now
	key := `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		return nil, err
	}
	client := dnr.Ropsten("https://ropsten.infura.io/")
	contractAddress := common.HexToAddress("0xF874c2b8Afaa199A81796746280Af9184cd0D75b")
	renContract := common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a")
	userConnection := dnr.NewDarkNodeRegistrar(context.Background(), &client, auth, &bind.CallOpts{}, contractAddress, renContract, nil)
	return userConnection, nil
}

func getDarkPool() []identity.ID {
	ids := make([]identity.ID, 8)
	for i := 0; i < 8; i++ {
		config, _ := LoadConfig(fmt.Sprintf("./test_configs/config-%d.json", i))
		ids[i] = config.MultiAddress.ID()
	}
	return ids
}
