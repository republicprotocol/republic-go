package topology_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		// Connect all peers to each other.
		for i := 0; i < numberOfPeers; i++ {
			for j := 0; j < numberOfPeers; j++ {
				if i == j {
					continue
				}
				peers[i].Peers = append(peers[i].Config.Peers, peers[j].Config.MultiAddress)
			}
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

		// Connect all peers to each other.
		for i := 0; i < numberOfPeers; i++ {
			for j := i + 1; j < numberOfPeers; j++ {
				peers[i].Config.Peers = append(peers[i].Config.Peers, peers[j].Config.MultiAddress)
			}
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
