package dnr_test

import (
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

const NumberOfTestNODES = 4

var _ = Describe("Dark nodes", func() {
	MockDarkNodeRegistrar := &MockDarkNodeRegistrar{
		registered:   make([][]byte, 0),
		toRegister:   make([][]byte, 0),
		toDeregister: make([][]byte, 0),
	}
	MockDarkNodeRegistrar.Epoch()

	var mu = new(sync.Mutex)
	var nodes []*node.DarkNode
	var configs []*node.Config
	var ethAddresses []*bind.TransactOpts

	startListening := func(nodes []*node.DarkNode) {
		// Fully connect the bootstrap nodes
		for i, iNode := range nodes {
			go iNode.Start()
			for j, jNode := range nodes {
				if i == j {
					continue
				}
				// log.Printf("%v pinging %v\n", iNode.MultiAddress.Address(), jNode.MultiAddress.Address())
				jNode.ClientPool.Ping(iNode.MultiAddress)
			}
		}
	}

	Context("nodes start up", func() {
		BeforeEach(func() {
			mu.Lock()

			configs = make([]*node.Config, NumberOfTestNODES)
			ethAddresses = make([]*bind.TransactOpts, NumberOfTestNODES)
			nodes = make([]*node.DarkNode, NumberOfTestNODES)

			for i := 0; i < NumberOfTestNODES; i++ {
				configs[i] = MockConfig()
				MockDarkNodeRegistrar.Register(
					configs[i].MultiAddress.ID(),
					append(configs[i].RepublicKeyPair.PublicKey.X.Bytes(), configs[i].RepublicKeyPair.PublicKey.Y.Bytes()...),
				)
				ethAddresses[i] = bind.NewKeyedTransactor(configs[i].EthereumKey.PrivateKey)
				nodes[i] = NewTestDarkNode(MockDarkNodeRegistrar, *configs[i])
			}

			MockDarkNodeRegistrar.Epoch()
			startListening(nodes)
		})

		AfterEach(func() {
			mu.Unlock()
		})

		It("WatchForDarkOceanChanges sends a new DarkOceanOverlay on a channel whenever the epoch changes", func() {
			channel := make(chan do.Option, 1)
			go darkocean.WatchForDarkOceanChanges(MockDarkNodeRegistrar, channel)
			MockDarkNodeRegistrar.Epoch()
			Eventually(channel).Should(Receive())
		})

		It("Registration checking returns the correct result", func() {
			id0 := nodes[0].MultiAddress.ID()
			pub := append(nodes[0].Config.RepublicKeyPair.PublicKey.X.Bytes(), nodes[0].Config.RepublicKeyPair.PublicKey.Y.Bytes()...)
			MockDarkNodeRegistrar.Deregister(id0)

			// Before epoch, should still be registered
			Ω(nodes[0].IsRegistered()).Should(Equal(true))
			Ω(nodes[0].IsPendingRegistration()).Should(Equal(false))

			MockDarkNodeRegistrar.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(false))

			MockDarkNodeRegistrar.Register(id0, pub)

			// Before epoch, should still be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(false))
			Ω(nodes[0].IsPendingRegistration()).Should(Equal(true))

			MockDarkNodeRegistrar.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(true))
		})
	})
})

var i = 0

func MockConfig() *node.Config {
	keypair, err := identity.NewKeyPair()
	if err != nil {
		panic(err)
	}

	// Long process to get this into the right format!
	ethereumPair, err := crypto.GenerateKey()
	ethereumKey := &keystore.Key{
		Address:    crypto.PubkeyToAddress(ethereumPair.PublicKey),
		PrivateKey: ethereumPair,
	}

	port := fmt.Sprintf("1851%v", i)
	i++
	host := "127.0.0.1"

	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, keypair.Address()))
	if err != nil {
		panic(err)
	}

	return &node.Config{
		NetworkOptions: network.Options{
			MultiAddress: multiAddress,
		},
		RepublicKeyPair: keypair,
		RSAKeyPair:      keypair,
		EthereumKey:     ethereumKey,
		Port:            port,
		Host:            host,
	}
}

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewTestDarkNode(registrar dnr.DarkNodeRegistrar, config node.Config) *node.DarkNode {
	if config.Prime == nil {
		config.Prime = node.Prime
	}

	// TODO: This should come from the DNR.
	k := int64(5)

	newNode := &node.DarkNode{Config: config}

	newNode.Logger = logger.NewLogger()
	newNode.ClientPool = rpc.NewClientPool(newNode.MultiAddress)
	newNode.DHT = dht.NewDHT(newNode.MultiAddress.Address(), newNode.MaxBucketLength)

	newNode.DeltaBuilder = compute.NewDeltaBuilder(k, newNode.Prime)
	newNode.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(newNode.Prime)
	newNode.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	newNode.OrderFragmentWorker = node.NewOrderFragmentWorker(newNode.OrderFragmentWorkerQueue, newNode.DeltaFragmentMatrix)
	newNode.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	newNode.DeltaFragmentWorker = node.NewDeltaFragmentWorker(newNode.DeltaFragmentWorkerQueue, newNode.DeltaBuilder)

	// options := network.Options{}
	newNode.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	newNode.Swarm = network.NewSwarmService(newNode, newNode.NetworkOptions, newNode.Logger, newNode.ClientPool, newNode.DHT)
	newNode.Dark = network.NewDarkService(newNode, newNode.NetworkOptions, newNode.Logger)

	newNode.Registrar = registrar

	return newNode
}
