package dnr_test

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network"
)

const NumberOfTestNODES = 4

var _ = Describe("Dark nodes", func() {

	var mu = new(sync.Mutex)
	var err error
	var nodes []*node.DarkNode
	var configs []*node.Config
	var ethAddresses []*bind.TransactOpts
	var MockDarkNodeRegistrar dnr.DarkNodeRegistrar

	startListening := func(nodes []*node.DarkNode) {
		// Fully connect the bootstrap nodes
		for _, n := range nodes {
			go func(n *node.DarkNode) {
				n.StartServices()
			}(n)
		}
		time.Sleep(1 * time.Second)
		for i, iNode := range nodes {
			for j, jNode := range nodes {
				if i == j {
					continue
				}
				// log.Printf("%v pinging %v\n", iNode.MultiAddress.Address(), jNode.MultiAddress.Address())
				jNode.ClientPool.Ping(iNode.NetworkOptions.MultiAddress)
			}
		}
	}

	stopListening := func(nodes []*node.DarkNode) {
		for _, n := range nodes {
			n.Stop()
		}
	}

	Context("nodes start up", func() {
		BeforeEach(func() {
			mu.Lock()

			MockDarkNodeRegistrar, err = dnr.NewMockDarkNodeRegistrar()
			MockDarkNodeRegistrar.Epoch()

			configs = make([]*node.Config, NumberOfTestNODES)
			ethAddresses = make([]*bind.TransactOpts, NumberOfTestNODES)
			nodes = make([]*node.DarkNode, NumberOfTestNODES)

			for i := 0; i < NumberOfTestNODES; i++ {

				configs[i] = MockConfig()
				MockDarkNodeRegistrar.Register(
					configs[i].NetworkOptions.MultiAddress.ID(),
					append(configs[i].RepublicKeyPair.PublicKey.X.Bytes(), configs[i].RepublicKeyPair.PublicKey.Y.Bytes()...),
					big.NewInt(100),
				)
				ethAddresses[i] = bind.NewKeyedTransactor(configs[i].EthereumKey.PrivateKey)
				nodes[i], err = node.NewDarkNode(*configs[i], MockDarkNodeRegistrar)
				Ω(err).ShouldNot(HaveOccurred())
			}

			MockDarkNodeRegistrar.Epoch()
			startListening(nodes)
		})

		AfterEach(func() {
			stopListening(nodes)
			mu.Unlock()
		})

		It("WatchForDarkOceanChanges sends a new DarkOceanOverlay on a channel whenever the epoch changes", func() {
			channel := make(chan do.Option, 1)
			MockDarkNodeRegistrar.Epoch()
			Eventually(channel).Should(Receive())
		})

		It("Registration checking returns the correct result", func() {
			id0 := nodes[0].NetworkOptions.MultiAddress.ID()
			pub := append(nodes[0].Config.RepublicKeyPair.PublicKey.X.Bytes(), nodes[0].Config.RepublicKeyPair.PublicKey.Y.Bytes()...)
			MockDarkNodeRegistrar.Deregister(id0)

			// Before epoch, should still be registered
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodeRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))

			MockDarkNodeRegistrar.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodeRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))

			MockDarkNodeRegistrar.Register(id0, pub, big.NewInt(100))

			// Before epoch, should still be deregistered
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodeRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))

			MockDarkNodeRegistrar.Epoch()

			// After epoch, should be deregistered
			Ω(nodes[0].DarkNodeRegistrar.IsDarkNodeRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))
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
		RepublicKeyPair: &keypair,
		RSAKeyPair:      &keypair,
		EthereumKey:     ethereumKey,
		Port:            port,
		Host:            host,
	}
}
