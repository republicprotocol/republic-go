package topology

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-swarm/rpc"
)

var _ = Describe("Line topology", func() {

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
			for i := 0; i < numberOfNodes ; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer wg.Done()

					if i != 0{
						err = swarm.Ping(nodes[i-1].MultiAddress,&rpc.MultiAddress{Multi:nodes[i].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}else{
						err = swarm.Ping(nodes[numberOfNodes-1].MultiAddress,&rpc.MultiAddress{Multi:nodes[0].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}

					if i != numberOfNodes - 1{
						err = swarm.Ping(nodes[i+1].MultiAddress,&rpc.MultiAddress{Multi:nodes[i].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}else{
						err = swarm.Ping(nodes[0].MultiAddress,&rpc.MultiAddress{Multi:nodes[i].MultiAddress.String()})
						Ω(err).ShouldNot(HaveOccurred())
					}
				}(i)
			}
			wg.Wait()
		})
	})
})
