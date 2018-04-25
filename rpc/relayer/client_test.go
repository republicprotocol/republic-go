package relayer_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

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
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse(OrderStatus_Open)
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(1))
		})

		It("should not update the orderbook with a nil entry", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := SyncResponse{
				Signature: []byte{},
				Entry:     nil,
			}
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})

		It("should not update the orderbook with an invalid order status", func() {
			orderbook := orderbook.NewOrderbook(100)
			syncResponse := getSyncResponse(-1)
			err := MergeEntry(&orderbook, &syncResponse)
			Ω(err).Should(HaveOccurred())
			Ω(len(orderbook.Blocks())).Should(Equal(0))
		})
	})

	Context("when syncing orders", func() {
		It("should sync the orderbook from a specified node", func() {
			crypter := crypto.NewWeakCrypter()
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := dht.NewDHT(address, 100)
			connPool := client.NewConnPool(100)
			client := NewClient(&crypter, &dht, &connPool)

			server := grpc.NewServer()
			listener, err := net.Listen("tcp", "127.0.0.1:3000")
			Ω(err).ShouldNot(HaveOccurred())

			book := orderbook.NewOrderbook(100)
			relayer := NewRelayer(&client, &book)
			relayer.Register(server)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := server.Serve(listener)
				Ω(err).ShouldNot(HaveOccurred())
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(2 * time.Second)
				server.Stop()
			}()

			stackint := stackint.FromUint(0)
			ord := order.NewOrder(1, 1, time.Time{}, 1, 2, stackint, stackint, stackint, stackint)
			entry := orderbook.NewEntry(*ord, order.Open)
			book.Open(entry)

			ctx := context.Background()
			multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%v", address))
			var epochHash [32]byte
			msgs, errs := client.SyncFrom(ctx, multiAddr, epochHash)
			// errs := client.Sync(ctx, &book, epochHash, 3)

			// FIXME:
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case msg := <-msgs:
					log.Println(msg)
				case err := <-errs:
					log.Println(err)
				}
			}()

			wg.Wait()
		})

		It("should send an error if the connection gets closed", func() {
			crypter := crypto.NewWeakCrypter()
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := dht.NewDHT(address, 100)
			connPool := client.NewConnPool(100)
			client := NewClient(&crypter, &dht, &connPool)

			server := grpc.NewServer()
			listener, err := net.Listen("tcp", "127.0.0.1:3000")
			Ω(err).ShouldNot(HaveOccurred())

			book := orderbook.NewOrderbook(100)
			relayer := NewRelayer(&client, &book)
			relayer.Register(server)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := server.Serve(listener)
				Ω(err).ShouldNot(HaveOccurred())
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(2 * time.Second)
				server.Stop()
			}()

			ctx := context.Background()
			multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%v", address))
			var epochHash [32]byte
			msgs, errs := client.SyncFrom(ctx, multiAddr, epochHash)

			msgCount, errCount := 0, 0
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-msgs:
					msgCount++
				case <-errs:
					errCount++
				}
			}()

			wg.Wait()
			Ω(msgCount).Should(Equal(0))
			Ω(errCount).Should(Equal(1))
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
