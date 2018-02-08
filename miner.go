package miner

import (
	"log"
	"math/big"
	"net"
	"runtime"
	"time"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
	"github.com/republicprotocol/go-xing"
	"google.golang.org/grpc"
)

// TODO: These variables should come from the Registrar contract.
var (
	N        = int64(5)
	K        = int64(3)
	Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
)

type Miner struct {
	Computer      *compute.ComputationMatrix
	Server        *grpc.Server
	Swarm         *swarm.Node
	Xing          *xing.Node
	Configuration *Config
	quit          chan struct{}
}

// NewMiner creates a new Miner, a new swarm.Node and xing.Node and assigns the
// new Miner as the delegate for both. Returns the new Miner, or an error.
func NewMiner(config *Config) (*Miner, error) {
	miner := &Miner{
		Computer:      compute.NewComputationMatrix(),
		Server:        grpc.NewServer(grpc.ConnectionTimeout(time.Minute)),
		Configuration: config,
		quit:          make(chan struct{}),
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
	swarmNode := swarm.NewNode(miner.Server, miner, swarmOptions)
	miner.Swarm = swarmNode

	xingOptions := xing.Options{
		Address:        config.MultiAddress.Address(),
		Debug:          xing.DebugHigh,
		Timeout:        30 * time.Second,
		TimeoutStep:    30 * time.Second,
		TimeoutRetries: 3,
		Concurrent:     true,
	}
	xingNode := xing.NewNode(miner.Server, miner, xingOptions)
	miner.Xing = xingNode

	return miner, nil
}

// Start mining for compute.Orders that are matched. It establishes connections
// to other peers in the swarm network by bootstrapping against a set of
// bootstrap swarm.Nodes.
func (miner *Miner) Start() {
	// Start both gRPC servers.
	miner.Swarm.Register()
	miner.Xing.Register()
	go func() {
		listener, err := net.Listen("tcp", miner.Configuration.Host+":"+miner.Configuration.Port)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on %s:%s\n", miner.Configuration.Host, miner.Configuration.Port)
		if err := miner.Server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the server to start and bootstrap the connections in the swarm.
	time.Sleep(time.Second)
	miner.Swarm.Bootstrap()

	// Start the compute.Order processing loop.
	loop := make(chan struct{}, 1)
	loop <- struct{}{}
	for {
		select {
		case <-miner.quit:
			miner.Server.Stop()
			return
		case <-loop:
			go func() {
				miner.ComputeAll()
				loop <- struct{}{}
			}()
		}
	}
}

// Stop mining.
func (miner *Miner) Stop() {
	miner.quit <- struct{}{}
}

// ComputeAll compute.Computations that can be handled by this Miner in a
// single step. It will process one compute.Computation per CPU available and
// broadcast the compute.ResultFragments.
func (miner Miner) ComputeAll() {
	numberOfCPUs := runtime.NumCPU()
	computations := miner.Computer.WaitForComputations(numberOfCPUs)
	resultFragments := make([]*compute.ResultFragment, len(computations))

	do.CoForAll(computations, func(i int) {
		resultFragment, err := miner.Compute(computations[i])
		if err != nil {
			return
		}
		resultFragments[i] = resultFragment
	})

	go func() {
		resultFragmentsOk := make([]*compute.ResultFragment, 0, len(resultFragments))
		for _, resultFragment := range resultFragments {
			if resultFragment != nil {
				resultFragmentsOk = append(resultFragmentsOk, resultFragment)
			}
		}
		miner.addResultFragments(resultFragmentsOk)
	}()
}

// Compute the required computation on two OrderFragments and send the result
// to all Miners in the M Network.
// TODO: Send computed order fragments to the M Network instead of all peers.
func (miner Miner) Compute(computation *compute.Computation) (*compute.ResultFragment, error) {
	resultFragment, err := computation.Sub(Prime)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, multiAddress := range miner.Swarm.DHT.MultiAddresses() {
			rpc.SendResultFragmentToTarget(multiAddress, multiAddress.Address(), miner.Swarm.MultiAddress(), resultFragment, miner.Swarm.Options.Timeout)
		}
	}()
	return resultFragment, nil
}

func (miner Miner) addResultFragments(resultFragments []*compute.ResultFragment) {
	results, _ := miner.Computer.AddResultFragments(resultFragments, K, Prime)
	for _, result := range results {
		if result.IsMatch(Prime) {
			log.Printf("%v computed [%s, %s] = match!\n", miner.Swarm.Address(), base58.Encode(result.BuyOrderID), base58.Encode(result.SellOrderID))
		} else {
			log.Printf("%v computed [%s, %s] = mismatch!\n", miner.Swarm.Address(), base58.Encode(result.BuyOrderID), base58.Encode(result.SellOrderID))
		}
	}
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
	miner.Computer.AddOrderFragment(orderFragment)
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
