package main

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/order"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
)

var _ = Describe("Relay", func() {
	Context("storing and updating orders", func() {
		It("should store entry in local orderbook", func() {
			book := orderbook.NewOrderbook(100)
			block := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := storeEntry(&block, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Open))
		})

		It("should store multiple entries in local orderbook", func() {
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := storeEntry(&fstBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = storeEntry(&sndBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(2))
			Ω(blocks[0].Status).Should(Equal(order.Open))
			Ω(blocks[1].Status).Should(Equal(order.Open))
		})

		It("should update entries with a higher status", func() {
			book := orderbook.NewOrderbook(100)
			openBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			confirmedBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Confirmed{
					Confirmed: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := storeEntry(&openBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = storeEntry(&confirmedBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
		})

		It("should not update entries with a lesser status", func() {
			book := orderbook.NewOrderbook(100)
			confirmedBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Confirmed{
					Confirmed: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			openBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := storeEntry(&confirmedBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = storeEntry(&openBlock, [32]byte{}, &book)
			Ω(err).Should(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
		})
	})

	Context("forwarding orders", func() {
		It("should forward orders read from the connection", func() {
			// Construct channels
			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(blocks)
			defer close(errs)

			connections := int32(1)
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				blocks <- &fstBlock
				blocks <- &sndBlock
				errs <- errors.New("connection lost")
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := forwardMessages(blocks, errs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(2))
			wg.Wait()
		})

		It("should forward orders from multiple connections", func() {
			// Construct channels
			fstBlocks, fstErrs := make(chan *rpc.SyncBlock), make(chan error)
			sndBlocks, sndErrs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(fstBlocks)
			defer close(fstErrs)
			defer close(sndBlocks)
			defer close(sndErrs)

			connections := int32(2)
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				fstBlocks <- &fstBlock
				fstErrs <- errors.New("connection lost")
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				sndBlocks <- &sndBlock
				sndErrs <- errors.New("connection lost")
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := forwardMessages(fstBlocks, fstErrs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			err = forwardMessages(sndBlocks, sndErrs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(2))
			wg.Wait()
		})

		It("should not forward orders with an invalid epoch hash", func() {
			// Construct channels
			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(blocks)
			defer close(errs)

			connections := int32(1)
			book := orderbook.NewOrderbook(100)
			block := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				blocks <- &block
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := forwardMessages(blocks, errs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(0))
			wg.Wait()
		})
	})
})
