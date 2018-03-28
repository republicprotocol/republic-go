package orderbook_test

import (
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

const NumberOfMessages = 100

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

var _ = Describe("sync message queue", func() {
	Context("receiving messages from the queue", func() {
		var mockStream MockStream
		var queue orderbook.SyncMessageQueue

		BeforeEach(func() {
			mockStream = NewMockStream()
			queue = orderbook.NewSyncMessageQueue(mockStream)

			go func() {
				err := queue.Run()
				Ω(err).ShouldNot(HaveOccurred())
			}()

			time.Sleep(1 *time.Second)
		})

		AfterEach(func() {
			Ω(queue.Shutdown()).ShouldNot(HaveOccurred())
			time.Sleep(1 *time.Second)
		})

		It("should be able to receive messages from the queue", func() {
			By("You should see the received order message in console")
			for i:=0 ;i <NumberOfMessages;i ++{
				ord := newOrder([]byte(fmt.Sprintf("%d",i)))
				block := newSyncBlock(ord)
				err := queue.Send(block)
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should not be able to receive wrong type message", func() {
			By("You should not see message in console")
			for i:=0 ;i <NumberOfMessages;i ++{
				err := queue.Send(fmt.Sprintf("message %d", i))
				Ω(err).Should(HaveOccurred())
			}
		})
	})

	Context("negative tests", func() {
		It("should panic when you try to receive message from the queue", func() {
			mockStream := NewMockStream()
			queue := orderbook.NewSyncMessageQueue(mockStream)
			Ω(func() {queue.Recv()}).Should(Panic())
		})
	})
})

func newSyncBlock(ord *order.Order) *rpc.SyncBlock{
	return &rpc.SyncBlock{
		Signature: []byte{},
		Timestamp: time.Now().Unix(),
		OrderBlock:	&rpc.SyncBlock_Open	{
			Open: rpc.SerializeOrder(ord),
		},
	}
}