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
	Number_Of_Nodes = 50
	Localhost       = "127.0.0.1"
)

var _ = Describe("Node", func() {

	Describe("Star topology", func() {
		var (
			nodes  [Number_Of_Nodes]*Node
			center identity.Address
		)

		BeforeEach(func() {
			// Initialize the center node of the network
			secp, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			center = secp.PublicAddress()
			nodes[0] = swarm.NewNode(Localhost, "7000", secp.PublicAddress())
			go nodes[0].StartListening()

			// Initialize other nodes
			for i := 1; i < Number_Of_Nodes; i++ {
				keyPair, err := identity.NewKeyPair()
				Ω(err).ShouldNot(HaveOccurred())
				nodes[i] = swarm.NewNode(Localhost, strconv.Itoa(7000+i), keyPair.PublicAddress())
				go nodes[i].StartListening()

				// Ping the center node to establish connection
				pong, err := nodes[i].PingNode(Localhost + ":" + "7000")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(pong.Address).Should(Equal(string(center)))
			}
		})

		Context("Test connection of two random nodes", func() {
			It("should be able to find each other in the network", func() {

				target := nodes[18].DHT.Address
				multi, err := nodes[16].FindNode(string(target))
				Ω(err).ShouldNot(HaveOccurred())
				mAddress, err := nodes[18].MultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(multi.Multis[0]).Should(Equal(mAddress.String()))
				//// Generate random number between [1,24]
				//rand1, rand2 := rand.Intn(Number_Of_Nodes-1)+1, rand.Intn(Number_Of_Nodes-1)+1
				//for rand1 == rand2 {
				//	rand2 = rand.Intn(Number_Of_Nodes-1) + 1
				//}
				//log.Println(rand1, rand2)
				//
				//// Node rand1 trying to find node rand2
				//target := nodes[rand2].DHT.Address
				//multi, err := nodes[rand1].FindNode(string(target))
				//Ω(err).ShouldNot(HaveOccurred())
				//
				//// Should get rand2's multiaddress
				//mAddress, err := nodes[rand2].MultiAddress()
				//Ω(err).ShouldNot(HaveOccurred())
				//Ω(multi.Multis[0]).Should(Equal(mAddress.String()))
			})
		})
	})
})
