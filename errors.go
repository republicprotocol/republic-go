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

type OrderComputationError string

// NewOrderComputationError returns a new OrderComputationError for two Orders
// that have the same buy index.
func NewOrderComputationError(lhs OrderBuySell) OrderComputationError {
	rhs := OrderBuy
	if lhs == OrderBuy {
		rhs = OrderSell
	}
	return OrderComputationError(fmt.Sprintf("expected buy = %v to be computed against buy = %v", lhs, rhs))
}

// Error implements the error interface.
func (err OrderComputationError) Error() string {
	return string(err)
}
