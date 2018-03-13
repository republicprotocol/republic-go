package network_test

import (
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

var (
	n        = int64(8)
	k        = int64(6)
	prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
)

var _ = Describe("Dark service", func() {
	var err error
	var mu = new(sync.Mutex)
	var darks []*network.DarkService
	var servers []*grpc.Server
	var numberOfNodes = 2
	var pool *rpc.ClientPool

	Context("rpc function calls", func() {
		BeforeEach(func() {
			mu.Lock()

			darks, servers, err = generateDarkServices(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())

			err = startDarkServices(servers, darks)
			Ω(err).ShouldNot(HaveOccurred())

			pool = rpc.NewClientPool(darks[0].MultiAddress)

			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			stopDarkServices(servers, darks)
			mu.Unlock()
		})

		It("should be able to handle Sync rpc", func() {
			_, err = pool.Sync(darks[1].MultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle OpenOrder rpc", func() {
			fragment, err := generateOrderFragment()
			Ω(err).ShouldNot(HaveOccurred())
			err = pool.OpenOrder(darks[1].MultiAddress, &rpc.OrderSignature{}, fragment)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle CancelOrder rpc", func() {
			err = pool.CancelOrder(darks[1].MultiAddress, &rpc.OrderSignature{})
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle SignOrderFragment rpc", func() {
			signature, err := pool.SignOrderFragment(darks[1].MultiAddress, &rpc.OrderFragmentSignature{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*signature).Should(Equal(rpc.OrderFragmentSignature{}))
		})

		It("should be able to handle RandomFragmentShares rpc", func() {
			fragments, err := pool.RandomFragmentShares(darks[1].MultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*fragments).Should(Equal(rpc.RandomFragments{}))
		})

		It("should be able to handle ResidueFragmentShares rpc", func() {
			fragments, err := pool.ResidueFragmentShares(darks[1].MultiAddress, &rpc.RandomFragments{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*fragments).Should(Equal(rpc.ResidueFragments{}))
		})

		It("should be able to handle ComputeResidueFragment rpc", func() {
			err = pool.ComputeResidueFragment(darks[1].MultiAddress, &rpc.ResidueFragments{})
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle BroadcastAlphaBetaFragment rpc", func() {
			fragment, err := pool.BroadcastAlphaBetaFragment(darks[1].MultiAddress, &rpc.AlphaBetaFragment{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*fragment).Should(Equal(rpc.AlphaBetaFragment{}))
		})
	})
})

func startDarkServices(servers []*grpc.Server, darks []*network.DarkService) error {
	for i, nd := range darks {
		server := grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
		servers[i] = server

		nd.Register(server)
		host, err := nd.MultiAddress.ValueForProtocol(identity.IP4Code)
		if err != nil {
			return err
		}
		port, err := nd.MultiAddress.ValueForProtocol(identity.TCPCode)
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", host+":"+port)

		if err != nil {
			nd.Logger.Error(logger.TagNetwork, err.Error())
		}
		go func() {
			if err := server.Serve(listener); err != nil {
				log.Fatal("fail to start the server")
			}
		}()
	}

	return nil
}

func stopDarkServices(servers []*grpc.Server, darks []*network.DarkService) {
	for _, server := range servers {
		server.Stop()
	}
	for _, dk := range darks {
		dk.Logger.Stop()
	}
}

func generateOrderFragment() (*rpc.OrderFragment, error) {
	fragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	return rpc.SerializeOrderFragment(fragments[0]), nil
}
