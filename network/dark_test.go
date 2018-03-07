package network_test

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
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/shamir"
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

func (delegate *mockDelegate) OnSync(from identity.MultiAddress) chan do.Option {
	syncBlock := make(chan do.Option, 1)
	syncBlock <- do.Ok(&rpc.SyncBlock{})
	return syncBlock
}

func (delegate *mockDelegate) SubscribeToLogs(channel chan do.Option) {
	go func() {
		channel <- do.Ok(&rpc.LogEvent{Type: []byte("type"), Message: []byte("message")})
	}()
}

func (delegate *mockDelegate) UnsubscribeFromLogs(chan do.Option) {
}

func (delegate *mockDelegate) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfReceivedOrderFragment++
}

func (delegate *mockDelegate) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
}

var _ = Describe("dark network", func() {
	var mu = new(sync.Mutex)

	startListening := func(nodes []*DarkOcean, listeners []net.Listener) {
		Ω(len(nodes)).Should(Equal(len(listeners)))
		for i, node := range nodes {
			go func(node *DarkOcean, listener net.Listener) {
				defer GinkgoRecover()
				node.Register()
				Ω(node.Server.Serve(listener)).ShouldNot(HaveOccurred())
			}(node, listeners[i])
		}
	}

	stopListening := func(nodes []*DarkOcean, listeners []net.Listener) {
		for _, node := range nodes {
			node.Server.Stop()
		}
		for _, listener := range listeners {
			listener.Close()
		}
	}

	sendOrderFragments := func(nodes []*DarkOcean, numberOfFragments int) {
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

	for _, numberOfNodes := range []int{4, 8, 16, 32} {
		for _, numberOfFragments := range []int{4, 8, 16, 32} {
			func(numberOfNodes, numberOfFragments int) {
				Context("when sending order fragment", func() {

					var nodes []*DarkOcean
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
					})
				})
			}(numberOfNodes, numberOfFragments)
		}
	}

	Context("when using a malformed configuration", func() {
		var nodes []*DarkOcean
		var listeners []net.Listener
		var delegate *mockDelegate
		var err error
		var numberOfNodes = 2
		var numberOfFragments = 2

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

		It("should print certain logs when debug option is greater or equal than DebugHigh", func() {
			nodes[0].Options.Debug = DebugHigh
			nodes[1].Options.Debug = DebugHigh

			listeners, err = createListener(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())
			startListening(nodes, listeners)
			time.Sleep(time.Second)

			sendOrderFragments(nodes, numberOfFragments)
		})

		It("should return error when we use wrong fragment", func() {
			from, to := nodes[0], nodes[1]
			fromMulti := "/ip4/" + DefaultIPAddress + "/tcp/3000/republic/" + from.Address().String()

			orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
			orderFragment.To = &rpc.Address{Address: to.Address().String()}
			orderFragment.From = &rpc.MultiAddress{Multi: fromMulti}
			orderFragment.MaxVolumeShare = []byte("")

			_, err = from.SendOrderFragment(context.Background(), orderFragment)
			Ω(err).Should(HaveOccurred())
		})

		It("should return error when we use wrong from address", func() {

			from, to := nodes[0], nodes[1]
			orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
			orderFragment.To = &rpc.Address{Address: to.Address().String()}
			_, err = from.SendOrderFragment(context.Background(), orderFragment)
		})

		It("should return error when we have error in context", func() {

			from, to := nodes[0], nodes[1]
			orderFragment := rpc.SerializeOrderFragment(randomOrderFragment())
			orderFragment.To = &rpc.Address{Address: to.Address().String()}

			canceledContext, cancel := context.WithCancel(context.Background())
			cancel()
			_, err = from.SendOrderFragment(canceledContext, orderFragment)
			Ω(err).Should(HaveOccurred())
		})
	})

	//Context("notifications of computation results", func() {
	//
	//	var (
	//		nodes                    []*DarkOcean
	//		listeners                []net.Listener
	//		server, client           *DarkOcean
	//		serverMulti, clientMulti identity.MultiAddress
	//		delegate                 *mockDelegate
	//		err                      error
	//		number_of_results        = 100
	//		numberOfNodes            = 2
	//	)
	//
	//	BeforeEach(func() {
	//		mu.Lock()
	//
	//		delegate = newMockDelegate()
	//		nodes, err = createNodes(delegate, numberOfNodes)
	//		Ω(err).ShouldNot(HaveOccurred())
	//		listeners, err = createListener(numberOfNodes)
	//		Ω(err).ShouldNot(HaveOccurred())
	//		server, client = nodes[0], nodes[1]
	//		server.Options.Debug = DebugHigh
	//		serverMulti, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort, server.Address()))
	//		Ω(err).ShouldNot(HaveOccurred())
	//		clientMulti, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort+1, client.Address()))
	//		Ω(err).ShouldNot(HaveOccurred())
	//
	//		startListening(nodes, listeners)
	//
	//		for i := 0; i < number_of_results; i++ {
	//			server.Notify(client.Address(), newResult(i))
	//		}
	//	})
	//
	//	AfterEach(func() {
	//		stopListening(nodes, listeners)
	//
	//		mu.Unlock()
	//	})
	//
	//	It("should be able to get all results", func() {
	//		results, err := rpc.GetFinalsFromTarget(serverMulti, clientMulti, DefaultTimeOut)
	//		Ω(err).ShouldNot(HaveOccurred())
	//		Ω(len(results)).Should(Equal(number_of_results))
	//	})
	//
	//	It("should be able to get new results", func() {
	//		resultsChan, quit := rpc.NotificationsFromTarget(serverMulti, clientMulti, DefaultTimeOut)
	//		results := make([]*compute.Final, 0)
	//
	//		var wg sync.WaitGroup
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			defer GinkgoRecover()
	//
	//			q := false
	//			for !q {
	//				select {
	//				case result := <-resultsChan:
	//					results = append(results, result.Ok.(*compute.Final))
	//				case <-quit:
	//					q = true
	//				default:
	//					continue
	//				}
	//			}
	//		}()
	//
	//		time.Sleep(time.Second * 2)
	//		quit <- struct{}{}
	//		wg.Wait()
	//		Ω(len(results)).Should(Equal(number_of_results))
	//	})
	//})

	Context("rpc call handlers", func() {
		var nodes []*DarkOcean
		var listeners []net.Listener
		var delegate *mockDelegate
		var err error
		var numberOfNodes = 2
		var from, to identity.MultiAddress

		BeforeEach(func() {
			mu.Lock()
			delegate = newMockDelegate()
			nodes, err = createNodes(delegate, numberOfNodes)
			nodes[0].Options.Debug = DebugHigh
			nodes[1].Options.Debug = DebugHigh
			Ω(err).ShouldNot(HaveOccurred())
			listeners, err = createListener(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())
			startListening(nodes, listeners)
			from, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort, nodes[0].Address()))
			Ω(err).ShouldNot(HaveOccurred())
			to, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", DefaultNodePort+1, nodes[1].Address()))
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			stopListening(nodes, listeners)
			mu.Unlock()
		})

		It("should be able to send logs to a log request", func() {
			logs, quit := rpc.Logs(to, DefaultTimeOut)
			_ = <-logs
			quit <- struct{}{}
			for len(logs) > 0 {
				<-logs
			}
		})

		//	It("should be able to respond to sync query", func() {
		//		option, quit := rpc.SyncWithTarget(to, from, DefaultTimeOut)
		//		_ = <-option
		//		quit <- struct{}{}
		//	})
		//
		//	It("should be able to respond to send order fragment commitment query", func() {
		//		err := rpc.SendOrderFragmentCommitmentToTarget(to, from, DefaultTimeOut)
		//		Ω(err).ShouldNot(HaveOccurred())
		//	})
		//
		//	It("should be able to respond to compute shard query", func() {
		//		err := rpc.AskToComputeShard(to, from, compute.Shard{}, DefaultTimeOut)
		//		Ω(err).ShouldNot(HaveOccurred())
		//	})
		//
		//	It("should be able to respond to elect shard query", func() {
		//		shard, err := rpc.StartElectShard(to, from, compute.Shard{}, DefaultTimeOut)
		//		Ω(err).ShouldNot(HaveOccurred())
		//		Ω(*shard).Should(Equal(rpc.Shard{}))
		//	})
		//
		//	It("should be able to respond to finalize shard query", func() {
		//		err := rpc.FinalizeShard(to, from, compute.DeltaShard{}, DefaultTimeOut)
		//		Ω(err).ShouldNot(HaveOccurred())
		//	})
	})
})

func randomOrderFragment() *compute.OrderFragment {
	sssShare := shamir.Share{Key: 1, Value: &big.Int{}}
	orderFragment := compute.NewOrderFragment(
		[]byte("orderID"),
		compute.OrderTypeIBBO,
		compute.OrderParityBuy,
		sssShare, sssShare, sssShare, sssShare, sssShare)
	return orderFragment
}

func createNodes(delegate DarkOceanDelegate, numberOfNodes int) ([]*DarkOcean, error) {
	nodes := make([]*DarkOcean, numberOfNodes)
	for i, _ := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nodes, err
		}
		node := NewDarkOcean(
			grpc.NewServer(),
			delegate,
			Options{
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
