package topology

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-x"
	"github.com/republicprotocol/go-x/rpc"
)

var _ = Describe("Fully connected mesh topologies", func() {

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
						err = x.Ping(nodes[j].MultiAddress, &rpc.MultiAddress{Multi: nodes[i].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}

				}(i)
			}
			wg.Wait()
		})
	})
})
