package relayer

import (
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// MarshalOrder into an RPC protobuf object
func MarshalOrder(ord *order.Order) *Order {
	// rpcOrder := new(Order)
	// rpcOrder.Id = &OrderId{
	// 	OrderId:   ord.ID,
	// 	Signature: ord.Signature,
	// }
	// rpcOrder.Type = int64(ord.Type)
	// rpcOrder.Parity = int64(ord.Parity)
	// rpcOrder.Expiry = ord.Expiry.Unix()

	// return rpcOrder
	panic("unimplemented")
}

// UnmarshalOrder from an RPC protobuf object.
func UnmarshalOrder(rpcOrder *Order) order.Order {
	ord := order.Order{}
	ord.ID = rpcOrder.OrderId
	ord.Type = order.Type(rpcOrder.Type)
	ord.Expiry = time.Unix(rpcOrder.Expiry, 0)
	return ord
}
