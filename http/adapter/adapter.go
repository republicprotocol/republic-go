package adapter

import (
	"github.com/republicprotocol/republic-go/order"
)

type OrderFragmentMapping map[string][]order.Fragment

type OpenOrderAdapter interface {
	OpenOrder(signature string, orderFragmentMapping OrderFragmentMapping) error
}

type CancelOrderAdapter interface {
	CancelOrder(signature string, orderID order.ID) error
}
