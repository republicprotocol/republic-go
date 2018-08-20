package oracle

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrMidpointPriceNotFound is returned when attempting to read price of a
// token pair that cannot be found.
var ErrMidpointPriceNotFound = errors.New("mid-point price not found")

// MidpointPriceStorer is used for retrieving/storing MidpointPrice.
type MidpointPriceStorer interface {

	// PutMidpointPrice stores the given midPointPrice into the storer.
	// It doesn't do any nonce check and will always overwrite previous record.
	PutMidpointPrice(MidpointPrice) error

	// MidpointPrices returns the latest mid-point price data.
	MidpointPrices() (MidpointPrice, error)

	// MidpointPrice returns the mid-point price for the given token.
	MidpointPrice(token order.Tokens) (uint64, error)
}
