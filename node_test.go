package xing_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
	"github.com/republicprotocol/go-xing"
)

const (
	DefaultNodePort  = 3000
	DefaultIPAddress = "127.0.0.1"
	DefaultTimeOut   = 5 * time.Second
)

type mockDelegate struct {
	mu                              *sync.Mutex
	numberOfReceivedOrderFragment   int
	numberOfForwardedOrderFragment  int
	numberOfReceivedResultFragment  int
	numberOfForwardedResultFragment int
}

func newMockDelegate() *mockDelegate {
	return &mockDelegate{
		mu: new(sync.Mutex),
		numberOfReceivedOrderFragment:   0,
		numberOfForwardedOrderFragment:  0,
		numberOfReceivedResultFragment:  0,
		numberOfForwardedResultFragment: 0,
	}
}

func (delegate *mockDelegate) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfReceivedOrderFragment++
}

func (delegate *mockDelegate) OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfReceivedResultFragment++
}

func (delegate *mockDelegate) OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfForwardedOrderFragment++
}

func (delegate *mockDelegate) OnResultFragmentForwarding(to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfForwardedResultFragment++
}

var _ = Describe("Xing overlay network", func() {
	var mu = new(sync.Mutex)

	startListening := func(nodes []*xing.Node) {
		for i, node := range nodes {
			listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", DefaultNodePort+i))
			Ω(err).ShouldNot(HaveOccurred())
			go func(node *xing.Node) {
				defer GinkgoRecover()
				Ω(node.Serve(listener)).ShouldNot(HaveOccurred())
			}(node)
		}
	}

	stopListening := func(nodes []*xing.Node) {
		for _, node := range nodes {
			node.Stop()
		}
	}

	sendOrderFragments := func(nodes []*xing.Node, numberOfFragments int) {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		to := keyPair.Address()

		for i := 0; i < numberOfFragments; i++ {
			orderFragment := randomOrderFragment()

			from, target := rand.Intn(len(nodes)), rand.Intn(len(nodes))
			for from == target {
				target = rand.Intn(len(nodes))
			}
			fromMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+from, nodes[from].Address().String())
			fromMultiAddress, err := identity.NewMultiAddressFromString(fromMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			targetMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+target, nodes[target].Address().String())
			targetMultiAddress, err := identity.NewMultiAddressFromString(targetMultiAddressString)

			Ω(err).ShouldNot(HaveOccurred())
			err = rpc.SendOrderFragmentToTarget(targetMultiAddress, to, fromMultiAddress, orderFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
			err = rpc.SendOrderFragmentToTarget(targetMultiAddress, nodes[target].Address(), fromMultiAddress, orderFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
		}
	}

	sendResultFragments := func(nodes []*xing.Node, numberOfFragments int) {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		to := keyPair.Address()

		for i := 0; i < numberOfFragments; i++ {
			resultFragment := randomResultFragment()

			from, target := rand.Intn(len(nodes)), rand.Intn(len(nodes))
			for from == target {
				target = rand.Intn(len(nodes))
			}
			fromMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+from, nodes[from].Address().String())
			fromMultiAddress, err := identity.NewMultiAddressFromString(fromMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			targetMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+target, nodes[target].Address().String())
			targetMultiAddress, err := identity.NewMultiAddressFromString(targetMultiAddressString)

			err = rpc.SendResultFragmentToTarget(targetMultiAddress, to, fromMultiAddress, resultFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
			err = rpc.SendResultFragmentToTarget(targetMultiAddress, nodes[target].Address(), fromMultiAddress, resultFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
		}
	}

	for _, numberOfNodes := range []int{4, 8, 16, 32} {
		for _, numberOfFragments := range []int{4, 8, 16, 32} {
			func(numberOfNodes, numberOfFragments int) {
				Context("when sending order fragment", func() {

					BeforeEach(func() {
						mu.Lock()
					})

					AfterEach(func() {
						mu.Unlock()
					})

					It("should either receive the order fragment or forward it to the target", func() {
						delegate := newMockDelegate()
						nodes, err := createNodes(delegate, numberOfNodes, DefaultNodePort)
						Ω(err).ShouldNot(HaveOccurred())
						startListening(nodes)
						sendOrderFragments(nodes, numberOfFragments)
						Ω(delegate.numberOfReceivedOrderFragment).Should(Equal(numberOfFragments))
						Ω(delegate.numberOfForwardedOrderFragment).Should(Equal(numberOfFragments))
						stopListening(nodes)
					})

					It("should either receive the result fragment or forward it to the target", func() {
						delegate := newMockDelegate()
						nodes, err := createNodes(delegate, numberOfNodes, DefaultNodePort)
						Ω(err).ShouldNot(HaveOccurred())
						startListening(nodes)
						sendResultFragments(nodes, numberOfFragments)
						Ω(delegate.numberOfReceivedResultFragment).Should(Equal(numberOfFragments))
						Ω(delegate.numberOfForwardedResultFragment).Should(Equal(numberOfFragments))
						stopListening(nodes)
					})
				})
			}(numberOfNodes, numberOfFragments)
		}
	}

	Context("when using a malformed configuration", func() {
		var numberOfNodes = 2
		var numberOfFragments = 2

		BeforeEach(func() {
			mu.Lock()
		})

		AfterEach(func() {
			mu.Unlock()
		})

		It("should print certain logs when debug option is greater or equal than DebugHigh", func() {
			delegate := newMockDelegate()
			nodes, err := createNodes(delegate, numberOfNodes, DefaultNodePort)
			Ω(err).ShouldNot(HaveOccurred())

			listener0, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", DefaultNodePort))
			Ω(err).ShouldNot(HaveOccurred())

			listener1, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", DefaultNodePort+1))
			Ω(err).ShouldNot(HaveOccurred())

			nodes[0].Options.Debug = xing.DebugHigh
			nodes[1].Options.Debug = xing.DebugHigh
			go func() {
				defer GinkgoRecover()
				Ω(nodes[0].Serve(listener0)).ShouldNot(HaveOccurred())
			}()
			go func() {
				defer GinkgoRecover()
				Ω(nodes[1].Serve(listener1)).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Second)

			sendOrderFragments(nodes, numberOfFragments)
			sendResultFragments(nodes, numberOfFragments)
			time.Sleep(1)

			stopListening(nodes)
		})

		It("should return error when we use wrong fragment", func() {
			delegate := newMockDelegate()
			nodes, err := createNodes(delegate, numberOfNodes, DefaultNodePort)
			Ω(err).ShouldNot(HaveOccurred())

			from, to := nodes[0], nodes[1]
			fromMulti := "/ip4/" + DefaultIPAddress + "/tcp/3000/republic/" + from.Address().String()
			orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
			orderFragment.To = &rpc.Address{Address: to.Address().String()}
			orderFragment.From = &rpc.MultiAddress{Multi: fromMulti}
			orderFragment.MaxVolumeShare = []byte("")
			resultFragment := rpc.SerializeResultFragment(randomResultFragment())
			resultFragment.From = &rpc.MultiAddress{Multi: fromMulti}
			resultFragment.To = &rpc.Address{Address: to.Address().String()}
			resultFragment.MaxVolumeShare = []byte("")

			_, err = from.SendOrderFragment(context.Background(), orderFragment)
			Ω(err).Should(HaveOccurred())
			_, err = from.SendResultFragment(context.Background(), resultFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return error when we use wrong from address", func() {
			delegate := newMockDelegate()
			nodes, err := createNodes(delegate, numberOfNodes, DefaultNodePort)
			Ω(err).ShouldNot(HaveOccurred())

			from, to := nodes[0], nodes[1]
			orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
			orderFragment.To = &rpc.Address{Address: to.Address().String()}
			resultFragment := rpc.SerializeResultFragment(randomResultFragment())
			resultFragment.To = &rpc.Address{Address: to.Address().String()}

			_, err = from.SendOrderFragment(context.Background(), orderFragment)
			Ω(err).Should(HaveOccurred())
			_, err = from.SendResultFragment(context.Background(), resultFragment)
			Ω(err).Should(HaveOccurred())
		})
	})
})

func randomOrderFragment() *compute.OrderFragment {
	sssShare := sss.Share{Key: 1, Value: &big.Int{}}
	orderFragment := compute.NewOrderFragment(
		[]byte("orderID"),
		compute.OrderTypeIBBO,
		compute.OrderParityBuy,
		sssShare, sssShare, sssShare, sssShare, sssShare)
	return orderFragment
}

func randomResultFragment() *compute.ResultFragment {
	sssShare := sss.Share{Key: 1, Value: &big.Int{}}
	resultFragment := compute.NewResultFragment(
		[]byte("butOrderID"),
		[]byte("sellOrderID"),
		[]byte("butOrderFragmentID"),
		[]byte("sellOrderFragmentID"),
		sssShare, sssShare, sssShare, sssShare, sssShare)
	return resultFragment
}

func createNodes(delegate xing.Delegate, numberOfNodes, port int) ([]*xing.Node, error) {
	nodes := make([]*xing.Node, numberOfNodes)
	for i, _ := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nodes, err
		}
		node := xing.NewNode(
			delegate,
			xing.Options{
				Address:        keyPair.Address(),
				Debug:          DefaultOptionsDebug,
				Timeout:        DefaultOptionsTimeout,
				TimeoutStep:    DefaultOptionsTimeoutStep,
				TimeoutRetries: DefaultOptionsTimeoutRetries,
				Concurrent:     DefaultOptionsConcurrent,
			},
		)
		nodes[i] = node
	}
	return nodes, nil
}
