package compute

import "fmt"

// An OrderFragmentationError occurs when a computation occurs over two
// OrderFragments that have different indices.
type OrderFragmentationError string

// NewOrderFragmentationError returns a new OrderFragmentationError for
// OrderFragments with indices i and j.
func NewOrderFragmentationError(i, j int64) OrderFragmentationError {
	return OrderFragmentationError(fmt.Sprintf("expected i = %v to be equal to j = %v", i, j))
}

// Error implements the error interface.
func (err OrderFragmentationError) Error() string {
	return string(err)
}

// An OrderParityError occurs when two OrderFragments have incompatible
// OrderParity values.
type OrderParityError string

// NewOrderParityError returns a new OrderParityError for two OrderFragments
// that have the same OrderParity.
func NewOrderParityError(lhs OrderParity) OrderParityError {
	rhs := OrderParityBuy
	if lhs == OrderParityBuy {
		rhs = OrderParitySell
	}
	return OrderParityError(fmt.Sprintf("expected buy/sell = %v to be computed against buy/sell = %v", lhs, rhs))
}

// Error implements the error interface.
func (err OrderParityError) Error() string {
	return string(err)
}
