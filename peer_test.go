package swarm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-swarm"

	"fmt"
	"math/rand"
)

const (
	Number_Of_Peers = 50
	Localhost       = "127.0.0.1"
)

var _ = Describe("Node", func() {

	Describe("Star topology", func() {
		var (
			nodes  [Number_Of_Peers]*Peer
			center identity.Address
		)

		BeforeEach(func() {
			peers := make([]*Peer, Number_Of_Peers)
			for i := 0; i < Number_Of_Peers; i++ {
				keyPair, err := identity.NewKeyPair()
				Ω(err).ShouldNot(HaveOccurred())

				multiAddress, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000+i, keyPair.PublicAddress()))
				Ω(err).ShouldNot(HaveOccurred())
				peers[i] = NewPeer(&Config{
					KeyPair:      keyPair,
					MultiAddress: multiAddress,
					Peers:        make([]identity.MultiAddress, 0, Number_Of_Peers-1),
				})
			}

			center = peers[0].Config.KeyPair.PublicAddress()
			// Connect all peers to the first peer.
			for i := 1; i < Number_Of_Peers; i++ {
				peers[0].Config.Peers = append(peers[0].Config.Peers, peers[i].Config.MultiAddress)
				peers[i].Config.Peers = append(peers[i].Config.Peers, peers[0].Config.MultiAddress)
			}
		})

		Context("Test connection of two random nodes", func() {
			It("should be able to find each other in the network", func() {

				// Generate random number between [1,24]
				rand1, rand2 := rand.Intn(Number_Of_Peers-1)+1, rand.Intn(Number_Of_Peers-1)+1
				for rand1 == rand2 {
					rand2 = rand.Intn(Number_Of_Peers-1) + 1
				}

				// Node rand1 trying to find node rand2
				target := nodes[rand2].DHT.Address
				multi, err := nodes[rand1].FindPeer(target)
				Ω(err).ShouldNot(HaveOccurred())

				// Should get rand2's multiaddress
				Ω(err).ShouldNot(HaveOccurred())
				Ω(multi.Multis[0].Multi).Should(Equal(nodes[rand2].Config.MultiAddress.String()))
			})
		})
	})
})
