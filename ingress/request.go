package ingress

import "github.com/republicprotocol/republic-go/order"

// Request is an interface implemented by components that can be interpreted by
// an Ingress as a request for an action to be performed, usually involving the
// Ethereum blockchain.
type Request interface {

	// IsRequest is implemented to explicitly mark that a type is a Request. An
	// implementation of this method must do nothing.
	IsRequest()
}

// An EpochRequest is a Request for the Ingress to trigger a new epoch on the
// Ethereum blockchain.
type EpochRequest struct {
}

// IsRequest implements the Request interface.
func (req EpochRequest) IsRequest() {}

// An OpenOrderRequest is a Request for the Ingress to open an order.Order on
// the Ethereum blockchain and forward order.Fragments to their respective
// Darknodes.
type OpenOrderRequest struct {
	signature               [65]byte
	orderID                 order.ID
	orderFragmentMapping    OrderFragmentMapping
	orderFragmentEpochDepth int
}

// IsRequest implements the Request interface.
func (req OpenOrderRequest) IsRequest() {}

// A CancelOrderRequest is a Request for the Ingress to cancel an order.Order
// on the Ethereum blockchain.
type CancelOrderRequest struct {
	signature [65]byte
	orderID   order.ID
}

// IsRequest implements the Request interface.
func (req CancelOrderRequest) IsRequest() {}
