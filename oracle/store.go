package oracle

// MidpointPriceStorer is used for retrieving/storing MidpointPrice.
type MidpointPriceStorer interface {

	// PutMidpointPrice stores the given midPointPrice into the storer.
	// It doesn't do any nonce check and will always overwrite previous record.
	PutMidpointPrice(MidpointPrice) error

	// MidpointPrice returns the latest mid-point price of the given token.
	MidpointPrice() (MidpointPrice, error)
}
