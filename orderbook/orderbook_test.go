package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/orderbook"
)

const maxConnections = 3

var _ = Describe("order book", func() {
	Context("creating new orderbook", func() {
		var book orderbook.Orderbook

		BeforeEach(func() {
			book = orderbook.NewOrderbook()
		})

		AfterEach(func() {
			book.Close()
		})

		It("listen and stop listening", func() {
			dones := [maxConnections]chan struct{}{}
			for i := 0; i < maxConnections; i++ {
				// stream := NewMockStream()
				dones[i] = make(chan struct{})
				defer close(dones[i])
				book.Listen(dones[i])
			}
		})
	})
})
