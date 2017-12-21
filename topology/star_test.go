package topology_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"time"
)

var _ = Describe("Star topologies", func() {

	BeforeEach(func() {
		μ.Lock()
	})

	AfterEach(func() {
		μ.Unlock()
	})

	It("should route messages in a two-sided connection", func() {

		// Initialize all peers.
		peers, err := generatePeers()
		Ω(err).ShouldNot(HaveOccurred())

		// Connect all peers to the first peer.
		for i := 1; i < numberOfPeers; i++ {
			host, err := peers[i].Config.MultiAddress.ValueForProtocol(identity.IP4Code)
			Ω(err).ShouldNot(HaveOccurred())
			port, err := peers[i].Config.MultiAddress.ValueForProtocol(identity.TCPCode)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = peers[0].PingPeer(host + ":" + port)
			Ω(err).ShouldNot(HaveOccurred())
		}

		for _, peer := range peers {
			go peer.StartListening()
		}
		time.Sleep(startTimeDelay)

		// Send messages through the topology
		err = sendMessages(peers)
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should route messages in a one-sided connection", func() {

		// Initialize all peers.
		peers, err := generatePeers()
		Ω(err).ShouldNot(HaveOccurred())

		// Connect all peers to the first peer.
		for i := 1; i < numberOfPeers; i++ {
			peers[0].Config.Peers = append(peers[0].Config.Peers, peers[i].Config.MultiAddress)
		}

		for _, peer := range peers {
			go peer.StartListening()
		}
		time.Sleep(startTimeDelay)

		// Send messages through the topology
		err = sendMessages(peers)
		Ω(err).ShouldNot(HaveOccurred())
	})
})
