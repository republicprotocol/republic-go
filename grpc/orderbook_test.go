package grpc_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Orderbook", func() {

	var serverMock *mockOrderbookServer
	var server *Server
	var service OrderbookService
	var serviceMultiAddr identity.MultiAddress
	var client orderbook.Client

	BeforeEach(func() {
		var err error

		client = NewOrderbookClient()

		serverMock = &mockOrderbookServer{}
		server = NewServer()
		service = NewOrderbookService(serverMock)
		service.Register(server)

		serviceEcdsaKey, err := crypto.RandomEcdsaKey()
		Expect(err).ShouldNot(HaveOccurred())

		serviceMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", serviceEcdsaKey.Address()))
		Expect(err).ShouldNot(HaveOccurred())

		go server.Start("0.0.0.0:18514")
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when opening order fragments", func() {

		It("should accept order fragments", func() {
			orderFragment, err := createEncryptedFragment()
			Expect(err).ShouldNot(HaveOccurred())
			err = client.OpenOrder(context.Background(), serviceMultiAddr, orderFragment)
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

})

type mockOrderbookServer struct {
	n int64
}

func (server *mockOrderbookServer) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	atomic.AddInt64(&server.n, 1)
	return nil
}

func createEncryptedFragment() (order.EncryptedFragment, error) {
	ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, order.SettlementRenEx, time.Now().Add(time.Hour), order.TokensBTCETH, order.NewCoExp(1, 1), order.NewCoExp(1, 1), order.NewCoExp(1, 1), 1)
	ordFragments, err := ord.Split(6, 4)
	if err != nil {
		return order.EncryptedFragment{}, err
	}
	rsaKey, err := crypto.RandomRsaKey()
	if err != nil {
		return order.EncryptedFragment{}, err
	}
	return ordFragments[0].Encrypt(rsaKey.PublicKey)
}
