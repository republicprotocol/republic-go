package dark_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-sss"
	"google.golang.org/grpc"
)

const (
	DefaultNodePort  = 3000
	DefaultIPAddress = "127.0.0.1"
	DefaultTimeOut   = 5 * time.Second
)

type mockDelegate struct {
	mu                             *sync.Mutex
	numberOfReceivedOrderFragment  int
	numberOfForwardedOrderFragment int
}

func newMockDelegate() *mockDelegate {
	return &mockDelegate{
		mu: new(sync.Mutex),
		numberOfReceivedOrderFragment:  0,
		numberOfForwardedOrderFragment: 0,
	}
}

func (delegate *mockDelegate) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfReceivedOrderFragment++
}

func (delegate *mockDelegate) OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfForwardedOrderFragment++
}

func (delegate *mockDelegate) OnSync(from identity.MultiAddress) chan do.Option {
	syncBlock := make(chan do.Option, 1)
	return syncBlock
}

func (delegate *mockDelegate) OnElectShard(from identity.MultiAddress, shard compute.Shard) compute.Shard {
	return compute.Shard{}
}

func (delegate *mockDelegate) OnComputeShard(from identity.MultiAddress, shard compute.Shard) {
}

func (delegate *mockDelegate) OnFinalizeShard(from identity.MultiAddress, finalShard compute.FinalShard) {

}

var _ = Describe("dark network", func() {
	var mu = new(sync.Mutex)

	startListening := func(nodes []*dark.Node, listeners []net.Listener) {
		Ω(len(nodes)).Should(Equal(len(listeners)))
		for i, node := range nodes {
			go func(node *dark.Node, listener net.Listener) {
				defer GinkgoRecover()
				node.Register()
				Ω(node.Server.Serve(listener)).ShouldNot(HaveOccurred())
			}(node, listeners[i])
		}
	}

	stopListening := func(nodes []*dark.Node, listeners []net.Listener) {
		for _, node := range nodes {
			node.Server.Stop()
		}
		for _, listener := range listeners {
			listener.Close()
		}
	}

	sendOrderFragments := func(nodes []*dark.Node, numberOfFragments int) {
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

	// sendResultFragments := func(nodes []*dark.Node, numberOfFragments int) {
	// 	keyPair, err := identity.NewKeyPair()
	// 	Ω(err).ShouldNot(HaveOccurred())
	// 	to := keyPair.Address()

	// 	for i := 0; i < numberOfFragments; i++ {
	// 		resultFragment := randomResultFragment()

	// 		from, target := rand.Intn(len(nodes)), rand.Intn(len(nodes))
	// 		for from == target {
	// 			target = rand.Intn(len(nodes))
	// 		}
	// 		fromMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+from, nodes[from].Address().String())
	// 		fromMultiAddress, err := identity.NewMultiAddressFromString(fromMultiAddressString)
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		targetMultiAddressString := fmt.Sprintf("/ip4/%s/tcp/%d/republic/%s", DefaultIPAddress, DefaultNodePort+target, nodes[target].Address().String())
	// 		targetMultiAddress, err := identity.NewMultiAddressFromString(targetMultiAddressString)

	// 		err = rpc.SendDeltaFragmentToTarget(targetMultiAddress, to, fromMultiAddress, resultFragment, DefaultTimeOut)
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		err = rpc.SendResultFragmentToTarget(targetMultiAddress, nodes[target].Address(), fromMultiAddress, resultFragment, DefaultTimeOut)
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 	}
	// }

	for _, numberOfNodes := range []int{4, 8, 16, 32} {
		for _, numberOfFragments := range []int{4, 8, 16, 32} {
			func(numberOfNodes, numberOfFragments int) {
				Context("when sending order fragment", func() {

					var nodes []*dark.Node
					var listeners []net.Listener
					var delegate *mockDelegate
					var err error

					BeforeEach(func() {
						mu.Lock()

						delegate = newMockDelegate()
						nodes, err = createNodes(delegate, numberOfNodes)
						Ω(err).ShouldNot(HaveOccurred())
						listeners, err = createListener(numberOfNodes)
						Ω(err).ShouldNot(HaveOccurred())

						startListening(nodes, listeners)
					})

					AfterEach(func() {
						stopListening(nodes, listeners)

						mu.Unlock()
					})

					It("should either receive the order fragment or forward it to the target", func() {
						sendOrderFragments(nodes, numberOfFragments)
						Ω(delegate.numberOfReceivedOrderFragment).Should(Equal(numberOfFragments))
						Ω(delegate.numberOfForwardedOrderFragment).Should(Equal(numberOfFragments))
					})

					// It("should either receive the result fragment or forward it to the target", func() {
					// 	sendDeltaFragments(nodes, numberOfFragments)
					// 	Ω(delegate.numberOfReceivedResultFragment).Should(Equal(numberOfFragments))
					// 	Ω(delegate.numberOfForwardedResultFragment).Should(Equal(numberOfFragments))
					// })
				})
			}(numberOfNodes, numberOfFragments)
		}
	}

	Context("when using a malformed configuration", func() {
		var nodes []*dark.Node
		var listeners []net.Listener
		var delegate *mockDelegate
		var err error
		var numberOfNodes = 2
		// var numberOfFragments = 2

		BeforeEach(func() {
			mu.Lock()

			delegate = newMockDelegate()
			nodes, err = createNodes(delegate, numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			stopListening(nodes, listeners)

			mu.Unlock()
		})

		// It("should print certain logs when debug option is greater or equal than DebugHigh", func() {
		// 	nodes[0].Options.Debug = dark.DebugHigh
		// 	nodes[1].Options.Debug = dark.DebugHigh

		// 	listeners, err = createListener(numberOfNodes)
		// 	Ω(err).ShouldNot(HaveOccurred())
		// 	startListening(nodes, listeners)
		// 	time.Sleep(time.Second)

		// 	sendOrderFragments(nodes, numberOfFragments)
		// 	sendResultFragments(nodes, numberOfFragments)
		// })

		// It("should return error when we use wrong fragment", func() {
		// 	from, to := nodes[0], nodes[1]
		// 	fromMulti := "/ip4/" + DefaultIPAddress + "/tcp/3000/republic/" + from.Address().String()

		// 	orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
		// 	orderFragment.To = &rpc.Address{Address: to.Address().String()}
		// 	orderFragment.From = &rpc.MultiAddress{Multi: fromMulti}
		// 	orderFragment.MaxVolumeShare = []byte("")

		// 	resultFragment := rpc.SerializeResultFragment(randomResultFragment())
		// 	resultFragment.From = &rpc.MultiAddress{Multi: fromMulti}
		// 	resultFragment.To = &rpc.Address{Address: to.Address().String()}
		// 	resultFragment.MaxVolumeShare = []byte("")

		// 	_, err = from.SendOrderFragment(context.Background(), orderFragment)
		// 	Ω(err).Should(HaveOccurred())
		// 	_, err = from.SendResultFragment(context.Background(), resultFragment)
		// 	Ω(err).Should(HaveOccurred())
		// })

		// It("should return error when we use wrong from address", func() {

		// 	from, to := nodes[0], nodes[1]
		// 	orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
		// 	orderFragment.To = &rpc.Address{Address: to.Address().String()}
		// 	resultFragment := rpc.SerializeResultFragment(randomResultFragment())
		// 	resultFragment.To = &rpc.Address{Address: to.Address().String()}

		// 	_, err = from.SendOrderFragment(context.Background(), orderFragment)
		// 	Ω(err).Should(HaveOccurred())
		// 	_, err = from.SendResultFragment(context.Background(), resultFragment)
		// 	Ω(err).Should(HaveOccurred())
		// })

		// It("should return error when we have error in context", func() {

		// 	from, to := nodes[0], nodes[1]
		// 	orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
		// 	orderFragment.To = &rpc.Address{Address: to.Address().String()}
		// 	resultFragment := rpc.SerializeResultFragment(randomResultFragment())
		// 	resultFragment.To = &rpc.Address{Address: to.Address().String()}

		// 	canceledContext, cancel := context.WithCancel(context.Background())
		// 	cancel()
		// 	_, err = from.SendOrderFragment(canceledContext, orderFragment)
		// 	Ω(err).Should(HaveOccurred())
		// 	_, err = from.SendResultFragment(canceledContext, resultFragment)
		// 	Ω(err).Should(HaveOccurred())
		// })
	})

	Context("notifications of computation results", func() {

		var (
			nodes                    []*dark.Node
			listeners                []net.Listener
			server, client           *dark.Node
			serverMulti, clientMulti identity.MultiAddress
			delegate                 *mockDelegate
			err                      error
			number_of_results        = 100
			numberOfNodes            = 2
		)

		BeforeEach(func() {
			mu.Lock()

			delegate = newMockDelegate()
			nodes, err = createNodes(delegate, numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())
			listeners, err = createListener(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())
			server, client = nodes[0], nodes[1]
			server.Options.Debug = dark.DebugHigh
			serverMulti, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort, server.Address()))
			Ω(err).ShouldNot(HaveOccurred())
			clientMulti, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort+1, client.Address()))
			Ω(err).ShouldNot(HaveOccurred())

			startListening(nodes, listeners)

			for i := 0; i < number_of_results; i++ {
				server.Notify(client.Address(), newResult(i))
			}
		})

		AfterEach(func() {
			stopListening(nodes, listeners)

			mu.Unlock()
		})

		It("should be able to get all results", func() {
			results, err := rpc.GetFinalsFromTarget(serverMulti, clientMulti, DefaultTimeOut)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(results)).Should(Equal(number_of_results))
		})

		It("should be able to get new results", func() {
			resultsChan, quit := rpc.NotificationsFromTarget(serverMulti, clientMulti, DefaultTimeOut)
			results := make([]*compute.Final, 0)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()

				q := false
				for !q {
					select {
					case result := <-resultsChan:
						results = append(results, result.Ok.(*compute.Final))
					case <-quit:
						q = true
					default:
						continue
					}
				}
			}()

			time.Sleep(time.Second * 2)
			quit <- struct{}{}
			wg.Wait()
			Ω(len(results)).Should(Equal(number_of_results))
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

// func randomResultFragment() *compute.DeltaFragment {
// 	sssShare := sss.Share{Key: 1, Value: &big.Int{}}
// 	resultFragment := compute.NewDeltaFragment(
// 		[]byte("butOrderID"),
// 		[]byte("sellOrderID"),
// 		[]byte("butOrderFragmentID"),
// 		[]byte("sellOrderFragmentID"),
// 		sssShare, sssShare, sssShare, sssShare, sssShare)
// 	return resultFragment
// }

func createNodes(delegate dark.Delegate, numberOfNodes int) ([]*dark.Node, error) {
	nodes := make([]*dark.Node, numberOfNodes)
	for i, _ := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nodes, err
		}
		node := dark.NewNode(
			grpc.NewServer(),
			delegate,
			dark.Options{
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

func createListener(numberOfListener int) ([]net.Listener, error) {
	listeners := make([]net.Listener, numberOfListener)
	for i := 0; i < numberOfListener; i++ {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", DefaultNodePort+i))
		if err != nil {
			return nil, err
		}
		listeners[i] = listener
	}

	return listeners, nil
}
