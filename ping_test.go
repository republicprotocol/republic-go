package x_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
)

type pingDelegate struct {
	numberOfPings int32
}

func newPingDelegate() *pingDelegate {
	return &pingDelegate{
		numberOfPings: 0,
	}
}

func (delegate *pingDelegate) OnPingReceived(peer identity.MultiAddress) {
	atomic.AddInt32(&delegate.numberOfPings, 1)
}

func (delegate *pingDelegate) OnOrderFragmentReceived() {
}

var _ = Describe("Ping RPC", func() {

	ping := func(nodes []*x.Node, topology map[identity.Address][]*x.Node) {
		var wg sync.WaitGroup
		wg.Add(len(nodes))

		for _, node := range nodes {
			go func(node *x.Node) {
				defer GinkgoRecover()
				defer wg.Done()

				peers := topology[node.DHT.Address]
				for _, peer := range peers {
					_, err := node.RPCPing(peer.MultiAddress)
					Ω(err).ShouldNot(HaveOccurred())
				}
			}(node)
		}

		wg.Wait()
	}

	run := func(numberOfNodes int) int {
		delegate := newPingDelegate()
		nodes, topology, err := generateFullyConnectedTopology(numberOfNodes, delegate)
		Ω(err).ShouldNot(HaveOccurred())
		for _, node := range nodes {
			go func(node *x.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
			defer func(node *x.Node) {
				defer GinkgoRecover()
				node.Stop()
			}(node)
		}
		time.Sleep(time.Second)
		ping(nodes, topology)
		return int(delegate.numberOfPings)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
			It("should update the DHT", func() {
				testMu.Lock()
				defer testMu.Unlock()
				numberOfPings := run(numberOfNodes)
				Ω(numberOfPings).Should(Equal(numberOfNodes * (numberOfNodes - 1)))
			})
		})
	}
})
