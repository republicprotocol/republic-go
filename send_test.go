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
	"github.com/republicprotocol/go-x/rpc"
)

type sendFragmentDelegate struct {
	numberOfFragments int32
}

func newSendFragmentDelegate() *sendFragmentDelegate {
	return &sendFragmentDelegate{
		numberOfFragments: 0,
	}
}

func (delegate *sendFragmentDelegate) OnPingReceived(peer identity.MultiAddress) {
}

func (delegate *sendFragmentDelegate) OnOrderFragmentReceived() {
	atomic.AddInt32(&delegate.numberOfFragments, 1)
}

var _ = Describe("Send order fragment", func() {

	send := func(nodes []*x.Node, numberOfFragments int) {
		var wg sync.WaitGroup
		wg.Add(numberOfFragments)

		for i := 0; i < numberOfFragments; i++ {
			defer wg.Done()
			from, to := randomNodes(nodes)
			address, err := to.MultiAddress.Address()
			Ω(err).ShouldNot(HaveOccurred())
			orderFragment := &rpc.OrderFragment{
				To:              string(address),
				From:            string(from.MultiAddress.String()),
				OrderID:         []byte("orderID"),
				OrderFragmentID: []byte("fragmentID"),
				OrderFragment:   []byte(address),
			}

			from.RPCSendOrderFragment(to.MultiAddress, orderFragment)
			//Ω(err).ShouldNot(HaveOccurred())
		}
		wg.Wait()
	}

	run := func(name string, numberOfNodes, numberOfFragment int) int {
		var nodes []*x.Node
		var topology map[identity.Address][]*x.Node
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

		err = ping(nodes, topology)
		Ω(err).ShouldNot(HaveOccurred())
		send(nodes, numberOfFragment)

		return int(delegate.numberOfPings)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
			It("should receive the order fragment", func() {
				numberOfMessages := numberOfNodes
				testMu.Lock()
				defer testMu.Unlock()
				numberOfFragments := run("full", numberOfNodes, numberOfMessages)
				Ω(numberOfFragments).Should(Equal(numberOfMessages))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
			It("should receive the order fragment", func() {
				numberOfMessages := numberOfNodes
				testMu.Lock()
				defer testMu.Unlock()
				numberOfFragments := run("star", numberOfNodes, numberOfMessages)
				Ω(numberOfFragments).Should(Equal(numberOfMessages))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
			It("should receive the order fragment", func() {
				numberOfMessages := numberOfNodes
				testMu.Lock()
				defer testMu.Unlock()
				numberOfFragments := run("line", numberOfNodes, numberOfMessages)
				Ω(numberOfFragments).Should(Equal(numberOfMessages))
			})
		})
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
			It("should receive the order fragment", func() {
				numberOfMessages := numberOfNodes
				testMu.Lock()
				defer testMu.Unlock()
				numberOfFragments := run("ring", numberOfNodes, numberOfMessages)
				Ω(numberOfFragments).Should(Equal(numberOfMessages))
			})
		})
	}

})
