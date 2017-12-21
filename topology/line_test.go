package topology_test

import (
. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
"time"
)

var _ = Describe("Line topologies", func() {

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

		// Connect each peer to the peer next to it
		for i := 0; i < numberOfPeers-1; i++ {
			peers[i].Config.Peers = append(peers[i].Config.Peers, peers[i+1].Config.MultiAddress)
		}
		for i := numberOfPeers-1; i>0 ;i-- {
			peers[i].Config.Peers = append(peers[i].Config.Peers, peers[i-1].Config.MultiAddress)
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

		// Connect each peer to next peer.
		for i := 0; i < numberOfPeers-1; i++ {
			peers[i].Config.Peers = append(peers[i].Config.Peers, peers[i+1].Config.MultiAddress)
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