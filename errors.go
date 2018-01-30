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

// A ResultFragmentationError occurs when two OrderFragments are combined that
// are incompatible.
type ResultFragmentationError string

// NewResultFragmentationError returns a new ResultFragmentationError for two
// OrderFragments that have the same buy-sell type.
func NewResultFragmentationError(lhs OrderParity) ResultFragmentationError {
	rhs := OrderParityBuy
	if lhs == OrderParityBuy {
		rhs = OrderParitySell
	}
	return ResultFragmentationError(fmt.Sprintf("expected buy/sell = %v to be computed against buy/sell = %v", lhs, rhs))
}

// Error implements the error interface.
func (err ResultFragmentationError) Error() string {
	return string(err)
}
