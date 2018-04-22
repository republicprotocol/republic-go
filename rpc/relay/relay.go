package relay

import (
	"fmt"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

type Relay struct {
	orderbook *orderbook.Orderbook
}

func NewRelay(orderbook *orderbook.Orderbook) Relay {
	return Relay{
		orderbook: orderbook,
	}
}

// Register the gRPC service to a grpc.Server.
func (relay *Relay) Register(server *grpc.Server) {
	RegisterRelayServer(server, relay)
}

func (relay *Relay) Sync(request *SyncRequest, stream Relay_SyncServer) error {
	entries := make(chan orderbook.Entry)
	defer close(entries)

	if err := relay.orderbook.Subscribe(entries); err != nil {
		return fmt.Errorf("cannot subscribe to orderbook: %v", err)
	}
	defer relay.orderbook.Unsubscribe(entries)

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case entry, ok := <-entries:
			if !ok {
				return nil
			}

			orderStatus := OrderStatus_Open
			switch entry.Status {
			case order.Open:
				orderStatus = OrderStatus_Open
			case order.Canceled:
				orderStatus = OrderStatus_Canceled
			case order.Unconfirmed:
				orderStatus = OrderStatus_Unconfirmed
			case order.Confirmed:
				orderStatus = OrderStatus_Confirmed
			case order.Settled:
				orderStatus = OrderStatus_Settled
			}

			syncResponse := &SyncResponse{
				Signature: []byte{},
				Epoch:     request.Epoch,
				Entry: &OrderbookEntry{
					Order: &Order{
						OrderId: entry.Order.ID,
						Expiry:  entry.Order.Expiry.Unix(),
						Type:    int32(entry.Order.Type),
						Tokens:  int32(0), // TODO: Use the correct token pair encoding
					},
					OrderStatus: orderStatus,
				},
			}
			if err := stream.Send(syncResponse); err != nil {
				return err
			}
		}
	}
}
