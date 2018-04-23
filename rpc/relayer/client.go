package relayer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc"

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
	dht      *dht.DHT
	connPool *client.ConnPool
}

// NewClient returns a Client that uses a dht.DHT and client.ConnPool when
// interacting with Relay services.
func NewClient(dht *dht.DHT, connPool *client.ConnPool) Client {
	return Client{
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
func (client *Client) Sync(ctx context.Context, orderbook *orderbook.Orderbook, epoch [32]byte, peers int) <-chan error {
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
				t := time.NewTimer(time.Minute)
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
			i := rand.Intn(len(multiAddrs))

			go func() {
				// Connect to the peer and begin synchronizing the
				// orderbook.Orderbook
				syncVals, syncErrs := client.SyncFrom(ctx, multiAddrs[i], epoch)
				atomic.AddInt64(&conns, 1)
				defer atomic.AddInt64(&conns, -1)

				for {
					select {
					case <-ctx.Done():
						return
					case val, ok := <-syncVals:
						if !ok {
							return
						}
						mergeEntry(orderbook, val)
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
func (client *Client) SyncFrom(ctx context.Context, multiAddr identity.MultiAddress, epoch [32]byte) (<-chan *SyncResponse, <-chan error) {
	responses := make(chan *SyncResponse)
	errs := make(chan error, 1)
	go func() {
		defer close(responses)
		defer close(errs)

		// FIXME: Provide verifiable signature
		relayClient := NewRelayClient(conn.ClientConn)
		request := &SyncRequest{
			Signature: []byte{},
			Epoch:     epoch[:],
		}
		stream, err := relayClient.Sync(ctx, request)
		if err != nil {
			errs <- err
			return
		}

		for {
			message, err := stream.Recv()
			if err != nil {
				errs <- err
				return
			}
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case responses <- err:
			}
		}
	}()
	return responses, errs
}

func mergeEntry(book *orderbook.Orderbook, val *SyncResponse) error {
	if val.GetEntry() == nil {
		return errors.New("cannot merge entry: nil entry")
	}
	if val.GetEpoch() == nil {
		return errors.New("cannot merge entry: nil error")
	}
	if len(val.GetEpoch()) != 32 {
		return fmt.Errorf("cannot merge entry: epoch %v", val.GetEpoch())
	}

	ord := rpc.UnmarshalOrder(val.GetEntry().GetOrder())
	epoch := [32]byte{}
	copy(epoch[:], val.GetEpoch())

	var err error
	switch val.GetEntry().GetOrderStatus() {
	case OrderStatus_Open:
		err = book.Open(orderbook.NewEntry(ord, order.Open, epoch))
	case OrderStatus_Canceled:
		err = book.Cancel(orderbook.NewEntry(ord, order.Canceled, epoch))
	case OrderStatus_Unconfirmed:
		err = book.Match(orderbook.NewEntry(ord, order.Unconfirmed, epoch))
	case OrderStatus_Confirmed:
		err = book.Confirm(orderbook.NewEntry(ord, order.Confirmed, epoch))
	case OrderStatus_Settled:
		err = book.Settle(orderbook.NewEntry(ord, order.Settled, epoch))
	default:
		return fmt.Errorf("cannot merge entry: status %v", status)
	}
	return err
}
