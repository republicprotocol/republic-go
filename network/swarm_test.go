package network_test

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dark-node"
)

var _ = Describe("Swarm nodes", func() {
	var err error
	var mu = new(sync.Mutex)
	var nodes []*node.DarkNode

	for _, numberOfNodes := range []int{5 /*, 16, 32, 64, 128*/} {
		for _, connectivity := range []int{100 /*, 80, 40*/} {
			func(numberOfNodes, connectivity int) {
				Context("when bootstrapping", func() {
					BeforeEach(func() {
						mu.Lock()

						nodes, err = generateNodes(numberOfNodes)
						Ω(err).ShouldNot(HaveOccurred())

						startNodes(nodes)

						err = connectNodes(nodes, connectivity)
						Ω(err).ShouldNot(HaveOccurred())

						watchNodes(nodes)
					})

					AfterEach(func() {
						stopNodes(nodes)
						mu.Unlock()
					})

					It("should reach consensus on an order match", func() {
						sendOrders(nodes)
					})
				})
			}(numberOfNodes, connectivity)
		}
	}
})
