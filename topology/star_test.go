package topology

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-x"
)

var _ = Describe("Star topology", func() {

	Context("when pinging", func() {
		It("should update their DHTs", func() {
			μ.Lock()
			defer μ.Unlock()

			// Initialize all nodes.
			nodes, err := generateNodes(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())

			// Start serving from all nodes.
			for _, n := range nodes {
				go func(node *x.Node) {
					defer GinkgoRecover()
					Ω(node.Serve()).ShouldNot(HaveOccurred())
				}(n)
				defer func(node *x.Node) {
					node.Stop()
				}(n)
			}
			time.Sleep(startTimeDelay)

			// All nodes ping the center node.
			var wg sync.WaitGroup
			wg.Add(numberOfNodes - 1)
			for i := 1; i < numberOfNodes; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer wg.Done()
					_, err = nodes[i].RPCPing(nodes[0].MultiAddress)
					Ω(err).ShouldNot(HaveOccurred())

				}(i)
			}
			wg.Wait()

			// Center node pings all the other nodes
			wg.Add(numberOfNodes - 1)
			for i := 1; i < numberOfNodes; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer wg.Done()
					_, err = nodes[0].RPCPing(nodes[i].MultiAddress)
					Ω(err).ShouldNot(HaveOccurred())

				}(i)
			}
			wg.Wait()
		})
	})
})
