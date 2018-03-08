package node_test

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/contracts/connection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/compute"
	darknode "github.com/republicprotocol/republic-go/dark-node"
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
	var mu = new(sync.Mutex)
	var nodes []*darknode.DarkNode
	var configs []*darknode.Config
	var ethAddresses []*bind.TransactOpts
	// var err error

	startListening := func(nodes []*darknode.DarkNode, bootstrapNodes int) {
		// Fully connect the bootstrap nodes
		for i := 0; i < bootstrapNodes; i++ {
			for j := 0; j < bootstrapNodes; j++ {
				if i == j {
					continue
				}
				nodes[j].ClientPool.Ping(nodes[i].MultiAddress)
			}
		}
	}

	Context("nodes start up", func() {
		BeforeEach(func() {
			mu.Lock()

			configs = make([]*darknode.Config, NumberOfTestNODES)
			ethAddresses = make([]*bind.TransactOpts, NumberOfTestNODES)
			nodes = make([]*darknode.DarkNode, NumberOfTestNODES)

			for i := 0; i < NumberOfTestNODES; i++ {
				configs[i] = MockConfig()
				keypair, err := configs[i].EthereumKeyPair()
				Ω(err).ShouldNot(HaveOccurred())
				ethAddresses[i] = bind.NewKeyedTransactor(keypair.PrivateKey)
			}

			clientDetails, err := connection.Simulated(ethAddresses...)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < NumberOfTestNODES; i++ {
				nodes[i] = NewTestDarkNode(clientDetails, *configs[i])
			}
			// nodes, err = generateNodes(NumberOfBootstrapNodes, NumberOfTestNODES)
			startListening(nodes, NumberOfTestNODES)
		})

		AfterEach(func() {
			mu.Unlock()
		})

		It("should be able to run startup successfully", func() {
		})
	})
})

func MockConfig() *darknode.Config {
	keypair, err := identity.NewKeyPair()
	if err != nil {
		panic(err)
	}

	ethereumPair, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	ethereumKey := hex.EncodeToString(ethereumPair.D.Bytes())

	port := "18514"
	host := "0.0.0.0"

	println(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, keypair.Address()))
	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, keypair.Address()))
	if err != nil {
		panic(err)
	}

	return &darknode.Config{
		RepublicKeyPair:    keypair,
		RSAKeyPair:         keypair,
		EthereumPrivateKey: ethereumKey,
		Port:               port,
		Host:               host,
		MultiAddress:       multiAddress,
	}
}

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewTestDarkNode(clientDetails connection.ClientDetails, config darknode.Config) *darknode.DarkNode {
	if config.Prime == nil {
		config.Prime = darknode.Prime
	}

	// TODO: This should come from the DNR.
	k := int64(5)

	node := &darknode.DarkNode{Config: config}

	node.Logger = logger.NewLogger()
	node.ClientPool = rpc.NewClientPool(node.MultiAddress)
	node.DHT = dht.NewDHT(node.MultiAddress.Address(), node.MaxBucketLength)

	node.DeltaBuilder = compute.NewDeltaBuilder(k, node.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(node.Prime)
	node.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	node.OrderFragmentWorker = darknode.NewOrderFragmentWorker(node.OrderFragmentWorkerQueue, node.DeltaFragmentMatrix)
	node.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentWorker = darknode.NewDeltaFragmentWorker(node.DeltaFragmentWorkerQueue, node.DeltaBuilder)

	// options := network.Options{}
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	node.Swarm = network.NewSwarmService(node, node.Options, node.Logger, node.ClientPool, node.DHT)
	node.Dark = network.NewDarkService(node, node.Options, node.Logger)

	registrar, err := darknode.ConnectToRegistrar(clientDetails, config)
	if err != nil {
		// TODO: Handler err
		panic(err)
	}
	node.Registrar = registrar

	return node
}
