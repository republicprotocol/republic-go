package network_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	sss "github.com/republicprotocol/go-sss"
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

func (delegate *sendFragmentDelegate) OnOrderFragmentReceived(orderFragment *compute.OrderFragment) {
	atomic.AddInt32(&delegate.numberOfFragments, 1)
}

func (delegate *sendFragmentDelegate) OnResultFragmentReceived(orderFragment *compute.ResultFragment) {
	atomic.AddInt32(&delegate.numberOfFragments, 1)
}

var _ = Describe("Send order fragment", func() {

	send := func(nodes []*network.Node, numberOfFragments int) {
		var wg sync.WaitGroup
		wg.Add(numberOfFragments)

		for i := 0; i < numberOfFragments; i++ {
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				from, to := randomNodes(nodes)
				address, err := to.MultiAddress.Address()
				Ω(err).ShouldNot(HaveOccurred())
				orderFragment := generateOrderFragment(address.String(), from.MultiAddress.String())

				_, err = from.SendOrderFragment(context.Background(), orderFragment)
				Ω(err).ShouldNot(HaveOccurred())
			}()
		}
		wg.Wait()
	}

	run := func(name string, numberOfNodes, numberOfFragment int) int {
		var nodes []*network.Node
		var topology map[identity.Address][]*network.Node
		var err error

		delegate := newSendFragmentDelegate()
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
		send(nodes, numberOfFragment)

		return int(delegate.numberOfFragments)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
			It("should send the order fragment to the right target", func() {
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
			It("should send the order fragment to the right target", func() {
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
			It("should send the order fragment to the right target", func() {
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
			It("should send the order fragment to the right target", func() {
				numberOfMessages := numberOfNodes
				testMu.Lock()
				defer testMu.Unlock()
				numberOfFragments := run("ring", numberOfNodes, numberOfMessages)
				Ω(numberOfFragments).Should(Equal(numberOfMessages))
			})
		})
	}

})

func generateOrderFragment(to, from string) *rpc.OrderFragment {
	return &rpc.OrderFragment{
		To:           &rpc.Address{Address: to},
		From:         &rpc.Address{Address: from},
		Id:           []byte(to),
		OrderType:    int64(0),
		OrderBuySell: int64(0),
		FstCodeShare: sss.ToBytes(sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		}),
		SndCodeShare: sss.ToBytes(sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		}),
		PriceShare: sss.ToBytes(sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		}),
		MaxVolumeShare: sss.ToBytes(sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		}),
		MinVolumeShare: sss.ToBytes(sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		}),
	}
}
