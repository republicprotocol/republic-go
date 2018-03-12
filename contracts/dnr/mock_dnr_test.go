package dnr_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/dark-ocean"
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
	mockDnr := &MockDNR{
		registered:   make([][]byte, 0),
		toRegister:   make([][]byte, 0),
		toDeregister: make([][]byte, 0),
	}
	mockDnr.Epoch()

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
				mockDnr.Register(
					configs[i].MultiAddress.ID(),
					append(configs[i].RepublicKeyPair.PublicKey.X.Bytes(), configs[i].RepublicKeyPair.PublicKey.Y.Bytes()...),
				)
				ethAddresses[i] = bind.NewKeyedTransactor(configs[i].EthereumKey.PrivateKey)
				nodes[i] = NewTestDarkNode(mockDnr, *configs[i])
			}

			mockDnr.Epoch()
			startListening(nodes)
		})

		AfterEach(func() {
			mu.Unlock()
		})

		It("WatchForDarkOceanChanges sends a new DarkOceanOverlay on a channel whenever the epoch changes", func() {
			channel := make(chan do.Option, 1)
			go darkocean.WatchForDarkOceanChanges(mockDnr, channel)
			mockDnr.Epoch()
			Eventually(channel).Should(Receive())
		})

		It("Registration checking returns the correct result", func() {
			id0 := nodes[0].MultiAddress.ID()
			pub := append(nodes[0].Config.RepublicKeyPair.PublicKey.X.Bytes(), nodes[0].Config.RepublicKeyPair.PublicKey.Y.Bytes()...)
			mockDnr.Deregister(id0)

			// Before epoch, should still be registered
			Ω(nodes[0].IsRegistered()).Should(Equal(true))
			Ω(nodes[0].IsPendingRegistration()).Should(Equal(false))

			mockDnr.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(false))

			mockDnr.Register(id0, pub)

			// Before epoch, should still be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(false))
			Ω(nodes[0].IsPendingRegistration()).Should(Equal(true))

			mockDnr.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].IsRegistered()).Should(Equal(true))
		})
	})
})

/*














 * Mock DarkNodeRegistar implementation

 */

type MockDNR struct {
	registered   [][]byte
	toRegister   [][]byte
	toDeregister [][]byte
	epoch        dnr.Epoch
}

func (mockDnr *MockDNR) Register(_darkNodeID []byte, _publicKey []byte) (*types.Transaction, error) {
	isRegistered, _ := mockDnr.IsDarkNodeRegistered(_darkNodeID)
	isPending, _ := mockDnr.IsDarkNodePendingRegistration(_darkNodeID)
	if isRegistered || isPending {
		return nil, errors.New("Must not be registered to register")
	}
	mockDnr.toRegister = append(mockDnr.toRegister, _darkNodeID)
	return nil, nil
}
func (mockDnr *MockDNR) Deregister(_darkNodeID []byte) (*types.Transaction, error) {
	for i, id := range mockDnr.toRegister {
		if string(_darkNodeID) == string(id) {
			mockDnr.toDeregister[i] = mockDnr.toDeregister[len(mockDnr.toDeregister)-1]
			mockDnr.toDeregister = mockDnr.toDeregister[:len(mockDnr.toDeregister)-1]
			return nil, nil
		}
	}
	if isRegistered, _ := mockDnr.IsDarkNodeRegistered(_darkNodeID); !isRegistered {
		return nil, errors.New("Must be registered to deregister")
	}
	mockDnr.toDeregister = append(mockDnr.toRegister, _darkNodeID)
	return nil, nil
}
func (mockDnr *MockDNR) GetBond(_darkNodeID []byte) (*big.Int, error) {
	return big.NewInt(1000), nil
}
func (mockDnr *MockDNR) IsDarkNodeRegistered(_darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.registered {
		if string(_darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}
func (mockDnr *MockDNR) IsDarkNodePendingRegistration(_darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.toRegister {
		if string(_darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}
func (mockDnr *MockDNR) CurrentEpoch() (dnr.Epoch, error) {
	return mockDnr.epoch, nil
}
func (mockDnr *MockDNR) Epoch() (*types.Transaction, error) {
	var b32 [32]byte

	_, err := rand.Read(b32[:])
	if err != nil {
		return nil, err
	}

	mockDnr.epoch = dnr.Epoch{
		Blockhash: b32,
		Timestamp: big.NewInt(time.Now().Unix()),
	}

	// Remove toRegister nodes
	for _, deregNode := range mockDnr.toDeregister {
		for i, node := range mockDnr.registered {
			if string(node) == string(deregNode) {
				mockDnr.registered[i] = mockDnr.registered[len(mockDnr.registered)-1]
				mockDnr.registered = mockDnr.registered[:len(mockDnr.registered)-1]
				break
			}
		}
	}

	mockDnr.registered = append(mockDnr.registered, mockDnr.toRegister...)

	mockDnr.toDeregister = make([][]byte, 0)
	mockDnr.toRegister = make([][]byte, 0)

	return nil, nil
}
func (mockDnr *MockDNR) GetCommitment(_darkNodeID []byte) ([32]byte, error) {
	var nil32 [32]byte
	return nil32, nil
}
func (mockDnr *MockDNR) GetOwner(_darkNodeID []byte) (common.Address, error) {
	var nil20 [20]byte
	return nil20, nil
}
func (mockDnr *MockDNR) GetPublicKey(_darkNodeID []byte) ([]byte, error) {
	return nil, nil
}
func (mockDnr *MockDNR) GetAllNodes() ([][]byte, error) {
	return mockDnr.registered, nil
}
func (mockDnr *MockDNR) MinimumBond() (*big.Int, error) {
	return big.NewInt(1000), nil
}
func (mockDnr *MockDNR) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(0), nil
}
func (mockDnr *MockDNR) Refund(_darkNodeID []byte) (*types.Transaction, error) {
	return nil, nil
}
func (mockDnr *MockDNR) WaitTillRegistration(_darkNodeID []byte) error {
	_, err := mockDnr.Epoch()
	return err
}

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
		Options: network.Options{
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
func NewTestDarkNode(registrar dnr.DarkNodeRegistrarInterface, config node.Config) *node.DarkNode {
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
	newNode.Swarm = network.NewSwarmService(newNode, newNode.Options, newNode.Logger, newNode.ClientPool, newNode.DHT)
	newNode.Dark = network.NewDarkService(newNode, newNode.Options, newNode.Logger)

	newNode.Registrar = registrar

	return newNode
}
