package topology

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	swarm "github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-swarm/rpc"
)

var _ = Describe("Pair topologies", func() {

	Context("when pinging", func() {
		It("should have updated DHTs", func() {
			μ.Lock()
			defer μ.Unlock()

			// Create the left Node.
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000, keyPair.Address()))
			Ω(err).ShouldNot(HaveOccurred())
			left, err := swarm.NewNode(&swarm.Config{
				KeyPair:      keyPair,
				MultiAddress: multiAddress,
				Peers:        make(identity.MultiAddresses, 0, numberOfNodes-1),
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Create the right Node.
			keyPair, err = identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err = identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3001, keyPair.Address()))
			Ω(err).ShouldNot(HaveOccurred())
			right, err := swarm.NewNode(&swarm.Config{
				KeyPair:      keyPair,
				MultiAddress: multiAddress,
				Peers:        make(identity.MultiAddresses, 0, numberOfNodes-1),
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Start the left and right Node.
			go func() {
				defer GinkgoRecover()
				Ω(left.Serve()).ShouldNot(HaveOccurred())
			}()
			go func() {
				defer GinkgoRecover()
				Ω(right.Serve()).ShouldNot(HaveOccurred())
			}()
			defer left.Stop()
			defer right.Stop()
			time.Sleep(startTimeDelay)

			// Ping the left Node from the right Node.
			err = swarm.Ping(left.MultiAddress, &rpc.MultiAddress{Multi: right.MultiAddress.String()})
			Ω(err).ShouldNot(HaveOccurred())

			// Ping the right Node from the left Node.
			err = swarm.Ping(right.MultiAddress, &rpc.MultiAddress{Multi: left.MultiAddress.String()})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(len(left.DHT.MultiAddresses())).Should(Equal(1))
			Ω(len(right.DHT.MultiAddresses())).Should(Equal(1))
		})
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
