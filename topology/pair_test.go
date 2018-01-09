package topology_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	identity "github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-swarm/rpc"
)

var _ = Describe("Pair topology", func() {

	It("should route messages in a two-sided connection", func() {
		μ.Lock()
		defer μ.Unlock()

		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		left, err := NewNode(&Config{
			KeyPair:      keyPair,
			MultiAddress: multiAddress,
			Peers:        make(identity.MultiAddresses, 0, numberOfNodes-1),
		})
		Ω(err).ShouldNot(HaveOccurred())

		keyPair, err = identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err = identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3001, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		right, err := NewNode(&Config{
			KeyPair:      keyPair,
			MultiAddress: multiAddress,
			Peers:        make(identity.MultiAddresses, 0, numberOfNodes-1),
		})
		Ω(err).ShouldNot(HaveOccurred())

		go func() {
			defer GinkgoRecover()
			Ω(left.Serve()).ShouldNot(HaveOccurred())
		}()
		go func() {
			defer GinkgoRecover()
			Ω(right.Serve()).ShouldNot(HaveOccurred())
		}()
		time.Sleep(startTimeDelay)

		err = Ping(left.MultiAddress, &rpc.MultiAddress{Multi: right.MultiAddress.String()})
		Ω(err).ShouldNot(HaveOccurred())

		err = Ping(right.MultiAddress, &rpc.MultiAddress{Multi: left.MultiAddress.String()})
		Ω(err).ShouldNot(HaveOccurred())

		Ω(len(left.DHT.MultiAddresses())).Should(Equal(1))
		Ω(len(right.DHT.MultiAddresses())).Should(Equal(1))

		left.Stop()
		right.Stop()
	})

	// It("should route messages in a one-sided connection", func() {

	// // Initialize all nodes.
	// nodes, err := generateNodes()
	// Ω(err).ShouldNot(HaveOccurred())

	// for _, node := range nodes {
	// 	go node.Serve()
	// }
	// time.Sleep(startTimeDelay)

	// // Connect all nodes to each other.
	// for i := 0; i < numberOfNodes; i++ {
	// 	for j := i + 1; j < numberOfNodes; j++ {
	// 		err := Ping(nodes[i].MultiAddress, &rpc.MultiAddress{Multi: nodes[j].MultiAddress.String()})
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 	}
	// }

	// // Send messages through the topology
	// err = sendMessages(nodes)
	// Ω(err).ShouldNot(HaveOccurred())
	// })
})
