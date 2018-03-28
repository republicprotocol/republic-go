package orderbook_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

const MaxConnections = 3

var _ = Describe("orderBookStreamer", func() {
	Context("subscribe and unsubscribe", func() {
		It("should be able to subscribe and unsubscribe ", func() {
			streamer := orderbook.NewOrderBookStreamer(MaxConnections)
			for i := 0; i < MaxConnections; i++ {
				go func(i int ) {
					client := NewMockStream()
					err := streamer.Subscribe(fmt.Sprintf("client%d", i), client)
					Ω(err).ShouldNot(HaveOccurred())
				}(i)
			}
			time.Sleep(1* time.Second)
			Ω(streamer.CurrentConnections()).Should(Equal(MaxConnections))

			for i := 0; i < MaxConnections; i++ {
				go func(i int) {
					client := NewMockStream()
					err := streamer.Subscribe(fmt.Sprintf("client%d", i+MaxConnections), client)
					Ω(err).Should(HaveOccurred())
				}(i)
			}
			time.Sleep(1* time.Second)
			Ω(streamer.CurrentConnections()).Should(Equal(MaxConnections))

			for i := 0; i < MaxConnections; i++ {
				streamer.Unsubscribe(fmt.Sprintf("client%d", i))
				time.Sleep(1 * time.Second)
				Ω(streamer.CurrentConnections()).Should(Equal(MaxConnections - i - 1))
			}
		})
	})

	Context("handle order status change", func() {
		var streamer orderbook.OrderBookStreamer

		BeforeEach(func() {
			streamer = orderbook.NewOrderBookStreamer(MaxConnections)
			for i := 0; i < MaxConnections; i++ {
				go func(i int ) {
					client := NewMockStream()
					err := streamer.Subscribe(fmt.Sprintf("client%d", i), client)
					Ω(err).ShouldNot(HaveOccurred())
				}(i)
			}
			time.Sleep(1 * time.Second)
		})

		It("should send message of open orders", func() {
			By("You should be able to see the orders in the console.")
			for i:=0;i <10;i ++{
				ord := newOrder(order.ID(fmt.Sprintf("%d", i )))
				streamer.Open(ord)
			}
			time.Sleep(2 * time.Second)
		})

		It("should send message of match orders", func() {
			By("You should be able to see the orders in the console.")
			for i:=0;i <10;i ++{
				ord := newOrder(order.ID(fmt.Sprintf("%d", i )))
				streamer.Match(ord)
			}
			time.Sleep(2 * time.Second)
		})

		It("should send message of confirming orders", func() {
			By("You should be able to see the orders in the console.")
			for i:=0;i <10;i ++{
				ord := newOrder(order.ID(fmt.Sprintf("%d", i )))
				streamer.Confirm(ord)
			}
			time.Sleep(2 * time.Second)
		})

		It("should send message of releasing orders", func() {
			By("You should be able to see the orders in the console.")
			for i:=0;i <10;i ++{
				ord := newOrder(order.ID(fmt.Sprintf("%d", i )))
				streamer.Release(ord)
			}
			time.Sleep(2 * time.Second)
		})

		It("should send message of settling orders", func() {
			By("You should be able to see the orders in the console.")
			for i:=0;i <10;i ++{
				ord := newOrder(order.ID(fmt.Sprintf("%d", i )))
				streamer.Settle(ord)
			}
			time.Sleep(2 * time.Second)
		})
	})
})
