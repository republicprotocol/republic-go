package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
)

var _ = Describe("Relay", func() {
	Context("synchronizing orders", func() {
		It("should store entry in local orderbook", func() {
			book := orderbook.NewOrderbook(100)
			block := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id:     &rpc.OrderId{},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			storeEntry(&block, [32]byte{}, &book)

			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
		})
	})
})
