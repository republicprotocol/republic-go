package contracts_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network"
)

const NumberOfTestNODES = 4

var _ = Describe("Dark nodes", func() {

	var mu = new(sync.Mutex)
	var nodes []*node.DarkNode
	var configs []*node.Config
	var ethAddresses []*bind.TransactOpts
	var DNR dnr.DarkNodeRegistry

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

			var err error
			DNR, err = dnr.TestnetDNR(nil)
			Ω(err).ShouldNot(HaveOccurred())

			DNR.WaitForEpoch()

			configs = make([]*node.Config, NumberOfTestNODES)
			ethAddresses = make([]*bind.TransactOpts, NumberOfTestNODES)
			nodes = make([]*node.DarkNode, NumberOfTestNODES)

			for i := 0; i < NumberOfTestNODES; i++ {

				bond := stackint.FromUint(10)
				// mul := stackint.FromUint(20)
				// bond = bond.Exp(&mul)
				configs[i] = MockConfig()
				auth := bind.NewKeyedTransactor(configs[i].EthereumKey.PrivateKey)

				dnr, err := dnr.TestnetDNR(auth)
				Ω(err).ShouldNot(HaveOccurred())

				err = dnr.ApproveRen(&bond)
				Ω(err).ShouldNot(HaveOccurred())

				dnr.SetGasLimit(300000)
				_, err = dnr.Register(
					configs[i].NetworkOptions.MultiAddress.ID(),
					append(configs[i].KeyPair.PublicKey.X.Bytes(), configs[i].KeyPair.PublicKey.Y.Bytes()...),
					&bond,
				)
				Ω(err).ShouldNot(HaveOccurred())
				dnr.SetGasLimit(0)

				ethAddresses[i] = bind.NewKeyedTransactor(configs[i].EthereumKey.PrivateKey)
				nodes[i], err = node.NewDarkNode(*configs[i], dnr)
				Ω(err).ShouldNot(HaveOccurred())
			}

			// Ω(nodes[0].DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))
			// Ω(nodes[0].DarkNodeRegistry.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))

			DNR.WaitForEpoch()

			Ω(nodes[0].DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))
			// Ω(nodes[0].DarkNodeRegistry.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))

			startListening(nodes)
		})

		AfterEach(func() {
			stopListening(nodes)
			mu.Unlock()
		})

		It("Registration checking returns the correct result", func() {
			// pub := append(nodes[0].Config.KeyPair.PublicKey.X.Bytes(), nodes[0].Config.KeyPair.PublicKey.Y.Bytes()...)

			dnr := nodes[0].DarkNodeRegistry
			_, err := dnr.Deregister(nodes[0].NetworkOptions.MultiAddress.ID())
			Ω(err).ShouldNot(HaveOccurred())

			// Before epoch, should still be registered
			Ω(dnr.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(BeTrue())
			// Ω(nodes[0].DarkNodeRegistry.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(false))

			DNR.WaitForEpoch()

			// After epoch, should be deregistered
			Ω(dnr.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(BeFalse())
			isDeregistered, err := dnr.IsDeregistered(nodes[0].NetworkOptions.MultiAddress.ID())
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isDeregistered).Should(BeTrue())

			dnr.SetGasLimit(300000)
			_, err = dnr.Refund(nodes[0].NetworkOptions.MultiAddress.ID())
			Ω(err).ShouldNot(HaveOccurred())
			dnr.SetGasLimit(0)

			bond := stackint.FromUint(10)
			err = dnr.ApproveRen(&bond)
			Ω(err).ShouldNot(HaveOccurred())

			dnr.SetGasLimit(300000)
			_, err = dnr.Register(
				nodes[0].NetworkOptions.MultiAddress.ID(),
				append(nodes[0].KeyPair.PublicKey.X.Bytes(), nodes[0].KeyPair.PublicKey.Y.Bytes()...),
				&bond,
			)
			Ω(err).ShouldNot(HaveOccurred())
			dnr.SetGasLimit(0)

			// Before epoch, should still be deregistered
			Ω(nodes[0].DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(BeFalse())
			// Ω(nodes[0].DarkNodeRegistry.IsDarkNodePendingRegistration(nodes[0].NetworkOptions.MultiAddress.ID())).Should(Equal(true))

			DNR.WaitForEpoch()

			// After epoch, should be deregistered
			Ω(nodes[0].DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())).Should(BeTrue())
		})

		It("can deregister all nodes", func() {
			for _, node := range nodes {
				dnr := node.DarkNodeRegistry

				_, err := dnr.Deregister(node.NetworkOptions.MultiAddress.ID())
				Ω(err).ShouldNot(HaveOccurred())

				// Before epoch, should still be registered
				Ω(dnr.IsRegistered(node.NetworkOptions.MultiAddress.ID())).Should(BeTrue())
				// Ω(node.DarkNodeRegistry.IsDarkNodePendingRegistration(node.NetworkOptions.MultiAddress.ID())).Should(Equal(false))

				DNR.WaitForEpoch()

				// After epoch, should be deregistered
				Ω(dnr.IsRegistered(node.NetworkOptions.MultiAddress.ID())).Should(BeFalse())
				isDeregistered, err := dnr.IsDeregistered(node.NetworkOptions.MultiAddress.ID())
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isDeregistered).Should(BeTrue())

				dnr.SetGasLimit(300000)
				_, err = dnr.Refund(node.NetworkOptions.MultiAddress.ID())
				Ω(err).ShouldNot(HaveOccurred())
				dnr.SetGasLimit(0)

				DNR.WaitForEpoch()
			}
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
		KeyPair:     keypair,
		EthereumKey: *ethereumKey,
		Port:        port,
		Host:        host,
	}
}
