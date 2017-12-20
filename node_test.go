package swarm_test

import (
	. "github.com/republicprotocol/go-swarm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm"
	"strconv"
)

const (
	Number_Of_Nodes = 25
	Localhost       = "127.0.0.1"
)

var _ = Describe("Node", func() {

	Describe("Star topology", func() {
		var (
			nodes      [Number_Of_Nodes]*Node
			center     identity.Address
			testingInt int
		)

		BeforeEach(func() {
			// Initialize the center node of the network
			secp, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			center = secp.PublicAddress()

			nodes[0] = swarm.NewNode(Localhost, "7000", secp.PublicAddress())
			go nodes[0].StartListen()

			// Initialize other nodes
			for i := 1; i < Number_Of_Nodes; i++ {
				keyPair, err := identity.NewKeyPair()
				Ω(err).ShouldNot(HaveOccurred())
				nodes[i] = swarm.NewNode(Localhost, strconv.Itoa(7000+i), keyPair.PublicAddress())
				go nodes[i].StartListen()

				// Ping the center node to establish connection
				pong, err := nodes[i].PingNode(Localhost + ":" + "7000")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(pong.Address).Should(Equal(string(center)))
			}
			testingInt = 1
		})

		Context("Test connection of two random nodes", func() {
			It("should be able to find each other in the network", func() {
				target := nodes[3].DHT.Address
				multi, err := nodes[24].FindNode(string(target))
				Ω(err).ShouldNot(HaveOccurred())
				mAddress, err := nodes[3].MultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(multi.Multis[0]).Should(Equal(mAddress.String()))
			})
		})
	})
})
