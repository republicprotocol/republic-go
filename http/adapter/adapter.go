package adapter

import (
	"github.com/republicprotocol/republic-go/order"
)

// An OrderFragmentMapping maps pods to encrypted order fragments represented
// as a JSON object. This representation is useful for HTTP drivers.
type OrderFragmentMapping map[string][]OrderFragment

// OrderFragment is an order.EncryptedFragment, encrypted by the trader. It
// stores the an index that identifies which index of shamir.Shares are stored
// in the OrderFragment. It is represented as a JSON object. This
// representation is useful for HTTP drivers.
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

// An OpenOrderAdapter can be used to open an order.Order by sending an
// OrderFragmentMapping to the Darknodes in the network.
type OpenOrderAdapter interface {
	OpenOrder(signature string, orderFragmentMapping OrderFragmentMapping) error
}

// A CancelOrderAdapter can be used to cancel an order.Order by sending a
// signed cancelation message to the Ethereum blockchain where all Darknodes in
// the network will observe it.
type CancelOrderAdapter interface {
	CancelOrder(signature string, orderID string) error
}
