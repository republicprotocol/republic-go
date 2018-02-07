package xing_test

import (
	"fmt"
	"math/big"
	"math/rand"
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
	DefaultNodePort = 3000
	DefaultTimeOut  = 5 * time.Second
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

var _ = Describe("nodes of Xing network", func() {
	var delegate *mockDelegate
	var nodes []*xing.Node

	startListening := func() {
		for _, node := range nodes {
			go node.Serve()
		}
	}

	sendOrderFragments := func(numberOfFragments int) {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		to := keyPair.Address()

		for i := 0; i < numberOfFragments; i++ {
			orderFragment := randomOrderFragment()

			from, target := rand.Intn(len(nodes)), rand.Intn(len(nodes))
			for from == target {
				target = rand.Intn(len(nodes))
			}
			fromMultiAddressString := "/ip4/" + nodes[from].Options.Host +
				"/tcp/" + nodes[from].Options.Port +
				"/republic/" + nodes[from].Address().String()
			targetMultiAddressString := "/ip4/" + nodes[target].Options.Host +
				"/tcp/" + nodes[target].Options.Port +
				"/republic/" + nodes[target].Address().String()
			fromMultiAddress, err := identity.NewMultiAddressFromString(fromMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			targetMultiAddress, err := identity.NewMultiAddressFromString(targetMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			err = rpc.SendOrderFragmentToTarget(targetMultiAddress, to, fromMultiAddress, orderFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
			// test forwar fragment
			// err = rpc.SendOrderFragmentToTarget(toMultiAddress, fromMultiAddress, orderFragment, DefaultTimeOut)
			// Ω(err).ShouldNot(HaveOccurred())
		}
	}

	AfterEach(func() {
		for _, node := range nodes {
			func(node *xing.Node) {
				node.Stop()
			}(node)
		}
	})

	sendResultFragments := func(numberOfFragments int) {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		to := keyPair.Address()

		for i := 0; i < numberOfFragments; i++ {
			resultFragment := randomResultFragment()

			from, target := rand.Intn(len(nodes)), rand.Intn(len(nodes))
			for from == target {
				target = rand.Intn(len(nodes))
			}
			fromMultiAddressString := "/ip4/" + nodes[from].Options.Host +
				"/tcp/" + nodes[from].Options.Port +
				"/republic/" + nodes[from].Address().String()
			targetMultiAddressString := "/ip4/" + nodes[target].Options.Host +
				"/tcp/" + nodes[target].Options.Port +
				"/republic/" + nodes[target].Address().String()
			fromMultiAddress, err := identity.NewMultiAddressFromString(fromMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			targetMultiAddress, err := identity.NewMultiAddressFromString(targetMultiAddressString)
			Ω(err).ShouldNot(HaveOccurred())
			err = rpc.SendResultFragmentToTarget(targetMultiAddress, to, fromMultiAddress, resultFragment, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
		}
	}

	for numberOfNodes := range []int{4, 8, 16, 32} {
		for numberOfFragments := range []int{4, 8, 16, 32} {
			Context("sending order fragment", func() {
				It("should  either receive the order fragment or forward it to the target", func() {
					delegate = newMockDelegate()
					nodes = createNodes(delegate, numberOfNodes, DefaultNodePort)
					startListening()
					sendOrderFragments(numberOfFragments)
					Ω(delegate.numberOfReceivedOrderFragment + delegate.numberOfForwardedOrderFragment).Should(Equal(numberOfFragments))
				})

				It("should  either receive the result fragment or forward it to the target", func() {
					delegate = newMockDelegate()
					nodes = createNodes(delegate, numberOfNodes, DefaultNodePort)
					startListening()
					sendResultFragments(numberOfFragments)
					Ω(delegate.numberOfReceivedResultFragment + delegate.numberOfForwardedResultFragment).Should(Equal(numberOfFragments))
				})
			})
		}
	}

	Context("negative tests for sending order/result fragments", func() {
		var numberOfNodes = 2
		var numberOfFragments = 2

		BeforeEach(func() {
			delegate = newMockDelegate()
			nodes = createNodes(delegate, numberOfNodes, DefaultNodePort)
		})

		It("can't use an occupied ip address and port", func() {
			nodes[1].Options.Port = fmt.Sprintf("%d", DefaultNodePort)
			nodes[0].Options.Port = fmt.Sprintf("%d", DefaultNodePort)
			go nodes[0].Serve()
			defer nodes[0].Stop()
			go func() {
				err := nodes[1].Serve()
				Ω(err).Should(HaveOccurred())
				By("get this ")
			}()
		})

		It("should print certain logs when debug option is greater or equal than DebugLow", func() {
			nodes[0].Options.Debug = xing.DebugLow
			startListening()
			sendOrderFragments(numberOfFragments)
			sendResultFragments(numberOfFragments)
			for _, node := range nodes {
				node.Stop()
			}
		})

		It("should print certain logs when debug option is greater or equal than DebugHigh", func() {
			nodes[0].Options.Debug = xing.DebugHigh
			startListening()
			sendOrderFragments(numberOfFragments)
			sendResultFragments(numberOfFragments)
		})

		It("should return error when we use wrong target address", func() {

		})

		It("should return error when we use wrong to address", func() {

		})

		It("should return error when we use wrong from address", func() {

		})
	})
})

func randomOrderFragment() *compute.OrderFragment {
	sssShare := sss.Share{Key: 1, Value: &big.Int{}}
	orderFragment := compute.NewOrderFragment([]byte("orderID"), compute.OrderTypeIBBO, compute.OrderParityBuy,
		sssShare, sssShare, sssShare, sssShare, sssShare)
	return orderFragment
}

func randomResultFragment() *compute.ResultFragment {
	sssShare := sss.Share{Key: 1, Value: &big.Int{}}
	resultFragment := compute.NewResultFragment([]byte("butOrderID"), []byte("sellOrderID"),
		[]byte("butOrderFragmentID"), []byte("sellOrderFragmentID"),
		sssShare, sssShare, sssShare, sssShare, sssShare)
	return resultFragment
}

func createNodes(delegate xing.Delegate, numberOfNodes, port int) []*xing.Node {
	nodes := make([]*xing.Node, numberOfNodes)
	for i, _ := range nodes {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		node := xing.NewNode(
			delegate,
			xing.Options{
				Address:        keyPair.Address(),
				Host:           "127.0.0.1",
				Port:           fmt.Sprintf("%d", port+i),
				Debug:          DefaultOptionsDebug,
				Timeout:        DefaultOptionsTimeout,
				TimeoutStep:    DefaultOptionsTimeoutStep,
				TimeoutRetries: DefaultOptionsTimeoutRetries,
				Concurrent:     DefaultOptionsConcurrent,
			},
		)
		nodes[i] = node
	}
	return nodes
}
