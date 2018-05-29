package adapter

import (
	"github.com/republicprotocol/republic-go/order"
)

type OrderFragmentMapping map[string][]OrderFragment

type OrderFragment struct {
	OrderSignature string       `json:"orderSignature"`
	OrderID        string       `json:"orderId"`
	OrderType      order.Type   `json:"orderType"`
	OrderParity    order.Parity `json:"orderParity"`
	OrderExpiry    int64        `json:"orderExpiry"`
	Index          int64        `json:"index"`
	ID             string       `json:"id"`
	Tokens         string       `json:"tokens"`
	Price          []string     `json:"price"`
	Volume         []string     `json:"volume"`
	MinimumVolume  []string     `json:"minimumVolume"`
}

type OpenOrderAdapter interface {
	OpenOrder(signature string, orderFragmentMapping OrderFragmentMapping) error
}

type CancelOrderAdapter interface {
	CancelOrder(signature string, orderID string) error
}
