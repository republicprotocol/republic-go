package topology

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-x"
	"github.com/republicprotocol/go-x/rpc"
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
					err = x.Ping(nodes[0].MultiAddress, &rpc.MultiAddress{Multi: nodes[i].MultiAddress.String()})
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
					err = x.Ping(nodes[i].MultiAddress, &rpc.MultiAddress{Multi: nodes[0].MultiAddress.String()})
					Ω(err).ShouldNot(HaveOccurred())

				}(i)
			}
			wg.Wait()
		})
	})
})
