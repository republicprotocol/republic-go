package relayer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
)

// Client is an abstraction over the gRPC RelayClient RPCs that implements high
// level client functionality. The Client uses a dht.DHT to load all
// identity.MultiAddresses during its interactions with Relay services, and a
// client.ConnPool to reuse gRPC connections. The Client is used to synchronize
// orderbook.Orderbooks from Relay services, and merging conflicting entries.
type Client struct {
	crypter  crypto.Crypter
	dht      *dht.DHT
	connPool *client.ConnPool
}

// NewClient returns a Client that uses a dht.DHT and client.ConnPool when
// interacting with Relay services.
func NewClient(crypter crypto.Crypter, dht *dht.DHT, connPool *client.ConnPool) Client {
	return Client{
		crypter:  crypter,
		dht:      dht,
		connPool: connPool,
	}
}

// Sync the orderbook.Orderbook from a number of random peers, for a specific
// epoch (see Client.SyncTo). A context can be used to cancel or expire the
// synchronization. This function will maintain connections to the desired
// number of peers, reconnecting to new random peer whenever a connection
// fails. The Client will use its dht.DHT to pick random peers. All entries
// returned during the synchronization are written to the orderbook.Orderbook.
// Users can subscribe to the orderbook.Orderbook to receive notifications when
// an entry is synchronized.
func (client *Client) Sync(ctx context.Context, orderbook *orderbook.Orderbook, peers int) <-chan error {
	errs := make(chan error, peers+1)
	go func() {
		defer close(errs)

		conns := int64(0)
		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
			}

			// If the number of desired peers are successfully synchronizing
			// then sleep for a minute and check again
			if atomic.LoadInt64(&conns) >= int64(peers) {
				t := time.NewTimer(10 * time.Second)
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-t.C:
					continue
				}
			}

			// Select a random peer to synchronize with
			multiAddrs := client.dht.MultiAddresses()
			if len(multiAddrs) == 0 {
				return
			}
			i := rand.Intn(len(multiAddrs))

			atomic.AddInt64(&conns, 1)
			go func() {
				defer atomic.AddInt64(&conns, -1)

				// Connect to the peer and begin synchronizing the
				// orderbook.Orderbook
				syncVals, syncErrs := client.SyncFrom(ctx, multiAddrs[i])

				for {
					select {
					case <-ctx.Done():
						return
					case val, ok := <-syncVals:
						if !ok {
							return
						}
						log.Println("RECV VAL")
						MergeEntry(orderbook, val)
					case err, ok := <-syncErrs:
						if !ok {
							return
						}
						errs <- err
						return
					}
				}
			}()
		}
	}()
	return errs
}

// SyncFrom a peer. A context can be used to cancel or expire the
// synchronization. Once this stops synchronizing, the cancellation and
// expiration of the Context will do nothing. In a background goroutine, the
// Client will synchronize from the given peer, expecting entries for a
// specific epoch. All entries and errors will be written to channels that are
// returned immediately.
func (client *Client) SyncFrom(ctx context.Context, multiAddr identity.MultiAddress) (<-chan *SyncResponse, <-chan error) {
	responses := make(chan *SyncResponse)
	errs := make(chan error, 1)
	go func() {
		defer close(responses)
		defer close(errs)

		conn, err := client.connPool.Dial(ctx, multiAddr)
		if err != nil {
			errs <- fmt.Errorf("cannot dial %v: %v", multiAddr, err)
			return
		}
		defer conn.Close()

		relayClient := NewRelayClient(conn.ClientConn)
		requestSignature, err := client.crypter.Sign(client.dht.Address)
		if err != nil {
			errs <- fmt.Errorf("cannot sign sync request: %v", err)
			return
		}
		request := &SyncRequest{
			Signature: requestSignature,
			Address:   client.dht.Address.String(),
		}
		stream, err := relayClient.Sync(ctx, request)
		if err != nil {
			errs <- fmt.Errorf("cannot sync with %v: %v", multiAddr, err)
			return
		}

		for {
			message, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				errs <- err
				return
			}
			if message == nil {
				continue
			}
			log.Println("RECV message")
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case responses <- message:
			}
		}
	}()
	return responses, errs
}

func MergeEntry(book *orderbook.Orderbook, val *SyncResponse) error {
	if val.GetEntry() == nil {
		return errors.New("cannot merge entry: nil entry")
	}

	ord := UnmarshalOrder(val.GetEntry().GetOrder())

	var err error
	switch val.GetEntry().GetOrderStatus() {
	case OrderStatus_Open:
		err = book.Open(orderbook.NewEntry(ord, order.Open))
	case OrderStatus_Canceled:
		err = book.Cancel(ord.ID)
	case OrderStatus_Unconfirmed:
		err = book.Match(orderbook.NewEntry(ord, order.Unconfirmed))
	case OrderStatus_Confirmed:
		log.Println("CONFIRM")
		err = book.Confirm(orderbook.NewEntry(ord, order.Confirmed))
	case OrderStatus_Settled:
		err = book.Settle(orderbook.NewEntry(ord, order.Settled))
	default:
		return fmt.Errorf("cannot merge entry: status %v", val.GetEntry().GetOrderStatus())
	}
	return err
}
