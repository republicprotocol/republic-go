package network_test

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

var (
	n        = int64(8)
	k        = int64(6)
	prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
)

var _ = Describe("Dark service", func() {
	var err error
	var mu = new(sync.Mutex)
	var darks []*network.DarkService
	var servers []*grpc.Server
	var numberOfNodes = 2
	var keypairs []*identity.KeyPair
	var pool *rpc.ClientPool

	Context("rpc function calls", func() {
		BeforeEach(func() {
			mu.Lock()

			darks, servers, keypairs, err = generateDarkServices(numberOfNodes)
			Ω(err).ShouldNot(HaveOccurred())

			err = startDarkServices(servers, darks)
			Ω(err).ShouldNot(HaveOccurred())

			// keypair = darks[0]
			multiAddressSignature, err := keypairs[0].Sign(darks[0].MultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
			pool = rpc.NewClientPool(darks[0].MultiAddress, multiAddressSignature)

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

		var fragment *order.Fragment
		It("should be able to handle OpenOrder rpc", func() {
			// Sign order fragment
			var err error
			fragment, err = generateOrderFragment()
			Ω(err).ShouldNot(HaveOccurred())
			err = fragment.Sign(*keypairs[0])
			Ω(err).ShouldNot(HaveOccurred())
			err = pool.OpenOrder(darks[1].MultiAddress, rpc.SerializeOrderFragment(fragment))
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle CancelOrder rpc", func() {
			err = pool.CancelOrder(darks[1].MultiAddress, fragment.OrderID, nil)
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
			nd.Logger.Error(err.Error())
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

func generateOrderFragment() (*order.Fragment, error) {
	price := stackint.FromUint(10)
	maxVolume := stackint.FromUint(1000)
	minVolume := stackint.FromUint(100)
	nonce := stackint.Zero()
	fragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
	if err != nil {
		return nil, err
	}
	return fragments[0], nil
}
