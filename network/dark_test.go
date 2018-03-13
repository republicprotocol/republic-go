package network_test

import (
	"log"
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

//var (
//	n = int64(8)
//	k = int64(6)
//	prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
//)

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

		// todo : test for unimplemented function
		It("should be able to handle Sync rpc", func() {
			_, err = pool.Sync(darks[1].MultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should be able to handle openOrder rpc", func() {
			fragment := &rpc.OrderFragment{
				Signature: []byte("signature"),
				Id:   []byte("id"),
				OrderId: []byte("orderId"),
				OrderType: int64(order.TypeLimit),
				OrderParity: int64(order.ParityBuy),
				FstCodeShare: []byte("first"),
				SndCodeShare: []byte("second"),
				PriceShare: []byte("price"),
				MaxVolumeShare: []byte("max"),
				MinVolumeShare: []byte("min"),
			}
			err = pool.OpenOrder(darks[1].MultiAddress, &rpc.OrderSignature{}, fragment)
			Ω(err).ShouldNot(HaveOccurred())
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
