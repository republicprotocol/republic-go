package rpc

import "github.com/republicprotocol/republic-go/order"

type OmeServer interface {
	OpenOrder(order.Fragment) error
}
