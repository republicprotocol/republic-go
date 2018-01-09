package topology

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-swarm/rpc"
)

var _ = Describe("Fully connected mesh topologies", func() {

	const numberOfNodes = 75

	Context("when pinging", func() {
		It("should update their DHTs", func() {
			μ.Lock()
			defer μ.Unlock()

			// Initialize all nodes.
			nodes, err := generateNodes(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())

			// Start serving from all nodes.
			for _, n := range nodes {
				go func(node *swarm.Node) {
					defer GinkgoRecover()
					Ω(node.Serve()).ShouldNot(HaveOccurred())
				}(n)
				defer func(node *swarm.Node) {
					node.Stop()
				}(n)
			}
			time.Sleep(startTimeDelay)

			// Connect all nodes to each other concurrently.
			var wg sync.WaitGroup
			wg.Add(numberOfNodes)
			for i := 0; i < numberOfNodes; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer wg.Done()

					for j := 0; j < numberOfNodes; j++ {
						if i == j {
							continue
						}
						err = swarm.Ping(nodes[j].MultiAddress, &rpc.MultiAddress{Multi: nodes[i].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}

				}(i)
			}
			wg.Wait()
		})
	})

	// It("should route messages in a one-sided connection", func() {

	// // Initialize all nodes.
	// nodes, err := generateNodes()
	// Ω(err).ShouldNot(HaveOccurred())

	// for _, node := range nodes {
	// 	go func(node *Node) {
	// 		defer GinkgoRecover()
	// 		Ω(node.Serve()).ShouldNot(HaveOccurred())
	// 	}(node)
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
