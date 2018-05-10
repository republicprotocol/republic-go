package relayer

import (
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// UnmarshalOrder from an RPC protobuf object.
func UnmarshalOrder(rpcOrder *Order) order.Order {
	ord := order.Order{}
	ord.ID = rpcOrder.OrderId
	ord.Type = order.Type(rpcOrder.Type)
	ord.Expiry = time.Unix(rpcOrder.Expiry, 0)
	return ord
}
