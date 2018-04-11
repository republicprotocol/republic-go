package orderbook_test

import (
	"fmt"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"
)

type MockStream struct {
	grpc.ServerStream
}

func NewMockStream() MockStream {
	return MockStream{}
}

func (mockStream MockStream) Send(block *rpc.SyncBlock) error {
	log.Println(block)
	return nil
}

const maxConnections = 3

var _ = Describe("order book", func() {
	Context("creating new orderbook", func() {
		var orderBook *orderbook.OrderBook

		BeforeEach(func() {
			orderBook = orderbook.NewOrderBook(maxConnections)
		})

		It("subscribe and unsubscribe", func() {
			for i := 0; i < maxConnections; i++ {
				stream := NewMockStream()
				queue := rpc.NewSyncerServerStreamQueue(stream, maxConnections)
				err := orderBook.Subscribe(fmt.Sprintf("%d", i), queue)
				Ω(err).ShouldNot(HaveOccurred())
			}

			for i := 0; i < maxConnections; i++ {
				stream := NewMockStream()
				queue := rpc.NewSyncerServerStreamQueue(stream, maxConnections)
				err := orderBook.Subscribe(fmt.Sprintf("%d", i), queue)
				Ω(err).Should(HaveOccurred())
			}

			for i := 0; i < maxConnections; i++ {
				orderBook.Unsubscribe(fmt.Sprintf("%d", i))
			}

		})
	})
})
