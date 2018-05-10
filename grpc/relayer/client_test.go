package relayer_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/relayer"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"google.golang.org/grpc"
)

var _ = Describe("Relayer", func() {

	Context("when merging entries", func() {
		It("should update the orderbook with a valid response", func() {
			orderbook := orderbook.NewOrderbook()
			syncResponse := getSyncResponse(OrderStatus_Open)
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(1))
		})

		It("should not update the orderbook with a nil entry", func() {
			orderbook := orderbook.NewOrderbook()
			syncResponse := SyncResponse{
				Signature: []byte{},
				Entry:     nil,
			}
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})

		It("should not update the orderbook with an invalid order status", func() {
			orderbook := orderbook.NewOrderbook()
			syncResponse := getSyncResponse(-1)
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})
	})

	Context("when syncing orders", func() {
		It("should sync orders from a given peer", func() {
			// Initialise the client
			crypter := crypto.NewWeakCrypter()
			dhtKey, err := crypto.RandomEcdsaKey()
			Ω(err).ShouldNot(HaveOccurred())
			dhtAddress := identity.Address(dhtKey.Address())
			dht := dht.NewDHT(dhtAddress, 100)
			connPool := client.NewConnPool(100)
			client := NewClient(&crypter, &dht, &connPool)

			// Set-up the server
			server := grpc.NewServer()
			listener, err := net.Listen("tcp", "127.0.0.1:3000")
			Ω(err).ShouldNot(HaveOccurred())

			// Create a relay and register it with the server
			book := orderbook.NewOrderbook()
			relayer := NewRelayer(&client, &book)
			relayer.Register(server)

			// Serve the server
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := server.Serve(listener)
				Ω(err).ShouldNot(HaveOccurred())
			}()

			// Add a new entry to the orderbook
			book.Open(order.Order{})

			// Synchronise the orderbook through the peer
			key, err := crypto.RandomEcdsaKey()
			Ω(err).ShouldNot(HaveOccurred())
			address := key.Address()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%s", address))
			Ω(err).ShouldNot(HaveOccurred())
			responses, errs := client.SyncFrom(context.Background(), multiAddr)

			// Check for any responses or errors
			resCount, errCount := 0, 0
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				select {
				case res := <-responses:
					submitID := []uint8(book.Blocks()[0].ID)
					resID := res.GetEntry().GetOrder().GetOrderId()
					Ω(submitID).Should(Equal(resID))
					resCount++
				case <-errs:
					errCount++
				}
			}()

			time.Sleep(time.Second)
			server.Stop()
			wg.Wait()

			// FIXME: Re-enable
			// Ω(len(book.Blocks())).Should(Equal(1))
			// Ω(resCount).Should(Equal(1))
			// Ω(errCount).Should(Equal(0))
		})

		It("should sync orders from random peers", func() {
			// Initialise the client
			crypter := crypto.NewWeakCrypter()
			dhtKey, err := crypto.RandomEcdsaKey()
			Ω(err).ShouldNot(HaveOccurred())
			dhtAddress := identity.Address(dhtKey.Address())
			dht := dht.NewDHT(dhtAddress, 100)
			connPool := client.NewConnPool(100)
			client := NewClient(&crypter, &dht, &connPool)

			// Set-up the server
			server := grpc.NewServer()
			listener, err := net.Listen("tcp", "127.0.0.1:3000")
			Ω(err).ShouldNot(HaveOccurred())

			// Create a relay and register it with the server
			relayBook := orderbook.NewOrderbook()
			relayer := NewRelayer(&client, &relayBook)
			relayer.Register(server)

			// Serve the server
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := server.Serve(listener)
				Ω(err).ShouldNot(HaveOccurred())
			}()

			// Add a new entry to the orderbook
			relayBook.Open(order.Order{})

			// Synchronise the orderbook through any peers
			key, err := crypto.RandomEcdsaKey()
			Ω(err).ShouldNot(HaveOccurred())
			address := identity.Address(key.Address())
			multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%v", address))
			Ω(err).ShouldNot(HaveOccurred())
			err = dht.UpdateMultiAddress(multiAddr)
			Ω(err).ShouldNot(HaveOccurred())

			ctx, cancel := context.WithCancel(context.Background())
			book := orderbook.NewOrderbook()
			errs := client.Sync(context.Background(), &book, 1)

			// Check for any responses or errors
			errCount := 0
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
				case err := <-errs:
					log.Println(err)
					errCount++
				}
			}()

			time.Sleep(time.Second)
			cancel()

			server.Stop()
			wg.Wait()

			// FIXME: Re-enable these checks
			// Ω(len(book.Blocks())).Should(Equal(1))
			// Ω(errCount).Should(Equal(0))
		})
	})
})

func getSyncResponse(orderStatus OrderStatus) SyncResponse {
	return SyncResponse{
		Signature: []byte{},
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
