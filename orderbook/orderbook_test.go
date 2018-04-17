package orderbook_test

import (
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
		var book orderbook.Orderbook

		BeforeEach(func() {
			book = orderbook.NewOrderbook(10)
		})

		AfterEach(func() {
			book.Close()
		})

		It("subscribe and unsubscribe", func() {

			var chans [maxConnections]chan orderbook.Entry

			for i := 0; i < maxConnections; i++ {
				// stream := NewMockStream()
				chans[i] = make(chan orderbook.Entry)
				defer close(chans[i])
				err := book.Subscribe(chans[i])
				Ω(err).ShouldNot(HaveOccurred())
			}

			for i := 0; i < maxConnections; i++ {
				book.Unsubscribe(chans[i])
			}

		})
	})
})
