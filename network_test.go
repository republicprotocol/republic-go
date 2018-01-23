package network_test

import (
	"fmt"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order-compute"
)

type pingDelegate struct {
	numberOfPings int32
}

func newPingDelegate() *pingDelegate {
	return &pingDelegate{
		numberOfPings: 0,
	}
}

func (delegate *pingDelegate) OnPingReceived(_ identity.MultiAddress) {
	atomic.AddInt32(&delegate.numberOfPings, 1)
}

func (delegate *pingDelegate) OnOrderFragmentReceived(_ identity.MultiAddress, _ *compute.OrderFragment) {
}

func (delegate *pingDelegate) OnResultFragmentReceived(_ identity.MultiAddress, _ *compute.ResultFragment) {
}

var _ = Describe("Ping RPC", func() {

	run := func(name string, numberOfNodes int) int {
		var nodes []*network.Node
		var topology map[identity.Address][]*network.Node
		var err error

		delegate := newPingDelegate()
		switch name {
		case "full":
			nodes, topology, err = generateFullyConnectedTopology(numberOfNodes, delegate)
		case "star":
			nodes, topology, err = generateStarTopology(numberOfNodes, delegate)
		case "line":
			nodes, topology, err = generateLineTopology(numberOfNodes, delegate)
		case "ring":
			nodes, topology, err = generateRingTopology(numberOfNodes, delegate)
		}
		Ω(err).ShouldNot(HaveOccurred())

		for _, node := range nodes {
			go func(node *network.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
			defer func(node *network.Node) {
				defer GinkgoRecover()
				node.Stop()
			}(node)
		}
		time.Sleep(time.Second)
		err = ping(nodes, topology)
		Ω(err).ShouldNot(HaveOccurred())

		return int(delegate.numberOfPings)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
			It("should update the DHT", func() {
				testMu.Lock()
				defer testMu.Unlock()
				numberOfPings := run("full", numberOfNodes)
				Ω(numberOfPings).Should(Equal(numberOfNodes * (numberOfNodes - 1)))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
			It("should update the DHT", func() {
				testMu.Lock()
				defer testMu.Unlock()
				numberOfPings := run("star", numberOfNodes)
				Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
			It("should update the DHT", func() {
				testMu.Lock()
				defer testMu.Unlock()
				numberOfPings := run("line", numberOfNodes)
				Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
			It("should update the DHT", func() {
				testMu.Lock()
				defer testMu.Unlock()
				numberOfPings := run("ring", numberOfNodes)
				Ω(numberOfPings).Should(Equal(2 * numberOfNodes))
			})
		})
	}

})
