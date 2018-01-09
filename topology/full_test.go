package topology_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-swarm/rpc"
)

var _ = Describe("Star topologies", func() {

	It("should route messages in a two-sided connection", func() {
		μ.Lock()
		defer μ.Unlock()

		// Initialize all nodes.
		nodes, err := generateNodes()
		Ω(err).ShouldNot(HaveOccurred())

		for _, n := range nodes {
			go func(node *Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(n)
			defer func(node *Node) {
				node.Stop()
			}(n)
		}
		time.Sleep(startTimeDelay)

		// Connect all nodes to each other.
		for i := 0; i < numberOfNodes; i++ {
			fmt.Println("sender", i, "...")
			client, conn, err := NewNodeClient(nodes[i].MultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
			for j := 0; j < numberOfNodes; j++ {
				if i == j {
					continue
				}
				fmt.Println("  ping", j)
				_, err = client.Ping(context.Background(), &rpc.MultiAddress{Multi: nodes[j].MultiAddress.String()})
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(conn.Close()).ShouldNot(HaveOccurred())
		}

		// // Send messages through the topology
		// err = sendMessages(nodes)
		// Ω(err).ShouldNot(HaveOccurred())
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
