package orderbook_test

import (
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/orderbook"
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

var _ = Describe("sync message queue", func() {
	Context("receiving messages from the queue", func() {

		It("should be able to receive messages from the queue", func() {
			mockStream := NewMockStream()
			queue := orderbook.NewSyncMessageQueue(mockStream)
			go func() {
				defer GinkgoRecover()
				err := queue.Run()
				Ω(err).ShouldNot(HaveOccurred())
			}()
			for i:=0 ;i <100;i ++{
				queue.Send(fmt.Sprintf("message %d", i))
			}
			time.Sleep(5 * time.Second)
			Ω(queue.Shutdown()).ShouldNot(HaveOccurred())
		})
	})
})


func newSyncBlock() *rpc.SyncBlock{
	return &rpc.SyncBlock{
		Signature: []byte{}
	}
}