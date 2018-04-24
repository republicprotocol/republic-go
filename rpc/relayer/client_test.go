package relayer_test

import (
	"context"
	"log"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
)

var _ = Describe("Relayer", func() {
	/* var conn client.Connection
	var darknodes darknode.Darknodes
	var bootstrapNodes []string

	BeforeSuite(func() {
		var err error

		// Connect to Ganache
		conn, err = ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())
		darknodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), conn, ganache.GenesisTransactor(), &bind.CallOpts{})
		Expect(err).ShouldNot(HaveOccurred())
		darknodeRegistry.SetGasLimit(1000000)

		// Create DarkNodes and contexts/cancels for running them
		darknodes, err = NewDarknodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
		Expect(err).ShouldNot(HaveOccurred())

		// Populate bootstrap nodes
		bootstrapNodes = make([]string, 5)
		for i := 0; i < NumberOfBootstrapDarkNodes; i++ {
			bootstrapNodes[i] = darknodes[i].MultiAddress().String()
		}

		// Assign trader address with the first dark node multiaddress
		traderMulti = darknodes[0].MultiAddress().String()

		// Register the Darknodes and trigger an epoch to accept their
		// registrations
		err = RegisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		var err error

		// Deregister the DarkNodes
		err = DeregisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())

		// Refund the DarkNodes
		err = RefundDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())
	}) */

	Context("when merging entries", func() {
		It("should update the orderbook with a valid response", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse(make([]byte, 32), OrderStatus_Open)
			err := mergeEntry(&orderbook, &syncResponse)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(1))
		})

		It("should not update the orderbook with a nil entry", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := SyncResponse{
				Signature: []byte{},
				Epoch:     make([]byte, 32),
				Entry:     nil,
			}
			err := mergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})

		It("should not update the orderbook with a nil epoch hash", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse(nil, OrderStatus_Open)
			err := mergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})

		It("should not update the orderbook with an epoch hash with invalid length", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse([]byte{}, OrderStatus_Open)
			err := mergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})

		It("should not update the orderbook with an invalid order status", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse(make([]byte, 32), -1)
			err := mergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})
	})

	Context("when syncing orders", func() {
		It("should sync the orderbook from a specified node", func() {
			dht := dht.NewDHT("", 100)
			connPool := client.NewConnPool(100)
			client := NewClient(&dht, &connPool)

			ctx := context.Background()
			var epochHash [32]byte
			// msgs, errs := client.SyncFrom(context.Background(), darknodes[0].MultiAddress, epochHash)
			msgs, errs := client.SyncFrom(ctx, identity.MultiAddress{}, epochHash)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case msg := <-msgs:
						log.Println(msg)
					case err := <-errs:
						log.Println(err)
					}
				}
			}()

			// TODO: Send order to connected node

			req := relay.OpenOrderRequest{
				Order: order.Order{
					Signature: identity.Signature{},
					ID:        []byte{},
					Type:      1,
					Parity:    1,
					Expiry:    time.Time{},
					FstCode:   1,
					SndCode:   2,
					Price:     stackint.FromUint(0),
					MaxVolume: stackint.FromUint(0),
					MinVolume: stackint.FromUint(0),
					Nonce:     stackint.FromUint(0),
				},
				OrderFragments: relay.OrderFragments{},
			}

			crypter := crypto.NewWeakCrypter()
			smpcerClient := smpcer.NewClient(&crypter, identity.MultiAddress{}, &connPool)
			// smpcerClient.OpenOrder(ctx, darknodes[0].MultiAddress, )

			wg.Wait()
		})
	})
})

func getSyncResponse(epochHash []byte, orderStatus OrderStatus) SyncResponse {
	return SyncResponse{
		Signature: []byte{},
		Epoch:     epochHash,
		Entry: &OrderbookEntry{
			Order: &Order{
				OrderId: []byte{},
				Expiry:  0,
				Type:    0,
				Tokens:  0,
			},
			OrderStatus: orderStatus,
		},
	}
}
